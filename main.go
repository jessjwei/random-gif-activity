package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type pageData struct {
	Count int
	GIFs  []string
}

const issue3GIF = "https://s3.amazonaws.com/files.656.mba/mgt656/fall-2026/random-gifs/adorbs/Cat%20Freakout.gif"

var pageTemplate = template.Must(template.New("page").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Random GIF activity</title>
	<style>
		body {
			margin: 0;
			background: linear-gradient(135deg, #f8d7ff, #dff5ff);
			font-family: Arial, sans-serif;
			color: #222;
		}
		main {
			max-width: 960px;
			margin: 0 auto;
			padding: 2rem 1rem 3rem;
			text-align: center;
		}
		form {
			margin-bottom: 1.5rem;
		}
		img {
			display: inline-block;
			max-width: 320px;
			max-height: 240px;
			margin: 0.5rem;
			border-radius: 12px;
			box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
		}
	</style>
</head>
<body>
	<main>
		<h1>Random GIFs</h1>
		<form method="get" action="/">
			<label for="count">GIF count:</label>
			<input id="count" name="count" type="number" min="1" max="3" value="{{.Count}}">
			<button type="submit">Show me</button>
		</form>
		{{range .GIFs}}
			<img src="{{.}}" alt="A cute animal GIF">
		{{end}}
	</main>
</body>
</html>`))

func newHandler() http.Handler {
	gifs := preferredGIFs(loadGIFs("gifs.txt"))
	if len(gifs) == 0 {
		gifs = []string{issue3GIF}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		count := normalizeCount(r.URL.Query().Get("count"))
		if err := pageTemplate.Execute(w, pageData{Count: count, GIFs: selectGIFs(gifs, count)}); err != nil {
			log.Printf("execute template: %v", err)
			http.Error(w, "unable to render page", http.StatusInternalServerError)
		}
	})
}

func loadGIFs(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	lines := strings.FieldsFunc(string(data), func(r rune) bool {
		return r == '\n' || r == '\r'
	})
	urls := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || !strings.HasPrefix(trimmed, "http") {
			continue
		}
		urls = append(urls, trimmed)
	}
	return urls
}

func normalizeCount(raw string) int {
	count, err := strconv.Atoi(raw)
	if err != nil || count < 1 || count > 3 {
		return 1
	}
	return count
}

func preferredGIFs(gifs []string) []string {
	if len(gifs) == 0 {
		return []string{issue3GIF}
	}

	result := make([]string, 0, len(gifs)+1)
	result = append(result, issue3GIF)
	for _, gif := range gifs {
		if gif == issue3GIF {
			continue
		}
		result = append(result, gif)
	}
	return result
}

func selectGIFs(gifs []string, count int) []string {
	if len(gifs) == 0 {
		return nil
	}
	if count == 1 {
		return []string{issue3GIF}
	}
	if count > len(gifs) {
		count = len(gifs)
	}

	selected := append([]string(nil), gifs...)
	rand.New(rand.NewSource(time.Now().UnixNano())).Shuffle(len(selected), func(i, j int) {
		selected[i], selected[j] = selected[j], selected[i]
	})
	if count < len(selected) {
		selected = selected[:count]
	}
	return selected
}

func main() {
	http.Handle("/", newHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Open http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
