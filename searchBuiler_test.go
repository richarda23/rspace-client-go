package rspace

import (
	"fmt"
	"testing"
)

func TestSearchBuilder(t *testing.T) {
	builder := &SearchQueryBuilder{}
	builder.operator(OR).addTerm("tag1", TAG).addTerm("formName", FORM)
	fmt.Println(builder)
	query := builder.build()
	fmt.Println("query is " + query.String())
	if len(query.Terms) != 2  {
		fail(t, fmt.Sprintf("should have 2 terms but was %d", len(query.Terms)))
	}
}

