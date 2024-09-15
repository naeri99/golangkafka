package main 

import (
    "fmt"
    "operation/login"
    "operation/storage"
)

func main() {
    // Get the Kubernetes clientset
    clientset:= login.GetClient()
   
    // Deploy storage resources
    firstprocess := storage.Deploying(clientset)
    
    // Wait for the process to complete
    <-firstprocess
    close(firstprocess)
    fmt.Println("hello")
}
