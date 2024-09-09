package concurrency

import (
	"fmt"
	"sync"


)

func TwoTestOne(){
	var wg sync.WaitGroup
	wg.Add(1)
	go sayHello(wg)
	fmt.Println("hello three")
	
}


func sayHello(wg sync.WaitGroup){
	defer wg.Done()
	fmt.Println("hello two")
}
