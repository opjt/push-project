package main

import "push/linker/internal/core/bootstrap"

func main() {
	bootstrap.RunServer(bootstrap.CommonModules)
}
