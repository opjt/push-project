package main

import (
	"push/dispatcher/internal/core"
)

func main() {
	core.RunServer(core.Modules)
}
