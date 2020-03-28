package rspace

import (
	"fmt"
	"testing"
	"os"
	"github.com/op/go-logging"
)

func TestMain(m *testing.M) {
	initLogging(logging.INFO)
	os.Exit(m.Run())
}

func fail(t *testing.T, message string) {
	t.Errorf(message)
}

func TestStatus(t *testing.T) {
	got := GetStatus()
	Log.Info(Marshal(got))
	if len(got.RSpaceVersion) == 0 {
		fail(t, "RSpaceVersion must be non-empty")
	}
	
}

func TestDocumentList(t *testing.T) {
	cfg := NewRecordListingConfig()
	got := Documents(cfg)
	Log.Info(Marshal(got))
	if got.TotalHits <= 1 {
		fail(t, fmt.Sprintf("Expected hits >= 1 but was %d", got.TotalHits))
	}
}

func TestDocumentNew(t *testing.T) {
	//post := DocumentPostNewBasicDocument("go12", "t1,t2,t3")
	var got = NewEmptyBasicDocument("go12", "tag1,tag2")
	Log.Info(Marshal(got))
	if got.Name != "go12" {
		fail(t, fmt.Sprintf("Expected 'go1' > 1 but was %s", got.Name))
	}
	if got.Tags != "tag1,tag2" {
		fail(t, fmt.Sprintf("Expected 'tag1,tags' > 1 but was %s", got.Tags))
	}
	var got2 = NewEmptyBasicDocument("nameonly", "")
	if got2.Tags != "" {
		fail(t, fmt.Sprintf("Expected '' > 1 but was %s", got2.Tags))
	}
	var got3 = NewBasicDocumentWithContent("n1", "t1", "<p> Some content </p")
	if got3 == nil {
		fail(t, "Doc3 is nil")
	}
	fullDoc := DocumentById(got3.Id)
	Log.Info(Marshal(fullDoc.Fields))

}
