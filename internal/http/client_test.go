package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

const (
	port          = 8080
	scheme        = "http"
	successAction = "/success"
	errorAction   = "/error"
	baseUrl       = "localhost"
)

type TestData struct {
	SimpleString     string
	SimpleInt        int
	IncludedTestData IncludedTestData
}

type IncludedTestData struct {
	IncludedSimpleString string
}

type ErrorData struct {
	Code int
	Text string
}

func TestClient_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		want := &TestData{
			SimpleString:     "simple string",
			SimpleInt:        1000,
			IncludedTestData: IncludedTestData{IncludedSimpleString: "included simple string"},
		}

		serve(want, nil)

		sut := NewClient(fmt.Sprintf("%s:%s:%v", scheme, baseUrl, port), time.Second*30)

		got := new(TestData)
		err := sut.Get(ctx, successAction, got)

		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func serve(successResult *TestData, errorResult *ErrorData) {
	http.HandleFunc(
		successAction,
		func(w http.ResponseWriter, req *http.Request) {
			b, err := json.Marshal(successResult)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(b)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		},
	)

	http.HandleFunc(
		errorAction,
		func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, errorResult.Text, errorResult.Code)
		},
	)

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
		if err != nil {
			panic(fmt.Sprintf("cannot serve test service: %v", err))
		}
	}()
}
