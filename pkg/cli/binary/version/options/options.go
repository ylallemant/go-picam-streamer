package options

var (
	Domain  = "picam-streamer"
	Current = NewOptions()
)

func NewOptions() *Options {
	options := new(Options)

	return options
}

type Options struct {
	Semver    bool
	Commit    bool
	Separator string
}
