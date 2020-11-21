package handlers

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(t *testing.T) {
	str, err := strconv.Atoi("Hello, world!")
	assert.NotNil(t, err)
	fmt.Println(str)
}
