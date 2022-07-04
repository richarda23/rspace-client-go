package rspace

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGetForms(t *testing.T) {
	cfg := NewRecordListingConfig()
	got, err := webClient.Forms(cfg)
	if err != nil {
		t.Fatalf("Error listing forms")
	}
	assertNotNil(t, got, "forms listing should not be nil")
	assertTrue(t, len(got.Forms) > 0, "must be at least 1 form")
	for _, v := range got.Forms {
		assertStringEquals(t, fmt.Sprintf("FM%d", v.Id), v.GlobalId, "")
	}
}

func TestCreateForm(t *testing.T) {
	filePath := "testdata/form.yaml"
	yamlFormDef, _ := os.Open(filePath)
	form, _ := webClient.CreateFormYaml(yamlFormDef)
	assertStringEquals(t, "MyForm", form.Name, "")
	form, _ = webClient.PublishForm(form.Id)
	fmt.Println(form)
	assertStringEquals(t, "PUBLISHED", form.FormState, "")

}

func TestSearchForms(t *testing.T) {
	cfg := NewRecordListingConfig()
	got, err := webClient.Forms(cfg)
	if err != nil {
		t.Fatalf("Error listing forms")
	}
	Log.Info("searching for '" + got.Forms[0].Name + "'")
	// there must be at least 1 form, now search by its name
	hits, _ := webClient.FormSearch(cfg, lower(got.Forms[0].Name))
	assertTrue(t, hits.TotalHits > 0, "expected at least 1 search hit")
	for _, v := range hits.Forms {
		fmt.Println(v)
		assertTrue(t, strings.Contains(lower(v.Name), lower(got.Forms[0].Name)), "")
	}
}

func lower(arg string) string {
	return strings.ToLower(arg)
}
