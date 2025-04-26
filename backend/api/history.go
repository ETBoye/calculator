package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/etboye/calculator/persistence"
)

// We copy the pattern from https://opensource.zalando.com/restful-api-guidelines/#pagination
type CalculationsPage struct { // Page as in pagination object, not web page
	Self  *string                  `json:"self"`
	First *string                  `json:"first"`
	Prev  *string                  `json:"prev"`
	Next  *string                  `json:"next"`
	Last  *string                  `json:"last"`
	Items []persistence.HistoryRow `json:"items"`
	Error *string                  `json:"error"`
}

type SessionHistoryHandler interface {
	GetSessionHistory(sessionId string, cursorQuery string) SimpleHttpResponse[CalculationsPage]
}

type StandardSessionHistoryHandler struct {
	persistenceClient persistence.PersistenceClient
}

func NewStandardSessionHistoryHandler(persistenceClient persistence.PersistenceClient) StandardSessionHistoryHandler {
	return StandardSessionHistoryHandler{persistenceClient: persistenceClient}
}

func (historyHandler StandardSessionHistoryHandler) GetSessionHistory(sessionId string, cursorQuery string) SimpleHttpResponse[CalculationsPage] {
	sessionIdValidationErr := validateSessionId(sessionId)

	if sessionIdValidationErr != nil {
		errorId := sessionIdValidationErr.Error()
		return SimpleHttpResponse[CalculationsPage]{ // TODO: Test
			Status:   http.StatusBadRequest,
			Response: CalculationsPage{Error: &errorId},
		}
	}

	var unmappedCalculationsPage persistence.CalculationsPageObject
	var persistenceError error

	if cursorQuery == "" {
		unmappedCalculationsPage, persistenceError = historyHandler.persistenceClient.GetSessionHistoryFirstPage(sessionId)
	} else {
		cursor, cursorParseErr := strconv.ParseInt(cursorQuery, 10, 64)

		if cursorParseErr != nil {
			errorId := "CURSOR_PARSING_ERROR"
			log.Printf("Could not parse %s as int", cursorQuery)
			return SimpleHttpResponse[CalculationsPage]{ // TODO: Test
				Status:   http.StatusBadRequest,
				Response: CalculationsPage{Error: &errorId},
			}
		}
		unmappedCalculationsPage, persistenceError = historyHandler.persistenceClient.GetSessionHistory(sessionId, cursor)
	}

	if persistenceError != nil {
		errorId := persistenceError.Error()
		return SimpleHttpResponse[CalculationsPage]{ // TODO: Test
			Status:   http.StatusInternalServerError,
			Response: CalculationsPage{Error: &errorId},
		}
	}

	return SimpleHttpResponse[CalculationsPage]{
		Status:   http.StatusOK,
		Response: mapPage(sessionId, unmappedCalculationsPage),
	}
}

func mapPage(sessionId string, unmappedCalculationsPage persistence.CalculationsPageObject) CalculationsPage {
	return CalculationsPage{
		Self:  getEndpointUrl(sessionId, unmappedCalculationsPage.Self),
		First: getEndpointUrl(sessionId, unmappedCalculationsPage.First),
		Prev:  getEndpointUrl(sessionId, unmappedCalculationsPage.Prev),
		Next:  getEndpointUrl(sessionId, unmappedCalculationsPage.Next),
		Last:  getEndpointUrl(sessionId, unmappedCalculationsPage.Last),
		Items: unmappedCalculationsPage.Items,
	}
}

func getEndpointUrl(sessionId string, cursor *int64) *string {
	if cursor == nil {
		return nil
	}
	result := fmt.Sprintf("/sessions/%s/history?cursor=%d", sessionId, *cursor)
	return &result
}
