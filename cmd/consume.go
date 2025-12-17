package main

import (
	"context"
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
)

func main() {
	container := dependency_injection.GetContainer()

	if err := container.Router.Run(context.Background()); err != nil {
		panic(err)
	}
}
