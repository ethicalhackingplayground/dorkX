package main

// Import the packages
import (
"time"
"net/http"
"strings"
"sync"
"flag"
"bufio"
"os"
"log"
"github.com/fatih/color"
)

// Main
func main () {
	
	// Print the banner
	color.Blue(`
	
_________                     ____  ___
\_   ___ \  ___________  _____\   \/  /
/    \  \/ /  _ \_  __ \/  ___/\     / 
\     \___(  <_> )  | \/\___ \ /     \ 
 \______  /\____/|__|  /____  >___/\  \
        \/                  \/      \_/

		V1.0
	`)

	// Arguments
	var concurrency int	
	flag.IntVar(&concurrency,"c", 10, "Set the concurrency for faster results")
	var help string
	flag.StringVar(&help, "h", "", "Display The help")
	flag.Parse()
	
	if help != "" {
		// Print the arguments
		flag.PrintDefaults()
	}else{

		color.Green("Checking for CORS...")

		var wg sync.WaitGroup
		for i:= 0; i < concurrency/2; i++ {		

			wg.Add(1)
			// Run the scanner
			go func() {
				runScanner()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

// The Cors main scanner
func runScanner() {

	// Create the results file	
	file,err := os.Create("cors.log")
	if err != nil {
		return
	}
	defer file.Close()

	// Setup the HTTP client
	timeout:= time.Duration(2 * time.Second)
	client:=http.Client {
		Timeout: timeout,
	}

	// Get the standard input
	scanner:=bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
	
		url:=scanner.Text()
		time.Sleep(time.Millisecond * 500)
		if strings.HasPrefix(url, "http") || strings.HasPrefix(url, "https") {

			color.White("[-] Checking: " + url + "\n")

			// Make a new GET Request
			reqGet,err:=http.NewRequest("GET", url, nil)
			// Check for errors
			if err !=nil { 
				return
			}

			// Make a null origin request
			reqGet.Header.Set("Origin", "null")

			// Make a new GET Request
                        reqPost,err:=http.NewRequest("POST", url, nil)
                        // Check for errors
                        if err !=nil {
                                return
                        }

                        // Make a null origin request
                        reqPost.Header.Set("Origin", "null")

			// Make the http request
			respGet,errGet:=client.Do(reqGet)
			if errGet !=nil {
				return
			}
			                

			// Make the http request
                        respPost,errPost:=client.Do(reqPost)
                        if errPost !=nil {
                                return
                        }

			// Check the response 'Access-Control-Allow-Origin' to see weather null get's reflected,
               		// Therefore if it is, it's vulnerable to CORS (Cross Origin Resource Sharing)
			for k,v := range respGet.Header {
				header:=string(k)
				value:=string(v[0])
				if header == "Access-Control-Allow-Origin" {
					 if value == "null" {
                               			// Vulnerabilitiy exists
                               			color.Red("VULNERABLE: " + url + "\n")
						color.Blue(header+":"+value)
						file.WriteString(url + "\n")	
						break
					 }else{
						color.Cyan("Not Vulnerable: " + url + "\n")
						break
					}
				}
			}


			
                        // Check the response 'Access-Control-Allow-Origin' to see weather null get's reflected,
                        // Therefore if it is, it's vulnerable to CORS (Cross Origin Resource Sharing)
                        for k,v := range respPost.Header {
                                header:=string(k)
                                value:=string(v[0])
                                if header == "Access-Control-Allow-Origin" {
                                         if value == "null" {
                                                // Vulnerabilitiy exists
                                                color.Red("VULNERABLE: " + url + "\n")
						color.Blue(header+":"+value)
                                                file.WriteString(url + "\n")
                                         }else{
                                                color.Cyan("Not Vulnerable: " + url + "\n")
                                        }
                                }
                        }
		} 
	}
	if err:=scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
