// goldenpath-api serves the organization's golden path: search and lookup over
// Markdown articles, plus a live activity page showing every request.
package main

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	goldenpathdemo "github.com/Solisol/goldenpath-demo"
)

func main() {
	content, err := fs.Sub(goldenpathdemo.Articles, "articles")
	if err != nil {
		log.Fatal(err)
	}
	// GOLDENPATH_CONTENT_DIR serves articles from disk instead, e.g. a mounted wiki
	// export or a git checkout, so content can change without rebuilding the image.
	if dir := os.Getenv("GOLDENPATH_CONTENT_DIR"); dir != "" {
		content = os.DirFS(dir)
	}
	lib, err := LoadLibrary(content)
	if err != nil {
		log.Fatalf("loading golden path: %v", err)
	}

	addr := ":" + envOr("PORT", "8080")
	log.Printf("goldenpath-api: %d articles, listening on %s", lib.Len(), addr)
	log.Fatal(http.ListenAndServe(addr, newServer(lib, NewActivity(200))))
}

func newServer(lib *Library, act *Activity) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /search", func(w http.ResponseWriter, r *http.Request) {
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		if q == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing query parameter q"})
			return
		}
		hits := lib.Search(q)
		act.Record(Event{Kind: "search", Query: q, Results: len(hits), Client: client(r)})
		writeJSON(w, http.StatusOK, map[string]any{"query": q, "results": hits})
	})

	mux.HandleFunc("GET /lookup/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		a, ok := lib.Lookup(id)
		act.Record(Event{Kind: "lookup", Query: id, Found: ok, Client: client(r)})
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "no article " + id, "hint": "use /search?q=<words> to find article ids"})
			return
		}
		writeJSON(w, http.StatusOK, a)
	})

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("GET /activity", act.ServeEvents)
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html; charset=utf-8")
		w.Write([]byte(activityPage))
	})
	return mux
}

// client identifies the caller, e.g. "goldenpath/0.1.0". Knowing which client
// versions are out there is half the reason the CLI is kept thin.
func client(r *http.Request) string {
	if c := r.Header.Get("X-Goldenpath-Client"); c != "" {
		return c
	}
	return r.UserAgent()
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
