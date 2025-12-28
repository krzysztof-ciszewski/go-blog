package main

import (
	"context"
	bootstrap "main/internal/Infrastructure/Bootstrap"
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
)

func main() {
	container := dependency_injection.GetContainer()
	defer container.Telemetry.Shutdown(context.Background())
	r := bootstrap.BootstrapGin(*container)
	r.Run()
}
