package main 

import (
    "fmt"
    "operation/login"
    "operation/storage"
    "operation/master"
    "operation/worker"
    "operation/controller"
)

func main() {
    // Get the Kubernetes clientset
    clientset:= login.GetClient()
   
    // Deploy storage resources
    firstprocess := storage.Deploying(clientset)
    secondprocess := master.DeployingMaster(clientset,firstprocess)
    thirdprocess := worker.DeployingWorker(clientset,secondprocess )
    fourprocess := controller.DeployingController(clientset,thirdprocess )
    // Wait for the process to complete
    <- fourprocess
    close(fourprocess)
    fmt.Println("hello")
}
