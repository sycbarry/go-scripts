package main

import (
    "io"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"crawler"
)

    /*
    1. GET the webpage
    2. parse all thelinks on teh package. (package we made before)
    3. build proper urls with our links.
    4. filter out any links w/ a different domain (non domain specific)
    5. find all the pages. (BFS)
    6. print out xml.
    */

func main() {
    //setup flags.
    urlFlag := flag.String("w", "https://google.com", "url we want to map")
    maxDepth := flag.Int("d", 3, "the depth of search we need")
    flag.Parse()

    pages := bfs(*urlFlag, *maxDepth)
    for _, page := range(pages) {
        fmt.Println(page)
    }


}



func bfs(url string, base int) []string {
    //the seen map will store the url's we've already seen. typical bfs
    //we use an empty struct here to designate no value and only a key.
    //this saves memory.
    seen :=make(map[string]struct{})
    var q map[string]struct{}
    nq := map[string]struct{}{
        url: struct{}{}, //struct{}{} means an empty struct.
    }
    for i := 0; i <= base; i++ {
        q, nq = nq, make(map[string]struct{})
        for curUrl, _ := range q {
            //if we've seen a link already during a previous pass.
            if _, ok := seen[curUrl]; ok {
                continue
            }
            seen[curUrl] = struct{}{} //mark it as seen.
            for _, link := range get(curUrl) {
                nq[link] = struct{}{}
            }
        }
    }
    ret := make([]string, 0, len(seen))
    for url, _ := range seen {
        ret = append(ret, url)
    }
    return ret
}



func filter(keepFn func(string) bool, links []string) []string {
    var ret []string
    //filter here
    for _, link := range links {
        if keepFn(link) {
            ret = append(ret, link)
        }
    }
    return ret
}


//closure used above
func withPrefix(pfx string) func(string) bool {
    return func(link string) bool {
        return strings.HasPrefix(link, pfx)
    }
}


func get(urlStr string) []string {


    r, err := http.Get(urlStr)
    if err != nil {
        panic(err)
    }

    defer r.Body.Close() //run this function whenever the current function i'm in
    //exits.

    //get teh url from the request (undecorated url)
    reqUrl := r.Request.URL

    //make an object of the URL
    baseUrl := &url.URL {
        Scheme: reqUrl.Scheme,
        Host: reqUrl.Host,
    }

    base := baseUrl.String() //make into string

    pages := hrefs(r.Body, base) //get the hrefs from the page

    for _, l := range pages { //loop through each of the links (pages)
        pages = append(pages, l)
    }

    //include closure for function
    //this filters the urls to match the base domain of the url we are searching
    //for.
    ret := filter(withPrefix(base), pages)
    return ret//return the pages

}


func hrefs(html io.Reader, base string) []string {
    //grab the links from the webpage received (parse the hrefs from the DOM)
    links, _ := crawler.Parse(html)

    //make slice of strings.
    var ret []string
    //loop through each link and cleanup the url.
    fmt.Println(base)
    for _, l := range links {
        switch {
        case strings.HasPrefix(l.Href, "/"):
            ret = append(ret, base + l.Href)
            continue
        case strings.HasPrefix(l.Href, "http"):
            ret = append(ret, l.Href)
            continue
        default:
            continue

        }
    }

    return ret
}



