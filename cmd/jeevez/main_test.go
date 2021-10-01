package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestMainFunc смок тест на ошибку в основной функции запуска cli
func TestMainFunc(t *testing.T) {
	assert.Panics(t, main)
}
