package util

import "net/http"

func GetAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(Attempts).(int); ok {
		return attempts
	}
	return 1
}

func GetRetriesFromContext(r *http.Request) int {
	if retries, ok := r.Context().Value(Retries).(int); ok {
		return retries
	}
	return 0
}
