package tools

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestCheck проверяет функцию-декоратор для обработки ошибки
func TestCheck(t *testing.T) {
	assert.PanicsWithError(
		t,
		"test error",
		func() { Check(errors.New("test error")) },
	)
}
