package main

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"testing/fstest"

	goldenpathdemo "github.com/Solisol/goldenpath-demo"
)

func testLibrary(t *testing.T) *Library {
	t.Helper()
	lib, err := LoadLibrary(fstest.MapFS{
		"storage.md": {Data: []byte("---\nname: Storing data\nkeywords: [storage, database]\n---\nSee `goldenpath lookup buckets`.\n")},
		"buckets.md": {Data: []byte("---\nname: Object storage (buckets)\nkeywords: [bucket, object storage, images]\n---\nUse STORAGE_BUCKET.\n")},
	})
	if err != nil {
		t.Fatal(err)
	}
	return lib
}

func TestParseArticle(t *testing.T) {
	a, err := parseArticle("x", "---\nname: \"Start here: how we build\"\nkeywords: [Go, ' web app ', ]\nstatus: draft\n---\n\nBody line\n")
	if err != nil {
		t.Fatal(err)
	}
	if a.Name != "Start here: how we build" {
		t.Errorf("name = %q", a.Name)
	}
	if len(a.Keywords) != 2 || a.Keywords[0] != "go" || a.Keywords[1] != "web app" {
		t.Errorf("keywords = %q", a.Keywords)
	}
	if a.Content != "Body line\n" {
		t.Errorf("content = %q", a.Content)
	}
}

func TestParseArticleErrors(t *testing.T) {
	for name, raw := range map[string]string{
		"no frontmatter": "just text\n",
		"not closed":     "---\nname: x\n",
		"no name":        "---\nkeywords: [a]\n---\n",
	} {
		if _, err := parseArticle("x", raw); err == nil {
			t.Errorf("%s: expected error", name)
		}
	}
}

func TestSearch(t *testing.T) {
	lib := testLibrary(t)
	cases := map[string][]string{
		"storage":      {"storage", "buckets"}, // exact keyword ranks above partial + name
		"STOR":         {"buckets", "storage"}, // partial, case-insensitive
		"image":        {"buckets"},
		"data":         {"storage"}, // matches the name "Storing data" and keyword "database"
		"nothing here": {},
	}
	for q, want := range cases {
		got := lib.Search(q)
		if len(got) != len(want) {
			t.Errorf("Search(%q) = %d hits, want %v", q, len(got), want)
			continue
		}
		for i := range want {
			if got[i].ID != want[i] {
				t.Errorf("Search(%q)[%d] = %s, want %s", q, i, got[i].ID, want[i])
			}
		}
	}
}

func TestHTTP(t *testing.T) {
	act := NewActivity(10)
	srv := httptest.NewServer(newServer(testLibrary(t), act))
	defer srv.Close()

	get := func(path string, wantStatus int, v any) {
		t.Helper()
		resp, err := http.Get(srv.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != wantStatus {
			t.Fatalf("GET %s = %d, want %d", path, resp.StatusCode, wantStatus)
		}
		if v != nil {
			if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
				t.Fatal(err)
			}
		}
	}

	var s struct{ Results []Hit }
	get("/search?q=bucket", 200, &s)
	if len(s.Results) != 1 || s.Results[0].ID != "buckets" || s.Results[0].Name == "" {
		t.Errorf("search results = %+v", s.Results)
	}
	var a Article
	get("/lookup/buckets", 200, &a)
	if a.Content != "Use STORAGE_BUCKET.\n" {
		t.Errorf("lookup content = %q", a.Content)
	}
	get("/lookup/BUCKETS", 200, nil)
	get("/lookup/missing", 404, nil)
	get("/search", 400, nil)
	get("/healthz", 200, nil)

	if n := len(act.recent); n != 4 { // 1 search + 3 lookups; bad requests aren't recorded
		t.Errorf("recorded %d events, want 4", n)
	}
	if act.recent[0].Kind != "search" || act.recent[0].Results != 1 {
		t.Errorf("first event = %+v", act.recent[0])
	}
}

// TestShippedContent checks the real golden path: it loads, and every
// `goldenpath lookup <id>` reference points at an article that exists.
func TestShippedContent(t *testing.T) {
	content, _ := fs.Sub(goldenpathdemo.Articles, "articles")
	lib, err := LoadLibrary(content)
	if err != nil {
		t.Fatal(err)
	}
	if lib.Len() == 0 {
		t.Fatal("no articles")
	}
	ref := regexp.MustCompile("goldenpath lookup ([a-z0-9-]+)")
	for id, a := range lib.articles {
		if len(a.Keywords) == 0 {
			t.Errorf("%s: no keywords", id)
		}
		for _, m := range ref.FindAllStringSubmatch(a.Content, -1) {
			if _, ok := lib.Lookup(m[1]); !ok {
				t.Errorf("%s: references missing article %q", id, m[1])
			}
		}
	}
}
