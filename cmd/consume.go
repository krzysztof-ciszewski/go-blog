package main

import (
	"context"
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
)

func main() {
	container := dependency_injection.GetContainer()
	defer container.Router.Close()
	defer container.Telemetry.Shutdown(context.Background())
	defer container.SessionStore.Close()
	if err := container.Router.Run(context.Background()); err != nil {
		panic(err)
	}
}
