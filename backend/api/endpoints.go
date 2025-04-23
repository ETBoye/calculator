package api

// I want services to control the status and body, but maybe not headers etc
type SimpleHttpResponse[Body any] struct {
	Status   int
	Response Body
}

type Endpoints struct {
	ComputationHandler ComputationHandler
}
