package concurrency

import (
	"fmt"
	"time"
    "math/rand"
	// "sync"
)

func TestThreeCheck() {
	fmt.Println("Test Start")

	// sig := make(chan struct{})
	// third := make(chan chan interface{})
	// oneLevel := func(receiver chan chan interface{}) chan chan interface{} {
	// 	go func() {
	// 		second := make(chan interface{})
	// 		var wg sync.WaitGroup // Add a WaitGroup
	// 		for i := 0; i < 10; i++ {
	// 			wg.Add(1) // Increment the WaitGroup counter
	// 			go func(ss chan interface{}, idx int) {
	// 				defer wg.Done() // Decrement the counter when the goroutine completes
	// 				rand.Seed(time.Now().UnixNano())
	// 				time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	// 				ss <- idx
	// 			}(second, i)
	// 		}

	// 		go func(to chan chan interface{}, from chan interface{}) {
	// 			wg.Wait() // Wait for all senders to finish
	// 			close(second) // Close the channel only after all sends are done
	// 			to <- from
	// 		}(receiver, second)
	// 	}()
	// 	return receiver
	// }

	// check := oneLevel(third)

	// go func() {
	// 	defer close(third)
	// 	select {
	// 	case <-sig:
	// 		return
	// 	default:
	// 		for ii := range check {
	// 			for jj := range ii {
	// 				fmt.Println(jj)
	// 			}
	// 		}
	// 	}
	// }()

	// time.Sleep(time.Duration(9) * time.Second)

	// sig <- struct{}{}
}



func TestThreefour() {
	fmt.Println("Test Start")

	sig := make(chan struct{})
	third :=make(chan chan interface{})
	second := make(chan interface{})
	for i := 0; i < 100; i++ {
		go func(ss chan interface{}, idx int) {
			rand.Seed(time.Now().UnixNano())
			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
			ss <- idx
		}(second, i)
	}

	go func(to chan chan interface{}, from chan interface{} ){
		to <- from
	}(third,second)

	go func(too chan chan interface{}, end chan struct{}){
		select{
			case <-end:
				close(sig)
				close(second)
				close(third)
				return 
			default:
				for jj := range too{
					for kk := range jj {
						fmt.Println(kk)
					}
				}
		}
		
	}(third, sig)


	time.Sleep(time.Duration(6) * time.Second)
	fmt.Println("end signal")


	go func(){
		time.Sleep(time.Duration(1) * time.Second)
		sig <- struct{}{}
	}()

}