package rspace

import (
	"fmt"
	"testing"
)

func TestGetForms(t *testing.T) {
	cfg := NewRecordListingConfig()
	got, err := webClient.FormS.Forms(cfg)
	fmt.Println(err)
	assertNotNil(t, got, "forms listing should not be nil")
	assertTrue(t, len(got.Forms) > 0, "must be at least 1 form")
	for _, v := range got.Forms {
		assertStringEquals(t, fmt.Sprintf("FM%d", v.Id), v.GlobalId, "")
	}
}
