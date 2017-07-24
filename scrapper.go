package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"os"
	"log"
	"strings"
)

func getHost(url string) (host string) {
	return "www.sapo.pt"
}

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}

//func processUrl(url string)(protocol string, host string, path string){

//}

// Extract all http** links from a given webpage
func crawl(iurl string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(iurl)
	//	protocol, host, path := processUrl(url)

	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("ERROR: Failed to crawl \"" + iurl + "\"")
		return
	}
	u, err := url.Parse(iurl)
	if err != nil {
		fmt.Println("ERROR: Parsing url" + iurl)
		log.Fatal(err)
	}

	b := resp.Body
	defer b.Close() // close Body when the function returns

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			ok, urln := getHref(t)
			if !ok {
				continue
			}

			// Make sure the url begines in http**
			if strings.Index(urln, "mailto") == 0 {
				// @implement
			} else if strings.Index(urln, "http") == 0 {
				ch <- urln
			} else if strings.Index(urln, "/") == 0 {
				ch <- u.Scheme + "://" + u.Host + urln
			} else if strings.Index(urln, "#") == 0 {
				// @implement ?
			} else {
				ch <- urln
			}
		}
	}
}

func main2() {
	url2 := "http://sapo.pt"
	u, err := url.Parse(url2)
	if err != nil {
		//log.Fatal(err)
	}

	fmt.Println(u.Host)

	u.Scheme = "https"
	u.Host = "google.com"
	q := u.Query()
	q.Set("q", "golang")
	u.RawQuery = q.Encode()
	fmt.Println(u)

}

func main() {
	foundUrls := make(map[string]bool)
	seedUrls := os.Args[1:]

	// Channels
	chUrls := make(chan string)
	chFinished := make(chan bool)

	// Kick off the crawl process (concurrently)
	for _, iurl := range seedUrls {
		go crawl(iurl, chUrls, chFinished)
	}

	// Subscribe to both channels
	for c := 0; c < len(seedUrls); {
		select {
		case iurl := <-chUrls:
			foundUrls[iurl] = true
		case <-chFinished:
			c++
		}
	}

	// We're done! Print the results...

	fmt.Println("\nFound", len(foundUrls), "unique urls:\n")

	for iurl, _ := range foundUrls {
		fmt.Println(" - " + iurl)
	}

	close(chUrls)
}
