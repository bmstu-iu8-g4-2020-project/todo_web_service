package handlers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	str, err := strconv.Atoi("Hello, world!")
	assert.NotNil(t, err)
	fmt.Println(str)
}
