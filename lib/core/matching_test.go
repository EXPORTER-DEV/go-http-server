package server

import (
	"reflect"
	"testing"
)

func TestMatchingExecute(t *testing.T) {
	t.Run("Check base matching", func(t *testing.T) {
		matching := Matching{
			RequestedPath: "/index",
			HandlerPath:   "/index",
		}

		result, params := matching.Execute(MatchingExecuteOptions{
			ParseParams: false,
			IsRegexp:    false,
		})

		if result != true {
			t.Fatalf("%v expected to match", matching)
		}

		if params != nil {
			t.Fatalf("%v expected to be nil", params)
		}
	})

	t.Run("Check params matching", func(t *testing.T) {
		matching := Matching{
			RequestedPath: "/index/123/321",
			HandlerPath:   "/index/:id/:style",
		}

		result, params := matching.Execute(MatchingExecuteOptions{
			ParseParams: true,
			IsRegexp:    false,
		})

		if result != true {
			t.Fatalf("%v expected to match", matching)
		}

		expectedParams := Params{
			"id":    "123",
			"style": "321",
		}

		if !reflect.DeepEqual(params, expectedParams) {
			t.Fatalf("%v expected to be equal %v", params, expectedParams)
		}
	})

	t.Run("Check params matching", func(t *testing.T) {
		matching := Matching{
			RequestedPath: "/index/123/321/1",
			HandlerPath:   "/index/:id/:style/:d",
		}

		result, params := matching.Execute(MatchingExecuteOptions{
			ParseParams: true,
			IsRegexp:    false,
		})

		if result != true {
			t.Fatalf("%v expected to match", matching)
		}

		expectedParams := Params{
			"id":    "123",
			"style": "321",
			"d":     "1",
		}

		if !reflect.DeepEqual(params, expectedParams) {
			t.Fatalf("%v expected to be equal %v", params, expectedParams)
		}
	})

	t.Run("Check params not matching", func(t *testing.T) {
		matching := Matching{
			RequestedPath: "/index/123//1",
			HandlerPath:   "/index/:id/:style/:d",
		}

		result, _ := matching.Execute(MatchingExecuteOptions{
			ParseParams: true,
			IsRegexp:    false,
		})

		if result == true {
			t.Fatalf("%v expected not match", matching)
		}
	})
}
