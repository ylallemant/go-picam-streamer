package api

const (
	DefaultDevice = "/dev/video0"
)

type Camera interface {
	ReadFrames() <-chan []byte
}

type Device interface {
	GetOutput() <-chan []byte
}

type CameraOption struct {
	CaptureHeight int
	CaptureWidth  int
}
