package persistence

import (
	"context"
	"log"
)

type cursorPositions struct {
	Self  *int64
	First *int64
	Prev  *int64
	Next  *int64
	Last  *int64
}

func (postgresClient PostgresClient) computeCursorPositions(sessionId string, currentCursorOrNil *int64) (cursorPositions, error) {
	calculationIds, err := postgresClient.getAllCalculationIdsInSessionInDescendingOrder(sessionId)

	if err != nil {
		return cursorPositions{}, err
	}

	if len(calculationIds) == 0 {
		return cursorPositions{
			First: nil,
			Prev:  nil,
			Next:  nil,
			Last:  nil,
		}, nil
	}

	first := calculationIds[0]

	var idxOfLast int
	if len(calculationIds)%pageMaxNumberOfItems == 0 {
		idxOfLast = len(calculationIds) - pageMaxNumberOfItems
	} else {
		idxOfLast = pageMaxNumberOfItems * (len(calculationIds) / pageMaxNumberOfItems)
	}

	last := calculationIds[idxOfLast]

	idxOfCurrent := 0

	if currentCursorOrNil != nil {
		currentCursor := *currentCursorOrNil
		log.Println("computeCursorPositions: Got current cursor", currentCursor)

		for idxOfCurrent < len(calculationIds) && calculationIds[idxOfCurrent] > currentCursor {
			idxOfCurrent++
		}

		log.Println("computeCursorPositions: Computed idxOfCurrent", idxOfCurrent)

		if idxOfCurrent == len(calculationIds) {
			// request is asking for a list with calculation ids below or equal to the lowest calculation id in the list
			// The result will always be an empty list

			log.Printf("Got into special case where idxOfCurrent=%v is out of bounds", idxOfCurrent)

			result := cursorPositions{
				Self:  &currentCursor,
				First: &first,
				Prev:  &last,
				Next:  nil,
				Last:  &last,
			}
			return result, nil
		}

		log.Printf("computeCursorPositions: currentCursor: %d, idxOfCurrent: %d", currentCursor, idxOfCurrent)

	}

	idxOfPrev := idxOfCurrent - pageMaxNumberOfItems
	if idxOfPrev < 0 {
		idxOfPrev = 0
	}

	idxOfNext := idxOfCurrent + pageMaxNumberOfItems

	if idxOfNext >= len(calculationIds) {
		idxOfNext = len(calculationIds) - 1
	}

	self := calculationIds[idxOfCurrent]

	prevResult := &calculationIds[idxOfPrev]
	if self == first {
		prevResult = nil
	}

	nextResult := &calculationIds[idxOfNext]
	if self == last {
		nextResult = nil
	}

	return cursorPositions{
		Self:  &self,
		First: &first,
		Prev:  prevResult,
		Next:  nextResult,
		Last:  &last,
	}, nil
}

func (postgresClient PostgresClient) getAllCalculationIdsInSessionInDescendingOrder(sessionId string) ([]int64, error) {
	rows, err := postgresClient.conn.Query(context.Background(), "SELECT calculationId FROM history WHERE sessionId=$1 ORDER BY calculationId DESC", sessionId)

	if err != nil {
		return []int64{}, err
	}

	result := make([]int64, 0)

	var calculationId int64
	for rows.Next() {
		err := rows.Scan(&calculationId)

		if err != nil {
			return []int64{}, err
		}

		result = append(result, calculationId)
	}

	return result, nil
}

func buildCalculationsPageObject(cursorPositions cursorPositions, items []HistoryRow) CalculationsPageObject {
	return CalculationsPageObject{
		Self:  cursorPositions.Self,
		First: cursorPositions.First,
		Prev:  cursorPositions.Prev,
		Next:  cursorPositions.Next,
		Last:  cursorPositions.Last,
		Items: items,
	}
}
