package rspace

import (
	"testing"
)

func TestStatus(t *testing.T) {
	got := GetStatus()
	if len(got.RSpaceVersion) == 0 {
		t.Errorf("RSpaceVersion must be non-empty")
	}
}

func TestDocumentList(t *testing.T) {
	cfg := NewRecordListingConfig()
	got := Documents(cfg)
	if got.TotalHits <= 1 {
		t.Errorf("Expected hits > 1 but was %d", got.TotalHits)
	}
}

func TestDocumentNewList(t *testing.T) {
	//post := DocumentPostNewBasicDocument("go12", "t1,t2,t3")
        var  got  = NewEmptyBasicDocument("go12", "tag1,tag2")
	if got.Name != "go12" {
		t.Errorf("Expected 'go1' > 1 but was %s", got.Name)
	}
}
