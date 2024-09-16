package main 

import (
    "fmt"
    "operation/login"
    "operation/storage"
    "operation/master"
)

func main() {
    // Get the Kubernetes clientset
    clientset:= login.GetClient()
   
    // Deploy storage resources
    firstprocess := storage.Deploying(clientset)
    secondprocess := master.DeployingMaster(clientset,firstprocess)
    // Wait for the process to complete
    <- secondprocess
    close(secondprocess)
    fmt.Println("hello")
}
