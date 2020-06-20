package main

import (
"io/ioutil"
"net/http"
"os"
"log"
"bufio"
"bytes"
"strings"
"flag"
"time"
"sync"
"github.com/fatih/color"
)

func main () {
	
	color.Cyan(`

   ___________ ____  _______  __
  / ____/ ___// __ \/ ____/ |/ /
 / /    \__ \/ /_/ / /_   |   / 
/ /___ ___/ / _, _/ __/  /   |  
\____//____/_/ |_/_/    /_/|_|  
                                

		V1.0
`)
	color.Green("\nSearching for CSRF Vulnerabiliies\n\n")

	var concurrency int
	flag.IntVar(&concurrency, "c", 10, "Set the concurrency for speed")
	flag.Parse()

	var wg sync.WaitGroup
	for i:= 0; i < concurrency/2; i++ {
		wg.Add(1)
		go func() {
		        // Detect potential CSRF headers,paramets in the request and response
        		detectCSRF()
			wg.Done()
		}()
		wg.Wait()
		
	}

}


func detectCSRF() {
	
	file,err := os.Create("csrf.log")
	if err != nil {
		return
	}
	defer file.Close()
	content, err := ioutil.ReadFile("tokens")
        if err != nil { 
               //Do something
               return
	}
	lines := strings.Split(string(content), "\n")
	client:=http.Client{}
	time.Sleep(500 * time.Millisecond)

	scanner:=bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url:=scanner.Text()
		requestGet,err:=http.NewRequest("GET", url, nil)
		if err!=nil {
			log.Fatal(err)
		}

		var reqgetbuffer bytes.Buffer		
		for k,v := range requestGet.Header {
			header:=k+":"+v[0]
			reqgetbuffer.WriteString(header)
			
		}
                if !strings.ContainsAny(url, strings.Join(lines, ",")) {
	               color.Red("[QUERY] Potential Vulnerability")
		       color.White("%s", url)
		       file.WriteString("[QUERY] " + url + "\n")
		}
		if !strings.ContainsAny(reqgetbuffer.String(), strings.Join(lines, ",")) {
			color.Red("[REQ GET] Potential Vulnerability")
			color.Cyan("%s", url)
			file.WriteString("[REQ GET] " + url + "\n")
		}
		responseGet,err:=client.Do(requestGet)
		if err!=nil {
			log.Fatal(err)
			
		}

		var respgetbuffer bytes.Buffer
                for k,v := range responseGet.Header {	
                        header:=k+":"+v[0]
			respgetbuffer.WriteString(header)
                }
		if !strings.ContainsAny(respgetbuffer.String(), strings.Join(lines, ",")) {
                        color.Red("[RESP GET] Potential Vulnerability")
			color.Cyan("%s", url)
			file.WriteString("[RESP GET] " + url + "\n")
                }

		requestPost,err:=http.NewRequest("POST", url, nil)
                if err!=nil {
                        log.Fatal(err)
                }

		var reqpostbuffer bytes.Buffer
                for k,v := range requestPost.Header {
                        header:=k+":"+v[0]
                        reqpostbuffer.WriteString(header)

                }
                if !strings.ContainsAny(reqpostbuffer.String(), strings.Join(lines, ",")) {
		        color.Red("[REQ POST] Potential Vulnerability") 
 			color.Cyan("%s", url) 
			file.WriteString("[REQ POST] " + url + "\n")

                }
		responsePost,err:=client.Do(requestPost)
                if err!=nil {
                        log.Fatal(err)

                }

                var resppostbuffer bytes.Buffer
                for k,v := range responsePost.Header {
                        header:=k+":"+v[0]
                        resppostbuffer.WriteString(header)

                }
		if !strings.ContainsAny(resppostbuffer.String(), strings.Join(lines, ",")) {
                        color.Red("[RESP POST] Potential Vulnerability")
			color.Cyan("%s", url)       
			file.WriteString("[RESP POST] " + url + "\n")
		}
	}	
	if err:=scanner.Err(); err != nil {
		return
        }
}
