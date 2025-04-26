package api

import (
	"errors"
	"log"
	"regexp"
)

var sessionIdPatternString string = `^[a-zA-Z0-9\-]*$`
var sessionIdPattern *regexp.Regexp = regexp.MustCompile(sessionIdPatternString)
var sessionIdMinLength int = 1
var sessionIdMaxLength int = 100
var sessionIdValidationErrorId = "SESSION_ID_VALIDATION_ERROR"

func validateSessionId(sessionId string) error {
	log.Println("Got sessionId: ", sessionId)

	if !sessionIdPattern.MatchString(sessionId) {
		log.Printf("sessionId must match pattern %s", sessionIdPatternString)
		return errors.New(sessionIdValidationErrorId)
	}

	if len(sessionId) > sessionIdMaxLength {
		log.Printf("sessionId must have length at most %d", sessionIdMaxLength) // TODO: test
		return errors.New(sessionIdValidationErrorId)
	}

	if len(sessionId) < sessionIdMinLength {
		log.Printf("sessionId must have length at minimum %d", sessionIdMinLength) // TODO: test
		return errors.New(sessionIdValidationErrorId)
	}

	return nil
}
