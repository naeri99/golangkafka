package concurrency

import (
	"fmt"
	"runtime"
	"math/rand"
	// "time"
	// "sync"
)

func PipeMan(){
	done := make(chan interface{})

	insert := func(test ...int ) chan int{
		basicInfo := make(chan int)
		go func(integers ...int){
			defer close(basicInfo)
			for _, i := range integers{
				basicInfo <-i
			}
		}(test...)
		return basicInfo
	}

	integerman:=make([]int,100)
	for i :=0; i<len(integerman); i++{
		integerman[i] = i 
	}

	add := func(original <-chan int, addPortion int) chan int{
		addInfo :=make(chan int)
		go func(origin <-chan int ){
			defer close(addInfo)
			for k := range origin{
				addInfo <- k+addPortion
			}
		}(original)
		return addInfo
	}


	mulit := func(original <-chan int, mPortion int) chan int{
		addInfo :=make(chan int)
		go func(origin <-chan int ){
			defer close(addInfo)
			for k := range origin{
				addInfo <- k*mPortion
			}
		}(original)
		return addInfo
	}


	sourceman:=insert(integerman...)

	addition:=mulit(add(sourceman,100) ,2)
	for k := range addition{
		fmt.Println(k)
	}


	close(done)
	// defer close(sourceman)


}


// func PipeMan() {
// 	done := make(chan interface{})

// 	insert := func(test ...int) chan int {
// 		basicInfo := make(chan int)
// 		go func(integers ...int) {
// 			defer close(basicInfo)
// 			for _, i := range integers {
// 				basicInfo <- i
// 			}
// 		}(test...)
// 		return basicInfo
// 	}

// 	integerman := make([]int, 100)
// 	for i := 0; i < len(integerman); i++ {
// 		integerman[i] = i
// 	}

// 	add := func(original <-chan int, addPortion int) chan int {
// 		addInfo := make(chan int)
// 		go func(origin <-chan int) {
// 			defer close(addInfo) // Close the channel only after processing all data
// 			for k := range origin {
// 				addInfo <- k + addPortion
// 			}
// 		}(original)
// 		return addInfo
// 	}

// 	sourceman := insert(integerman...)

// 	addition := add(sourceman, 100)
// 	for k := range addition {
// 		fmt.Println(k)
// 	}

// 	close(done)
// }


// func Inception(){
// 	higher := func() <-chan <-chan interface{}{
// 		middle :=make( chan <-chan interface{})
// 		go func (m  chan <-chan interface{}){
// 			defer close(middle)
// 			for i :=0 ; i <10 ; i++{
// 				lower :=make(chan interface{})
// 				go func(i int, x chan interface{}){
// 					defer close(lower)
// 					x <- i
// 				}(i, lower)
// 				m<-lower
// 			}
// 		}(middle)
// 		return middle

// 	}

// 	sample:=higher()

// 	for ii := range sample{
// 		for val := range ii { // Receive the value from the lower channel
//             fmt.Println(val)
//         }
// 	}
// }


func RepeatTest(){
	// repeatFn := func(
	// 	done <-chan interface{},
	// 	fn func() interface{},
	// 	) <-chan interface{} {
	// 		valueStream := make(chan interface{})
	// 		go func() {
	// 			defer close(valueStream)
	// 			for {
	// 				select {
	// 					case <-done:
	// 						return
	// 					case valueStream <- fn():
	// 					}
	// 				}	
	// 		}()	
	// 	return valueStream
	// }
	numFinders := runtime.NumCPU()
	fmt.Println(numFinders)


	randss := func() interface{} { return rand.Intn(50) }
	count:=0
	for {
		if count == 100{
			break
		}
			index:=randss() 
			fmt.Println(index)
			count++
	}
	
}

// func Inception(){
// 	genVals := func() <-chan <-chan interface{} {
// 		chanStream := make(chan (<-chan interface{}))
// 		go func() {
// 			defer close(chanStream)
// 				for i := 0; i < 10; i++ {
// 					stream := make(chan interface{})
// 					go func( m chan interface{}){
// 					defer close(stream)
// 						m <- i
// 					}(stream)
					
// 					chanStream <- stream
// 				}
// 			}()
// 		return chanStream
// 	}

// 	testst :=genVals()
// 	for i := range testst{
// 		for j := range i{
// 			fmt.Println(j)
// 		}
// 	}
// }


func Inception() {
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{})
				go func(m chan interface{}, val int) {
					defer close(stream)
					m <- val
				}(stream, i) // Pass `i` as an argument
				chanStream <- stream
			}
		}()
		return chanStream
	}

	testst := genVals()
	for i := range testst {
		for j := range i {
			fmt.Println(j)
		}
	}
}
