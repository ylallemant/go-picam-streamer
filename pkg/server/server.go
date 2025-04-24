package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/ylallemant/go-picam-streamer/pkg/api"
	"github.com/ylallemant/go-picam-streamer/pkg/camera"
)

func New(serverOptions *api.ServerOptions, cameraOptions *api.CameraOption) (*server, error) {
	svr := new(server)
	svr.mux = http.NewServeMux()
	svr.port = serverOptions.Port
	svr.binding = serverOptions.Address

	ctx, cancel := context.WithCancel(context.Background())
	svr.ctx = ctx
	svr.cancelFunc = cancel

	cam, err := camera.New(svr.ctx, cameraOptions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open camera device")
	}

	svr.camera = cam
	log.Info().Msgf("camera started")
	svr.frames = svr.camera.ReadFrames()

	var staticFS = fs.FS(staticFiles)
	htmlContent, err := fs.Sub(staticFS, "static")
	if err != nil {
		return nil, errors.Wrap(err, "failed mount static files")
	}
	fileserver := http.FileServer(http.FS(htmlContent))

	svr.mux.Handle("/", fileserver)
	svr.mux.HandleFunc("/stream", svr.imageServ)

	svr.http = &http.Server{
		Handler: svr.mux,
	}

	return svr, nil
}

//go:embed static
var staticFiles embed.FS

type server struct {
	http       *http.Server
	mux        *http.ServeMux
	camera     api.Camera
	ctx        context.Context
	cancelFunc context.CancelFunc
	port       string
	binding    string
	frames     <-chan []byte
}

func (i *server) Start() error {
	addr := fmt.Sprintf("%s:%s", i.binding, i.port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrapf(err, "failed to initiate listener on %s", addr)
	}

	log.Info().Msgf("Serving images: [%s/stream]", addr)
	return i.http.Serve(listener)
}

func (i *server) imageServ(w http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("request stream")
	mimeWriter := multipart.NewWriter(w)
	w.Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", mimeWriter.Boundary()))
	partHeader := make(textproto.MIMEHeader)
	partHeader.Add("Content-Type", "image/jpeg")

	var frame []byte
	for frame = range i.frames {
		log.Trace().Msgf("process frame")
		partWriter, err := mimeWriter.CreatePart(partHeader)
		if err != nil {
			log.Printf("failed to create multi-part writer: %s", err)
			return
		}

		if _, err := partWriter.Write(frame); err != nil {
			log.Printf("failed to write image: %s", err)
		}
	}
}
