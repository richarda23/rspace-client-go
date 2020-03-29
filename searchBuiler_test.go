package rspace

import (
	"fmt"
	"testing"
	"encoding/json"
	"strings"
)

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

func TestSearchBuilder(t *testing.T) {
	builder := &SearchQueryBuilder{}
	builder.operator(OR).addTerm("tag1", TAG).addTerm("formName", FORM)
	fmt.Println(builder)
	query := builder.build()
	fmt.Println("query is " + query.String())
	assertIntEquals(t, 2, len(query.Terms),"")
	json, _ := json.Marshal(query)
	fmt.Println(string(json))
}
func TestGlobalSearchBuilder(t *testing.T) {
	builder := &SearchQueryBuilder{}
	builder.addGlobalTerm("anything")
	query := builder.build()
	fmt.Println("global query is " + query.String())
	if query.Terms[0].QueryType != "global" {
		fail(t, fmt.Sprintf("search  should be global"))
	}

}

