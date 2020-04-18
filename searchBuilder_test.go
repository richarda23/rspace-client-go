package rspace

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSearchBuilder(t *testing.T) {
	builder := &SearchQueryBuilder{}
	builder.Operator(Or).AddTerm("tag1", TAG).AddTerm("formName", FORM)
	query := builder.Build()
	fmt.Println("query is " + query.String())
	assertIntEquals(t, 2, len(query.Terms), "")
	json, _ := json.Marshal(query)
	fmt.Println(string(json))
}
func TestGlobalSearchBuilder(t *testing.T) {
	builder := &SearchQueryBuilder{}
	builder.AddGlobalTerm("anything")
	query := builder.Build()
	fmt.Println("global query is " + query.String())
	if query.Terms[0].QueryType != "global" {
		fail(t, fmt.Sprintf("search  should be global"))
	}
}
