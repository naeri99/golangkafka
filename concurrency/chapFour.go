package concurrency


import (
    "net/http"
	"fmt"
)


type Result struct{
	Error error
	Response *http.Response
}

func TestSol(){
	
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result{
		results:=make(chan Result)
		go func(){
			defer close(results)
			for _, url := range urls{
				var result Result 
				resp, err := http.Get(url)
				result = Result{Error: err, Response: resp}
				select {
					case <-done:
						return
					case results <- result:
				}
			}
		}()
		return results
	}

	done := make(chan interface{})
	defer close(done)

	errorCount :=0 
	urls := []string{"https://www.google.com", "https://badhost", "A", "b","c", "d" }
	for result := range checkStatus(done, urls...) {
		if result.Error != nil { 
			fmt.Printf("error: %v", result.Error)
			errorCount++
			if errorCount>=3 {
				fmt.Printf("Too many error breaking\n")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}

