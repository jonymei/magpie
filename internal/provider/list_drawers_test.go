package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

// The add form's "Fetch models" (List) must offer the vendor's image models
// too: a relay that serves gpt-image (sub2api does) is picked from here, and
// Settings → Images is what draws with it.
func TestListOffersImageModels(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"object":"list","data":[{"id":"gpt-5.6-sol"},{"id":"gpt-image-2.5-flare"}]}`))
	}))
	defer server.Close()

	p := Provider{ID: "relay", Name: "Relay", Responses: server.URL + "/v1", Key: "k"}
	ms, err := p.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, m := range ms {
		ids = append(ids, m.ID)
	}
	if !slices.Contains(ids, "gpt-image-2.5-flare") {
		t.Fatalf("List = %v, want gpt-image-2.5-flare among them", ids)
	}
}