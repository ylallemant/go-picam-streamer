package options

var (
	Current = NewOptions()
)

func NewOptions() *Options {
	options := new(Options)

	options.Port = "8080"
	options.Address = "0.0.0.0"

	options.CaptureHeight = 520
	options.CaptureWidth = 960

	return options
}

type Options struct {
	Port          string
	Address       string
	CaptureHeight int
	CaptureWidth  int
}
