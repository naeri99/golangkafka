package concurrency

import (
	"fmt"
    "sync"
	"time"
)




func TestOne(){
	var data int 
	go func(){
		data++
	}()
	time.Sleep(1*time.Second)
	if data == 0 {
		fmt.Printf("the value is %v.\n",data)
	}

}


func TestTwo(){
	var memoryAccess sync.Mutex
	var value int 
	go func(){
		time.Sleep(500*time.Millisecond)
		memoryAccess.Lock()
		value++
		memoryAccess.Unlock()
	}()

	time.Sleep(400*time.Millisecond)
	memoryAccess.Lock()
	if value ==0 {
		fmt.Printf("the value is %v.\n", value)
	}else{
		fmt.Printf("the value is %v.\n", value)
	}
	memoryAccess.Unlock()

}

type value struct{
	mu sync.Mutex
	value int 
}

func TestThree(){
	var wg sync.WaitGroup 
	printSum := func(v1, v2 * value){
		defer wg.Done() 
		v1.mu.Lock()
		defer v1.mu.Unlock()
		
		time.Sleep(2*time.Second)
		v2.mu.Lock()
		defer v2.mu.Unlock()

		fmt.Printf("sum=%v\n", v1.value+v2.value)
	}

	var a,b value
	wg.Add(2)
	go printSum(&a, &b)
	go printSum(&b, &a)
	wg.Wait()

}