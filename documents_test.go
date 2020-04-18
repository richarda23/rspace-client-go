package rspace

import (
	"fmt"
	"github.com/op/go-logging"
	"net/url"
	"os"
	"testing"
)

// Main entry point into testing, sets in env variables and creates web client
func TestMain(m *testing.M) {
	initLogging(logging.INFO)
	url, _ := url.Parse(getenv(BASE_URL_ENV_NAME))
	fmt.Println("url is " + url.String())
	apikey := getenv(APIKEY_ENV_NAME)
	fmt.Println("api is " + apikey)
	webClient = NewWebClient(url, apikey)
	os.Exit(m.Run())
}

func TestStatus(t *testing.T) {
	got, err := webClient.documentS.GetStatus()
	if err != nil {
		Log.Error(err)
	}
	Log.Info(Marshal(got))
	assertTrue(t, len(got.RSpaceVersion) > 0, "RSpaceVersion must be non-empty")
}

func TestDocumentList(t *testing.T) {
	cfg := NewRecordListingConfig()
	got, _ := webClient.documentS.Documents(cfg)
	if got.TotalHits <= 1 {
		fail(t, fmt.Sprintf("Expected hits >= 1 but was %d", got.TotalHits))
	}
}

func TestDocumentBasicSearch(t *testing.T) {
	name := randomAlphanumeric(6)
	tag := randomAlphanumeric(6)
	cfg := NewRecordListingConfig()
	created := webClient.documentS.NewEmptyBasicDocument(name, tag)
	results, _ := webClient.documentS.SearchDocuments(cfg, name)
	assertIntEquals(t, 1, results.TotalHits, "")
	assertIntEquals(t, created.Id, results.Documents[0].Id, "")
}
func TestDocumentAdvancedSearch(t *testing.T) {
	// given
	name := randomAlphanumeric(6)
	tag := randomAlphanumeric(6)
	tag2 := randomAlphanumeric(6)
	created := webClient.documentS.NewEmptyBasicDocument(name, tag)
	cfg := NewRecordListingConfig()

	builder := &SearchQueryBuilder{}
	builder.Operator(And).AddTerm(name, NAME).AddTerm(tag, TAG)
	query := builder.Build()
	//when
	results, _ := webClient.documentS.AdvancedSearchDocuments(cfg, query)
	//then
	assertIntEquals(t, 1, results.TotalHits, "")
	assertIntEquals(t, created.Id, results.Documents[0].Id, "")
	assertStringEquals(t, created.Name, results.Documents[0].Name, "")
	// and doesn't match here
	builder2 := &SearchQueryBuilder{}
	builder2.Operator(And).AddTerm(name, NAME).AddTerm(tag2, TAG)
	query2 := builder2.Build()
	results2, _ := webClient.documentS.AdvancedSearchDocuments(cfg, query2)
	assertIntEquals(t, 0, results2.TotalHits, "")
	// but or does
	builder3 := &SearchQueryBuilder{}
	builder3.Operator(Or).AddTerm(name, NAME).AddTerm(tag2, TAG)
	query3 := builder3.Build()
	results3, _ := webClient.documentS.AdvancedSearchDocuments(cfg, query3)
	assertIntEquals(t, 1, results3.TotalHits, "")

}

func TestDocumentNew(t *testing.T) {
	//post := DocumentPostNewBasicDocument("go12", "t1,t2,t3")
	var got = webClient.documentS.NewEmptyBasicDocument("go12", "tag1,tag2")
	Log.Info(Marshal(got))
	if got.Name != "go12" {
		fail(t, fmt.Sprintf("Expected 'go1' > 1 but was %s", got.Name))
	}
	if got.Tags != "tag1,tag2" {
		fail(t, fmt.Sprintf("Expected 'tag1,tags' > 1 but was %s", got.Tags))
	}
	var got2 = webClient.documentS.NewEmptyBasicDocument("nameonly", "")
	if got2.Tags != "" {
		fail(t, fmt.Sprintf("Expected '' > 1 but was %s", got2.Tags))
	}
	var got3 = webClient.documentS.NewBasicDocumentWithContent("n1", "t1", "<p> Some content </p")
	if got3 == nil {
		fail(t, "Doc3 is nil")
	}
	// now delete
	rs, _ := webClient.documentS.DeleteDocument(got3.Id)
	assertTrue(t, rs, "Delete document failed:")
}
