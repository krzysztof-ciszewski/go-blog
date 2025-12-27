package config

import "os"

type TelemetryConfig struct {
	ServiceName        string
	ServiceVersion     string
	ServiceEnvironment string
}

func GetTelemetryConfig() *TelemetryConfig {
	return &TelemetryConfig{
		ServiceName:        os.Getenv("SERVICE_NAME"),
		ServiceVersion:     os.Getenv("SERVICE_VERSION"),
		ServiceEnvironment: os.Getenv("SERVICE_ENVIRONMENT"),
	}
}
