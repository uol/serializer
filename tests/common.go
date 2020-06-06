package tests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// CheckNullErrorValidation - checks for a null warning message
func CheckNullErrorValidation(t *testing.T, err error) bool {

	if !assert.Error(t, err, "expected validation error") {
		return false
	}

	return assert.True(t, strings.Contains(err.Error(), "is null"), "expected \"is null\" on message")
}
