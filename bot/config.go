package bot

import "os"

type Config struct {
	Token string
	Timeout int
	CoronaApi
}

type CoronaApi struct {
	XRapidAPIKey string
}

func NewConfig() *Config {
	return &Config{
		Token: os.Getenv("CORONA_BOT_TOKEN"),
		Timeout: 50,
		CoronaApi: CoronaApi{ XRapidAPIKey: os.Getenv("CORONA_API_KEY") },
	}
}