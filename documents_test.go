package rspace

import (
	"fmt"
	"github.com/op/go-logging"
	"os"
	"testing"
	"time"
)

var ds *DocumentService = &DocumentService{
	BaseService: BaseService{
		Delay: time.Duration(100) * time.Millisecond}}

func TestMain(m *testing.M) {
	initLogging(logging.INFO)
	os.Exit(m.Run())
}


func TestStatus(t *testing.T) {
	got, err := ds.GetStatus()
	if err != nil {
		Log.Error(err)
	}
	Log.Info(Marshal(got))
	if len(got.RSpaceVersion) == 0 {
		fail(t, "RSpaceVersion must be non-empty")
	}
}

func TestDocumentList(t *testing.T) {
	cfg := NewRecordListingConfig()
	got, _ := ds.Documents(cfg)
	Log.Info(Marshal(got))
	if got.TotalHits <= 1 {
		fail(t, fmt.Sprintf("Expected hits >= 1 but was %d", got.TotalHits))
	}
}

func TestDocumentBasicSearch(t *testing.T) {
	name := randomAlphanumeric(6)
	tag := randomAlphanumeric(6)
	cfg := NewRecordListingConfig()
	created := ds.NewEmptyBasicDocument(name, tag)
	results, _ := ds.SearchDocuments(cfg, name)
	assertIntEquals(t, 1, results.TotalHits, "")
	assertIntEquals(t, created.Id, results.Documents[0].Id, "")
}
func TestDocumentAdvancedSearch(t *testing.T) {
	// given
	name := randomAlphanumeric(6)
	tag := randomAlphanumeric(6)
	tag2 := randomAlphanumeric(6)
	created := ds.NewEmptyBasicDocument(name, tag)
	cfg := NewRecordListingConfig()

	builder := &SearchQueryBuilder{}
	builder.operator(and).addTerm(name, NAME).addTerm(tag, TAG)
	query := builder.build()
	//when
	results, _ := ds.AdvancedSearchDocuments(cfg, query)
	//then
	assertIntEquals(t, 1, results.TotalHits, "")
	assertIntEquals(t, created.Id, results.Documents[0].Id, "")
	assertStringEquals(t, created.Name, results.Documents[0].Name, "")
	// and doesn't match here
	builder2 := &SearchQueryBuilder{}
	builder2.operator(and).addTerm(name, NAME).addTerm(tag2, TAG)
	query2 := builder2.build()
	results2, _ := ds.AdvancedSearchDocuments(cfg, query2)
	assertIntEquals(t, 0, results2.TotalHits, "")
	// but or does
	builder3 := &SearchQueryBuilder{}
	builder3.operator(or).addTerm(name, NAME).addTerm(tag2, TAG)
	query3 := builder3.build()
	results3, _ := ds.AdvancedSearchDocuments(cfg, query3)
	assertIntEquals(t, 1, results3.TotalHits, "")

}

func TestDocumentNew(t *testing.T) {
	//post := DocumentPostNewBasicDocument("go12", "t1,t2,t3")
	var got = ds.NewEmptyBasicDocument("go12", "tag1,tag2")
	Log.Info(Marshal(got))
	if got.Name != "go12" {
		fail(t, fmt.Sprintf("Expected 'go1' > 1 but was %s", got.Name))
	}
	if got.Tags != "tag1,tag2" {
		fail(t, fmt.Sprintf("Expected 'tag1,tags' > 1 but was %s", got.Tags))
	}
	var got2 = ds.NewEmptyBasicDocument("nameonly", "")
	if got2.Tags != "" {
		fail(t, fmt.Sprintf("Expected '' > 1 but was %s", got2.Tags))
	}
	var got3 = ds.NewBasicDocumentWithContent("n1", "t1", "<p> Some content </p")
	if got3 == nil {
		fail(t, "Doc3 is nil")
	}
	fullDoc, _ := ds.DocumentById(got3.Id)
	Log.Info(Marshal(fullDoc.Fields))

}
