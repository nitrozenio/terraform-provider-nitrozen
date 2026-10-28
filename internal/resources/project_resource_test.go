package resources

import "testing"

func TestParseProjectResponsePublicURL(t *testing.T) {
	body := []byte(`{"data":{"id":7,"name":"TextPro","slug":"textproxyz","description":"x","public_url":"https://nitrozen.io/c/textproxyz"}}`)
	got, err := parseProjectResponse(body)
	if err != nil {
		t.Fatal(err)
	}
	if got.Data.ID != 7 || got.Data.PublicURL != "https://nitrozen.io/c/textproxyz" {
		t.Fatalf("unexpected parse result: %+v", got.Data)
	}
}
