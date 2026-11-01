package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

// Article is one golden-path page: a Markdown file with a small frontmatter block.
//
//	---
//	name: Object storage (buckets)
//	keywords: [storage, bucket, files]
//	---
//	Markdown body...
type Article struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Keywords []string `json:"keywords"`
	Content  string   `json:"content,omitempty"`
}

// Hit is a search result: the article without its content, plus what matched.
type Hit struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Keywords []string `json:"keywords"`
	Matched  []string `json:"matched"`
}

type Library struct {
	articles map[string]Article
}

// LoadLibrary reads every *.md file in fsys. The article id is the file name without .md.
func LoadLibrary(fsys fs.FS) (*Library, error) {
	lib := &Library{articles: map[string]Article{}}
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || path.Ext(p) != ".md" {
			return err
		}
		raw, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		a, err := parseArticle(strings.TrimSuffix(path.Base(p), ".md"), string(raw))
		if err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		if _, dup := lib.articles[a.ID]; dup {
			return fmt.Errorf("%s: duplicate article id %q", p, a.ID)
		}
		lib.articles[a.ID] = a
		return nil
	})
	return lib, err
}

func parseArticle(id, raw string) (Article, error) {
	a := Article{ID: id}
	sc := bufio.NewScanner(strings.NewReader(raw))
	if !sc.Scan() || strings.TrimSpace(sc.Text()) != "---" {
		return a, fmt.Errorf("missing frontmatter (first line must be ---)")
	}
	closed := false
	for sc.Scan() {
		line := sc.Text()
		if strings.TrimSpace(line) == "---" {
			closed = true
			break
		}
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		val = strings.TrimSpace(val)
		switch strings.TrimSpace(key) {
		case "name":
			a.Name = strings.Trim(val, `"'`)
		case "keywords":
			for _, k := range strings.Split(strings.Trim(val, "[]"), ",") {
				if k = strings.ToLower(strings.TrimSpace(strings.Trim(strings.TrimSpace(k), `"'`))); k != "" {
					a.Keywords = append(a.Keywords, k)
				}
			}
		}
	}
	if !closed {
		return a, fmt.Errorf("frontmatter not closed with ---")
	}
	if a.Name == "" {
		return a, fmt.Errorf("frontmatter has no name")
	}
	var body strings.Builder
	for sc.Scan() {
		body.WriteString(sc.Text())
		body.WriteByte('\n')
	}
	a.Content = strings.TrimLeft(body.String(), "\n")
	return a, sc.Err()
}

// Search does case-insensitive partial matching of each query word against an
// article's keywords and name. An exact keyword match scores highest, then a match
// in the name, then a partial keyword match; best score first.
func (l *Library) Search(query string) []Hit {
	words := strings.Fields(strings.ToLower(query))
	type scored struct {
		hit   Hit
		score int
	}
	var found []scored
	for _, a := range l.articles {
		var matched []string
		score := 0
		for _, w := range words {
			for _, k := range a.Keywords {
				switch {
				case k == w:
					matched, score = append(matched, k), score+4
				case strings.Contains(k, w):
					matched, score = append(matched, k), score+1
				}
			}
			if strings.Contains(strings.ToLower(a.Name), w) {
				matched, score = append(matched, "name"), score+2
			}
		}
		if score > 0 {
			found = append(found, scored{Hit{ID: a.ID, Name: a.Name, Keywords: a.Keywords, Matched: dedupe(matched)}, score})
		}
	}
	sort.Slice(found, func(i, j int) bool {
		if found[i].score != found[j].score {
			return found[i].score > found[j].score
		}
		return found[i].hit.ID < found[j].hit.ID
	})
	hits := make([]Hit, len(found))
	for i, f := range found {
		hits[i] = f.hit
	}
	return hits
}

func (l *Library) Lookup(id string) (Article, bool) {
	a, ok := l.articles[strings.ToLower(id)]
	return a, ok
}

func (l *Library) Len() int { return len(l.articles) }

func dedupe(s []string) []string {
	seen := map[string]bool{}
	out := s[:0]
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
