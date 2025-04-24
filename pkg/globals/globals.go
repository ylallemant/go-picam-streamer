package globals

import "github.com/rs/zerolog"

var (
	Current = new(Globals)
)

type Globals struct {
	ConfigPath     string
	FallbackConfig bool
	Debug          bool
	LogLevel       string
	NonBlocking    bool
}

func ProcessGlobals() {
	if Current.Debug {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else if Current.LogLevel != "" {
		switch Current.LogLevel {
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		default:
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		}
	}
}
