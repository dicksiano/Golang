package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

var safeVisited = struct {
	v   map[string]bool
	sync.Mutex
} { v: make(map[string]bool) }

func Crawl(url string, depth int, fetcher Fetcher) {
	MyCrawl(url, depth, fetcher, depth)
}

func prettyPrinter(diff int) {
	for i := 0; i < diff; i++ {
		fmt.Print("		")
	}
}

func MyCrawl(url string, depth int, fetcher Fetcher, maxDepth int) {
	if depth <= 0 {
		fmt.Println("Depth 0!")
		return
	}
	
	safeVisited.Lock()

	if safeVisited.v[url] {
		safeVisited.Unlock()
		prettyPrinter(maxDepth - depth)
		fmt.Println("Done: ", url, " Already fetched!")
		return
	}

	safeVisited.v[url] = true	
	safeVisited.Unlock()

	body, urls, err := fetcher.Fetch(url)

	if err !=  nil {
		prettyPrinter(maxDepth - depth)
		fmt.Println("Url: ", url, " -  Error: ", err)
		return
	}

	fmt.Println()
	prettyPrinter(maxDepth - depth)
	fmt.Println("Url: ", url, "       Body: ", body, "\n")

	done := make(chan bool)

	for i,u := range urls {
		prettyPrinter(maxDepth - depth)
		fmt.Println("Crawling ", i, "/", len(urls), " of ", url, " : ", u)

		go func( url string) {
			MyCrawl(url, depth-1, fetcher, maxDepth)
			done <- true
		} (u)
	}		
	
	for u := range urls {
		prettyPrinter(maxDepth - depth)
		fmt.Println("Parent: ", url, " ---  waiting for child: ", u)
		<- done
	}
	prettyPrinter(maxDepth - depth)
	fmt.Println("Url: ", url, " finished!")
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f *fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := (*f)[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = &fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

