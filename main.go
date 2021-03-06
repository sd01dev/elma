package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	workerLimit  = 5
	intendedWord = "Go"
	golangURL    = "https://golang.org"
	goWiki       = "https://en.wikipedia.org/wiki/Go_(programming_language)"
	githubGo     = "https://github.com/golang/go"
	habrGo       = "https://habr.com/ru/hub/go/"
	freeCodeCamp = "https://www.freecodecamp.org/news/what-is-go-programming-language/"
)

type site struct {
	url string
}

type matches struct {
	amount int
}

func main() {
	jobs := make(chan site, workerLimit)
	results := make(chan matches, workerLimit)
	urlList := []string{goWiki, golangURL, githubGo, habrGo, freeCodeCamp, githubGo}
	sum := 0

	for w := 1; w <= workerLimit; w++ {
		go countIntendedWordNumber(jobs, results, intendedWord)
	}

	for _, currentUrl := range urlList {
		jobs <- site{url: currentUrl}
	}
	close(jobs)

	for i := 1; i <= len(urlList); i++ {
		amount := <-results
		sum += amount.amount
	}
	log.Printf("Total: %d", sum)
}

func wordCount(text, intendedWord string) (counter int) {
	words := strings.Fields(text)
	for _, word := range words {
		if word == intendedWord {
			counter += 1
		}

	}
	return counter
}

func getData(url string) (string, error) {
	req, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("get url error: %w", err)
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
func countIntendedWordNumber(jobs <-chan site, results chan<- matches, intendedWord string) {
	for site := range jobs {
		data, err := getData(site.url)
		if err != nil {
			return
		}
		repeatedWord := wordCount(data, intendedWord)
		log.Printf("Count for %s: %v\n", site.url, repeatedWord)
		results <- matches{amount: repeatedWord}
	}
}
