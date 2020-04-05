package rspace

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func fail(t *testing.T, message string) {
	t.Errorf(message)
}

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

func randomAlphanumeric(length int) string {
	return stringWithCharset(length, alphanumeric)
}

type Testable interface {
	IsEqual() bool
	String() string
}
type IntTestResult struct {
	Expected int
	Actual   int
}

func (r IntTestResult) String() string {
	return fmt.Sprintf("Expected [%d] but was [%d]", r.Expected, r.Actual)
}
func (r IntTestResult) IsEqual() bool {
	return r.Expected == r.Actual
}

type StringTestResult struct {
	Expected string
	Actual   string
}

func (r StringTestResult) String() string {
	return fmt.Sprintf("Expected [%s] but was [%s]", r.Expected, r.Actual)
}
func (r StringTestResult) IsEqual() bool {
	return r.Expected == r.Actual
}

func assertIntEquals(t *testing.T, expected int, actual int, message string) {
	result := IntTestResult{expected, actual}
	_assertEquals(t, result, message)
}
func assertStringEquals(t *testing.T, expected string, actual string, message string) {
	result := StringTestResult{expected, actual}
	_assertEquals(t, result, message)
}

func assertNotNil(t *testing.T, toTest interface{}, message string) {
	if toTest == nil {
		fail(t, message)
	}
}

func assertNil(t *testing.T, toTest interface{}, message string) {
	if toTest != nil {
		fail(t, message)
	}
}
func assertTrue(t *testing.T, toTest bool, message string) {
	if toTest == false {
		fail(t, message)
	}
}

func _assertEquals(t *testing.T, testable Testable, message string) {
	var b strings.Builder
	var isFail bool = false
	if !testable.IsEqual() {
		isFail = true
		b.WriteString(testable.String())
	}
	if len(message) > 0 {
		b.WriteString("\n" + message)
	}
	if isFail {
		fail(t, b.String())
	}
}
