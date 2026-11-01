// goldenpath is a deliberately thin client for the golden-path service.
//
// Everything that may change (content, search, guidelines) lives in the service,
// which the platform team controls. This binary only forwards the question and
// prints the answer, so it rarely needs updating once it's on a laptop.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var version = "dev" // set at build time: -ldflags "-X main.version=0.1.0"

const usage = `goldenpath: how we build things here.

Usage:
  goldenpath search <words>   find golden-path articles by keyword or name
  goldenpath lookup <id>      show one article
  goldenpath version

Flags:
  --json        print the raw JSON response
  --url <url>   service URL (default $GOLDENPATH_URL or http://localhost:8080)
`

func main() {
	fs := flag.NewFlagSet("goldenpath", flag.ExitOnError)
	fs.Usage = func() { fmt.Fprint(os.Stderr, usage) }
	asJSON := fs.Bool("json", false, "")
	base := fs.String("url", envOr("GOLDENPATH_URL", "http://localhost:8080"), "")
	fs.Parse(reorder(os.Args[1:]))
	args := fs.Args()

	if len(args) == 0 {
		fs.Usage()
		os.Exit(2)
	}
	var path string
	switch args[0] {
	case "search":
		if len(args) < 2 {
			fail("usage: goldenpath search <words>")
		}
		path = "/search?q=" + url.QueryEscape(strings.Join(args[1:], " "))
	case "lookup":
		if len(args) != 2 {
			fail("usage: goldenpath lookup <id>")
		}
		path = "/lookup/" + url.PathEscape(args[1])
	case "version":
		fmt.Println("goldenpath", version)
		return
	default:
		fs.Usage()
		os.Exit(2)
	}

	body, status := get(strings.TrimRight(*base, "/") + path)
	if *asJSON {
		os.Stdout.Write(body)
		if status >= 400 {
			os.Exit(1)
		}
		return
	}
	if status == http.StatusNotFound && args[0] == "lookup" {
		fail(fmt.Sprintf("No golden-path article %q. Find article ids with: goldenpath search <words>", args[1]))
	}
	if status >= 400 {
		var e struct{ Error string }
		json.Unmarshal(body, &e)
		fail(fmt.Sprintf("golden-path service error (%d): %s", status, e.Error))
	}
	if args[0] == "search" {
		printSearch(body)
	} else {
		printArticle(body)
	}
}

func get(u string) ([]byte, int) {
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("X-Goldenpath-Client", "goldenpath/"+version)
	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		fail(fmt.Sprintf("cannot reach the golden-path service: %v\n(set GOLDENPATH_URL or --url)", err))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body, resp.StatusCode
}

func printSearch(body []byte) {
	var r struct {
		Query   string
		Results []struct {
			ID, Name string
			Matched  []string
		}
	}
	if err := json.Unmarshal(body, &r); err != nil {
		fail("unexpected response from service: " + err.Error())
	}
	if len(r.Results) == 0 {
		fmt.Printf("No golden-path articles match %q. Try other words.\n", r.Query)
		return
	}
	width := 0
	for _, h := range r.Results {
		width = max(width, len(h.ID))
	}
	for _, h := range r.Results {
		fmt.Printf("%-*s  %s  (matched: %s)\n", width, h.ID, h.Name, strings.Join(h.Matched, ", "))
	}
	fmt.Println("\nRead one with: goldenpath lookup <id>")
}

func printArticle(body []byte) {
	var a struct {
		ID, Name, Content string
		Keywords          []string
	}
	if err := json.Unmarshal(body, &a); err != nil {
		fail("unexpected response from service: " + err.Error())
	}
	fmt.Printf("# %s\n(id: %s · keywords: %s)\n\n%s", a.Name, a.ID, strings.Join(a.Keywords, ", "), a.Content)
}

// reorder moves flags before positional arguments, so `goldenpath search storage --json` works.
func reorder(args []string) []string {
	var flags, rest []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--json" || a == "-json":
			flags = append(flags, a)
		case a == "--url" || a == "-url":
			if i+1 < len(args) {
				flags = append(flags, a, args[i+1])
				i++
			}
		case strings.HasPrefix(a, "--url=") || strings.HasPrefix(a, "-url="):
			flags = append(flags, a)
		default:
			rest = append(rest, a)
		}
	}
	return append(flags, rest...)
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
