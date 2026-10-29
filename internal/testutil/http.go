package testutil

import (
	"context"
	"net/http"
	"time"
)

func CompareHeaders(m1, m2 http.Header) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		v2, ok := m2[k]

		if !ok {
			return false
		}

		if len(v1) != len(v2) {
			return false
		}

		for i := range v1 {
			if v1[i] != v2[i] {
				return false
			}
		}
	}
	return true
}

func GetNeededCtx(deadlineContext bool) (context.Context, context.CancelFunc) {
	if deadlineContext {
		return context.WithDeadline(context.Background(), time.Time{})
	}
	return context.WithCancel(context.Background())

}
