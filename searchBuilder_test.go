package rspace

import (
	"fmt"
	"testing"
	"encoding/json"
)


func TestSearchBuilder(t *testing.T) {
	builder := &SearchQueryBuilder{}
	builder.operator(or).addTerm("tag1", TAG).addTerm("formName", FORM)
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

