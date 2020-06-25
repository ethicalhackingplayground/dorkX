package main

import (
"github.com/PuerkitoBio/goquery"
"net/http"
"strings"
"fmt"
"sync"
"time"
"flag"
"bufio"
"os"
"log"
)


func main () {
        var concurrency=flag.Int("concurrency", 10, "Set the concurrency of the tool")
        var dork=flag.String("dork", "", "The dork to search for")
        var dorks=flag.String("dorks", "", "The list of dorks to search for")
	flag.Parse()
        var wg sync.WaitGroup
	if *dork == "" && *dorks == "" {
		flag.PrintDefaults()
		return
	}
        for i:=0; i<*concurrency/2; i++ {
                wg.Add(1)
                go func() {
			if *dork != "" && *dorks == "" {
                        	dorkScanner(*dork)
                        	wg.Done()	
			}else {
				file, err := os.Open(*dorks)
    				if err != nil {
        				log.Fatal(err)
    				}
    				defer file.Close()

				scanner:=bufio.NewScanner(file)
				for scanner.Scan() {
					dorkScanner(scanner.Text())
				}
				if err:=scanner.Err(); err!= nil {
					return
				}
				wg.Done()
			}
                }()
                wg.Wait()
        }

}

func dorkScanner(dork string) {
        time.Sleep(time.Millisecond * 50)

        res,_:=GoogleScrape(dork, "uk", "com")
	for _,item:=range res{
		url:=item.ResultURL
		fmt.Println(url)
        }
	time.Sleep(time.Second * 30)
}

type GoogleResult struct{
	ResultURL string
}




var googleDomains = map[string]string{
	"com": "https://www.google.com/search?q=",
	"uk": "https://www.google.co.uk/search?q=",
	"ru": "https://www.google.ru/search?q=",
	"fr": "https://www.google.fr/search?q=",
}

func buildGoogleUrl(searchTerm string, countryCode string, languageCode string) string {
        searchTerm = strings.Trim(searchTerm, " ")
        searchTerm = strings.Replace(searchTerm, " ", "+", -1)
        if googleBase, found := googleDomains[countryCode]; found {
                return fmt.Sprintf("%s%s&num=100&hl=%s", googleBase, searchTerm, languageCode)
        } else {
                return fmt.Sprintf("%s%s&num=100&hl=%s", googleDomains["com"], searchTerm, languageCode)
        }
}

func googleRequest(searchURL string) (*http.Response, error) {

	baseClient := &http.Client{}

	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")

	res, err := baseClient.Do(req)

	if err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func googleResultParser(response *http.Response) ([]GoogleResult, error){
        doc, err := goquery.NewDocumentFromResponse(response)
        if err != nil {
                return nil, err
        }
        results := []GoogleResult{}
        sel := doc.Find("div.g")
        rank := 1
        for i := range sel.Nodes {
                item := sel.Eq(i)
                linkTag := item.Find("a")
                link, _ := linkTag.Attr("href")
                link = strings.Trim(link, " ")
                if link != "" && link != "#"{
                        result := GoogleResult{
                                link,
                        }
                        results = append(results, result)
                        rank += 1
                }
        }
        return results, err
}

func GoogleScrape(searchTerm string, countryCode string, languageCode string) ([]GoogleResult, error) {
        googleUrl := buildGoogleUrl(searchTerm, countryCode, languageCode)
        res, err := googleRequest(googleUrl)
        if err != nil {
                return nil, err
        }
        scrapes, err := googleResultParser(res)
        if err != nil {
                return nil, err
        } else {
                return scrapes, nil
        }
}
