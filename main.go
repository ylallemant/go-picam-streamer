package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	frames <-chan []byte
)

//go:embed static
var staticFiles embed.FS

func imageServ(w http.ResponseWriter, req *http.Request) {
	mimeWriter := multipart.NewWriter(w)
	w.Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", mimeWriter.Boundary()))
	partHeader := make(textproto.MIMEHeader)
	partHeader.Add("Content-Type", "image/jpeg")

	var frame []byte
	for frame = range frames {
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

func main() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		// TODO find a solution for bad output "../../../../../../../../pkg/nlp/tokenizer.go:229"
		// start := strings.Index(file, "/pkg")
		// if start > -1 {
		// 	return file[start:] + ":" + strconv.Itoa(line)
		// }
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()

	port := "9090"
	binding := "0.0.0.0"
	devName := "/dev/video0"

	fileserver := http.FileServer(http.Dir("./static"))

	flag.StringVar(&devName, "d", devName, "device name (path)")
	flag.StringVar(&port, "p", port, "webcam service port")
	flag.StringVar(&binding, "b", binding, "dinding address")

	camera, err := device.Open(
		devName,
		//device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 640, Height: 480}),
		device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 960, Height: 520}),
	)
	if err != nil {
		log.Fatal().Msgf("failed to open device: %s", err)
	}
	defer camera.Close()

	log.Info().Msgf("device name:            %s", camera.Name())
	log.Info().Msgf("device file descriptor: %v", camera.Fd())
	log.Info().Msgf("device capability:      %s", camera.Capability())
	log.Info().Msgf("device buffer type:     %v", camera.BufferType())

	if err := camera.Start(context.TODO()); err != nil {
		log.Fatal().Msgf("camera start: %s", err)
	}

	frames = camera.GetOutput()

	log.Info().Msgf("Serving images: [%s/stream]", fmt.Sprintf("%s:%s", binding, port))
	http.Handle("/", fileserver)
	http.HandleFunc("/stream", imageServ)
	log.Fatal().Msgf("%s", http.ListenAndServe(fmt.Sprintf("%s:%s", binding, port), nil))
}
