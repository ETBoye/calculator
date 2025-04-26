package main

import (
	"log"
	"os"

	"github.com/etboye/calculator/api"
	"github.com/etboye/calculator/calculation"
	"github.com/etboye/calculator/persistence"
	"github.com/etboye/calculator/server"
)

func main() {
	server := server.GinServer{}

	persistenceClient, err := persistence.InitPostgresClient()

	if err != nil {
		log.Fatalf("Received error initialising persistence client: %s", err.Error())
		os.Exit(1)
	}

	calculator := calculation.NewDefaultExpressionCalculator()

	endpoints := api.Endpoints{
		ComputationHandler:    api.NewStandardComputationHandler(&calculator, persistenceClient),
		SessionHistoryHandler: api.NewStandardSessionHistoryHandler(persistenceClient),
	}
	server.StartServer(endpoints)
}
