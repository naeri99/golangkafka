package concurrency

import ( 
	"fmt"
	"time"
)

func Parent(){
	
	stringman := func(done <-chan interface{}) chan string{
		sddd :=make( chan string)
		go func(){
			var coutIdx int 
			container:= []string{"a","b", "c", "d" ,"e"}
			for {
				select{
					case <-done:
						return
					default:
						selected := container[coutIdx%5]
						sddd <- selected
						coutIdx++
				}

			}
		}()
		return sddd

	}
	
	doWork := func( done <-chan interface{} , strings <-chan string) <-chan interface{}{
		completed :=make(chan interface{})
		go func(){
			defer fmt.Println("do work done")
			defer close(completed)
			for {
				select{
				case s:= <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return completed
	}
	done := make(chan interface{})
	chnnn :=doWork(done, stringman(done))

	go func() { 
		// Cancel the operation after 1 second.
		 time.Sleep(5 * time.Second)
		 fmt.Println("Canceling doWork goroutine...")
		 close(done)
	}()
	
	for i := range chnnn{
		fmt.Println(i)
	}
	fmt.Println("done")
}


func ORChannel(){
	var or func(channels ...<-chan interface{}) <-chan interface{}
	
	or = func(channels ...<-chan interface{})<-chan interface{} {
		switch len(channels){
			case 0: 
				return nil 
			case 1:
				return channels[0]
		}
		orDone := make(chan interface{})
	
		go func() { 
			defer close(orDone)
			switch len(channels) {
				case 2: 
					select {
					case <-channels[0]:
					case <-channels[1]:
					}
				default: 
					select {
					case <-channels[0]:
					case <-channels[1]:
					case <-channels[2]:
					case <-or(append(channels[3:], orDone)...): 
					}
				}
			}()
		 return orDone
	}

	
}

func TestSimpe(){
	var or func(indexman int, channels ...chan interface{}) chan interface{}

	lastEmit := make(chan interface{})

	defer close(lastEmit)

	or = func(indexman int, channels ...chan interface{}) chan interface{} {
		if indexman == 0 {
			if len(channels) == 0 {
				return nil
			}
			indicator := make(chan interface{})
			go func() {
				// defer close(indicator)
				
				finish := make(chan interface{})
				defer close(finish)
				select {
					case <-finish:
						return
					case <-or(1, append(channels, indicator)...):
						finish <- 1
				}
				
			}()
			return indicator
		} else {

			if len(channels) == 0 {
				return nil
			} else if len(channels) == 1 {
				return channels[0]
			} else {
				lastElement := channels[len(channels)-1]

				go func() {
					defer close(lastElement)

					// Do not close `lastElement` here, since it's the responsibility of the sender
					select {
					case <-channels[0]:
						return
					case <-or(1, channels[1:]...):
						return 
					}
				}()
				return lastElement
			}
		}
	}

	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})
	ch4 := make(chan interface{})
	ch5 := make(chan interface{})
	ch6 := make(chan interface{})


	// go func() {
	// 	ch1 <- "Message from ch1"
	// 	close(ch1)
	// }()

	// go func() {
	// 	ch2 <- "Message from ch2"
	// 	close(ch2)
	// }()

	// go func() {
	// 	ch3 <- "Message from ch2"
	// 	close(ch3)
	// }()
	// go func() {
	// 	ch4 <- "Message from ch2"
	// 	close(ch2)
	// }()	
	
	// go func() {
	// 	ch5 <- "Message from ch2"
	// 	close(ch2)
	// }()	
	go func() {
		time.Sleep(4*time.Second)
		ch6 <- "Message from ch2"
		close(ch6)
		lastEmit<-1
	}()

	out := or(0, ch1, ch2, ch2,ch3,ch4,ch5, ch6)
	

	fmt.Println(<-lastEmit)
	
	value, ok:=<-out
	if !ok {
		fmt.Println(ok) // Will print the first message received
	}else{
		fmt.Println(value)
	}
	


}



// func OrDoneTwo(){
// 	var or func(channels ... chan interface{}, indexman int) chan interface{}

// 	or = func(indexman int, channels ... chan interface{}) <-chan interface{}{
// 		if indexman ==0 {
// 			if len(channels) ==0 {
// 				return nil
// 			}
// 			indicator := make(chan interface{})
// 			go func(){
// 				defer close(indicator)
// 				<-or(1, append(channels, indicator)...)
// 			}()
// 			return indicator
// 		}else{
// 			if len(channels) ==0 {
// 				return nil
// 			}else if len(channels) == 1 {
// 				return channels[0]
// 			}else{
// 				lastElement := channels[len(channels)-1]
				
// 				go func(){
// 					defer close(lastElement)
// 					switch len(channels) {
// 						case 1 :
// 							<-lastElement
// 						default:
// 							<-or(1, channels[1:])
// 					}
// 				}()

// 				return lastElement
// 			}
// 		}
// 	}

// }