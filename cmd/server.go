package main

import (
	"context"
	bootstrap "main/internal/Infrastructure/Bootstrap"
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
)

func main() {
	container := dependency_injection.GetContainer()
	defer container.Telemetry.Shutdown(context.Background())
	_, span := container.Telemetry.TraceStart(context.Background(), "main")
	defer span.End()
	r := bootstrap.BootstrapGin(*container)
	r.Run()
}
