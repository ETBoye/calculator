package main

import (
	"github.com/etboye/calculator/api"
	"github.com/etboye/calculator/calculation"
	"github.com/etboye/calculator/server"
)

type Application struct {
	calculator calculation.Calculator
}

func main() {
	server := server.GinServer{}

	endpoints := api.Endpoints{
		ComputationHandler: api.NewStandardComputationHandler(&calculation.StubCalculator{}),
	}
	server.StartServer(endpoints)
}
