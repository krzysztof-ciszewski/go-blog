package main

import (
	bootstrap "main/internal/Infrastructure/Bootstrap"
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
)

func main() {
	container := dependency_injection.GetContainer()
	r := bootstrap.BootstrapGin(*container)
	r.Run()
}
