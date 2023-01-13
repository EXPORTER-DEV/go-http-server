package utils

func Find[T any](arr []T, matcher func(item T, index int) bool) (bool, T) {
	for index, value := range arr {
		if matcher(value, index) {
			return true, value
		}
	}

	var emptyResult T

	return false, emptyResult
}

func Some[T any](arr []T, matcher func(item T, index int) bool) bool {
	for index, value := range arr {
		if matcher(value, index) {
			return true
		}
	}

	return false
}

func Every[T any](arr []T, matcher func(item T, index int) bool) bool {
	if len(arr) == 0 {
		return false
	}

	for index, value := range arr {
		if !matcher(value, index) {
			return false
		}
	}

	return true
}

func Filter[T any](arr []T, matcher func(item T, index int) bool) []T {
	var result []T

	for index, value := range arr {
		if matcher(value, index) {
			result = append(result, value)
		}
	}

	return result
}
