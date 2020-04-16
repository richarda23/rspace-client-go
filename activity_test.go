package rspace

import (
	"fmt"
	"testing"
	"time"
)

func TestActivityGet(t *testing.T) {
	var builder ActivityQueryBuilder = ActivityQueryBuilder{}
	var err error
	var result *ActivityList
	var q *ActivityQuery
	builder.Domain("RECORD")
	q, _ = builder.Build()
	result, err = webClient.ActivityS.Activities(q)
	if err != nil {
		fmt.Println(err)
	}

	//get non-existent results
	builder.DateFrom(time.Now().AddDate(1, 0, 0))
	q, _ = builder.Build()
	result, err = webClient.ActivityS.Activities(q)
	if err != nil {
		fmt.Println(err)	
	}
	assertIntEquals(t, 0, result.TotalHits, "")
	// too far in the past
	builder = ActivityQueryBuilder{}
	builder.DateTo(time.Now().AddDate(-10, 0, 0))
	q, _ = builder.Build()
	result, err = webClient.ActivityS.Activities(q)
	assertIntEquals(t, 0, result.TotalHits, "")
}
func TestActivityForDocumentGet(t *testing.T) {
	name := randomAlphanumeric(6)
	created := webClient.NewEmptyBasicDocument(name, "")
	builder := ActivityQueryBuilder{}
	q, _ := builder.Oid(GlobalId(created.GlobalId)).Build()
	result, err := webClient.ActivityS.Activities(q)
	assertNil(t, err, "error should be nil")
	assertIntEquals(t, 1, result.TotalHits, "")
	assertStringEquals(t, "CREATE", result.Activities[0].Action, "")
	assertStringEquals(t, "RECORD", result.Activities[0].Domain, "")
	timestamp, _ := result.Activities[0].TimestampTime()
	assertTrue(t, timestamp.Before(time.Now()), "timestamp parsing is invalid")
	assertIntEquals(t, 1, len(result.Links), "")
	assertStringEquals(t, "/api/v1/activity", result.Links[0].Link.Path, "")
}
