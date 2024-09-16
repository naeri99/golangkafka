package main 

import (
    "fmt"
    "operation/login"
    "operation/storage"
    "operation/master"
    "operation/worker"
)

func main() {
    // Get the Kubernetes clientset
    clientset:= login.GetClient()
   
    // Deploy storage resources
    firstprocess := storage.Deploying(clientset)
    secondprocess := master.DeployingMaster(clientset,firstprocess)
    thirdprocess := worker.DeployingWorker(clientset,secondprocess )
    // Wait for the process to complete
    <- thirdprocess
    close(thirdprocess)
    fmt.Println("hello")
}
