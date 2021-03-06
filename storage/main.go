package main

import (
	"storage/component"
	"storage/component/dsn"

	"github.com/x-research-team/contract"
)

// Init Load plugin with all components
func Init() contract.KernelModule {
	return component.New(
		component.ConnectTo(dsn.Parse()),
	)
}

func main() {
}
