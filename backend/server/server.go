package server

import (
	"github.com/etboye/calculator/api"
)

type Server interface {
	StartServer(endpoints api.Endpoints) error
}
