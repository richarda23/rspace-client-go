package rspace

import (
	"fmt"
	"testing"
)
var formService = webClient.formService()

func TestGetForms(t *testing.T) {
	cfg := NewRecordListingConfig()
	got, _ := formService.Forms(cfg)
	assertNotNil(t, got, "forms listing should not be nil")
	for _, v := range got.Forms {
		assertStringEquals(t, fmt.Sprintf("FM%d", v.Id), v.GlobalId, "")
	}
}
