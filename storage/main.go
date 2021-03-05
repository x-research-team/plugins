package main

import (
	"storage/component"

	"github.com/x-research-team/contract"
)

// Init Load plugin with all components
func Init() contract.KernelModule {
	return component.New(
		component.ConnectTo("Server=localhost;Database=kernel;Uid=root;Pwd=root;"),
	)
}

func main() {
}
