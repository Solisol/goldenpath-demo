package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Event is one request, as shown on the activity page.
type Event struct {
	Time    time.Time `json:"time"`
	Kind    string    `json:"kind"` // search | lookup
	Query   string    `json:"query"`
	Results int       `json:"results"`
	Found   bool      `json:"found"`
	Client  string    `json:"client,omitempty"`
}

// Activity keeps the most recent events and streams new ones to the activity page.
// For the platform team this is product telemetry: what do developers and agents
// actually ask the golden path for, and what do they not find?
type Activity struct {
	mu     sync.Mutex
	recent []Event
	max    int
	subs   map[chan Event]struct{}
}

func NewActivity(max int) *Activity {
	return &Activity{max: max, subs: map[chan Event]struct{}{}}
}

func (a *Activity) Record(e Event) {
	e.Time = time.Now()
	if e.Kind == "search" {
		log.Printf("search %q results=%d client=%q", e.Query, e.Results, e.Client)
	} else {
		log.Printf("lookup %q found=%t client=%q", e.Query, e.Found, e.Client)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.recent = append(a.recent, e)
	if len(a.recent) > a.max {
		a.recent = a.recent[len(a.recent)-a.max:]
	}
	for ch := range a.subs {
		select {
		case ch <- e:
		default: // slow viewer: drop rather than block requests
		}
	}
}

// ServeEvents streams events as server-sent events: the recent history first, then live.
func (a *Activity) ServeEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "text/event-stream")
	w.Header().Set("cache-control", "no-cache")

	ch := make(chan Event, 32)
	a.mu.Lock()
	history := append([]Event(nil), a.recent...)
	a.subs[ch] = struct{}{}
	a.mu.Unlock()
	defer func() {
		a.mu.Lock()
		delete(a.subs, ch)
		a.mu.Unlock()
	}()

	send := func(e Event) {
		b, _ := json.Marshal(e)
		fmt.Fprintf(w, "data: %s\n\n", b)
	}
	for _, e := range history {
		send(e)
	}
	flusher.Flush()

	keepalive := time.NewTicker(15 * time.Second)
	defer keepalive.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case e := <-ch:
			send(e)
			flusher.Flush()
		case <-keepalive.C:
			fmt.Fprint(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}

const activityPage = `<!doctype html>
<html lang="en"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Golden path activity</title>
<style>
  :root { --bg:#101216; --fg:#ececec; --dim:#8b919b; --search:#7cc4ff; --lookup:#b7f0a0; --miss:#ff9a8a; --line:#23262d; }
  body { background:var(--bg); color:var(--fg); font:26px/1.4 ui-monospace, "SF Mono", Menlo, monospace; margin:0; padding:28px 32px; }
  header { display:flex; justify-content:space-between; align-items:baseline; margin-bottom:20px; }
  h1 { font-size:20px; color:var(--dim); font-weight:normal; letter-spacing:.08em; text-transform:uppercase; margin:0; }
  #count { color:var(--dim); font-size:20px; }
  ol { list-style:none; margin:0; padding:0; }
  li { display:grid; grid-template-columns:5.5em 4.5em 1fr auto; gap:0 18px; padding:10px 0; border-bottom:1px solid var(--line); animation:in .35s ease-out; }
  .t, .r { color:var(--dim); } .r { font-size:20px; align-self:center; }
  .search .k { color:var(--search); } .lookup .k { color:var(--lookup); } .miss .k, .miss .q { color:var(--miss); }
  .q { overflow-wrap:anywhere; }
  #empty { color:var(--dim); padding:24px 0; }
  @keyframes in { from { opacity:0; transform:translateY(-6px); } }
</style></head>
<body>
<header><h1>Golden path activity</h1><span id="count"></span></header>
<p id="empty">Waiting for the first search or lookup…</p>
<ol id="log"></ol>
<script>
  const log = document.getElementById("log"), count = document.getElementById("count");
  let n = 0;
  const es = new EventSource("/activity");
  es.onmessage = (msg) => {
    const e = JSON.parse(msg.data);
    document.getElementById("empty").hidden = true;
    const miss = (e.kind === "search" && e.results === 0) || (e.kind === "lookup" && !e.found);
    const li = document.createElement("li");
    li.className = e.kind + (miss ? " miss" : "");
    const time = new Date(e.time).toLocaleTimeString([], { hour12: false });
    const result = e.kind === "search" ? e.results + " hit" + (e.results === 1 ? "" : "s") : (e.found ? "found" : "not found");
    for (const [cls, text] of [["t", time], ["k", e.kind], ["q", e.query], ["r", result]]) {
      const span = document.createElement("span");
      span.className = cls; span.textContent = text; li.appendChild(span);
    }
    log.prepend(li);
    count.textContent = ++n + " requests";
  };
</script>
</body></html>`
