package rspace

import (
	"math/rand"
	"time"
	"testing"
	"strings"
	"fmt"
)
var seededRand *rand.Rand = rand.New(
	  rand.NewSource(time.Now().UnixNano()))

const alphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_0123456789"

func stringWithCharset(length int, charset string) string {
  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}

func randomAlphanumeric (length int) string {
	return stringWithCharset(length, alphanumeric)
}


func assertIntEquals(t *testing.T, expected int, actual int, message string) {
	var b strings.Builder
	var isFail bool = false
	if actual != expected {
		isFail = true
		b.WriteString(fmt.Sprintf("Expected [%d] but was [%d]", expected, actual))
	}
	if len(message) > 0 {
		b.WriteString("\n" +message)
	}
	if isFail {
		fail(t, b.String())
	}
}


