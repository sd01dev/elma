package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	maxGoroutines = 5
	intendedWord  = "Go"
	golangURL     = "https://golang.org"
	goWiki        = "https://en.wikipedia.org/wiki/Go_(programming_language)"
)

func main() {
	ch := make(chan int, maxGoroutines)
	sum := 0
	urlList := []string{goWiki, golangURL, goWiki, golangURL, goWiki, golangURL}

	for _, value := range urlList {
		go countIntendedWordNumber(ch, value, intendedWord)
	}

	for i := len(urlList); i != 0; i-- {
		sum += <-ch
	}
	fmt.Printf("Total: %d", sum)
}

func wordCount(text, intendedWord string) int {
	words := strings.Fields(text)
	wordMap := make(map[string]int)
	for _, word := range words {
		wordMap[word] += 1
	}
	return wordMap[intendedWord]
}

func getData(url string) (string, error) {
	req, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("get URL error: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(req.Body)

	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", fmt.Errorf("read content error: %w", err)
	}
	return string(content), nil
}

// Counts only intended word covered with spaces. Example: " Go ".
// If needed count with special symbols before\after desired word, should use regexp.
func countIntendedWordNumber(ch chan int, url, intendedWord string) {
	data, err := getData(url)
	if err != nil {
		return
	}
	matches := wordCount(data, intendedWord)
	fmt.Printf("Count for %s: %v\n", url, matches)
	ch <- matches
}
