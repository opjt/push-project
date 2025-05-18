package main

import "push/linker/core/bootstrap"

func main() {
	bootstrap.RunServer(bootstrap.CommonModules)
}
