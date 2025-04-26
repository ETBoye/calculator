package api

type SimpleHttpResponse[Body any] struct {
	Status   int
	Response Body
}

type Endpoints struct {
	ComputationHandler    ComputationHandler
	SessionHistoryHandler SessionHistoryHandler
}
