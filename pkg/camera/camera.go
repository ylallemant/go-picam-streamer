package camera

import (
	"context"

	"github.com/pkg/errors"
	"github.com/ylallemant/go-picam-streamer/pkg/api"
)

func New(ctx context.Context, options *api.CameraOption) (*camera, error) {
	instance := new(camera)

	cam, err := Device(ctx, options)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initialise camera device %s", api.DefaultDevice)
	}

	instance.v4l2 = cam

	return instance, nil
}

var _ api.Camera = &camera{}

type camera struct {
	v4l2 api.Device
}

func (i *camera) ReadFrames() <-chan []byte {
	return i.v4l2.GetOutput()
}
