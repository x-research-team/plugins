package main

import (
	"queue/component"

	"github.com/x-research-team/contract"
)

// Init Load plugin with all components
func Init() contract.KernelModule {
	return component.New()
}

func main() {
}
