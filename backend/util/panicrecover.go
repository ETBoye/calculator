package util

import "fmt"

type withError[T any] struct {
	t   T
	err error
}

func RecoverFromPanicWithError[T any](f func() (T, error), recoverItem T, recoverError error, recoverMessage string) (T, error) {
	resultWithError := RecoverFromPanic[withError[T]](func() withError[T] {
		t, err := f()

		result := withError[T]{t: t, err: err}
		return result
	}, withError[T]{t: recoverItem, err: recoverError}, recoverMessage)

	return resultWithError.t, resultWithError.err
}

func RecoverFromPanic[T any](f func() T, recoverItem T, recoverMessage string) T {
	c := make(chan T)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(recoverMessage, r)
				c <- recoverItem
			}
		}()

		c <- f()
	}()

	result := <-c
	return result
}
