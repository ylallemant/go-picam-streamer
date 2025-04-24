//go:build linux

package camera

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
	"github.com/ylallemant/go-picam-streamer/pkg/api"
)

func Device(ctx context.Context, options *api.CameraOption) (*device.Device, error) {
	cam, err := device.Open(
		api.DefaultDevice,
		//device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 640, Height: 480}),
		device.WithPixFormat(v4l2.PixFormat{
			PixelFormat: v4l2.PixelFmtMJPEG,
			Width:       options.CaptureWidth,
			Height:      options.CaptureHeight,
		}),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open camera device %s", api.DefaultDevice)
	}

	log.Info().Msgf("device name:            %s", cam.Name())
	log.Info().Msgf("device file descriptor: %v", cam.Fd())
	log.Info().Msgf("device capability:      %s", cam.Capability())
	log.Info().Msgf("device buffer type:     %v", cam.BufferType())

	if err := cam.Start(ctx); err != nil {
		log.Fatal().Msgf("camera start: %s", err)
	}

	return cam, nil
}
