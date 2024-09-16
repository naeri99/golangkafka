package main 

import (
    "fmt"
    "operation/login"
    "operation/storage"
    "operation/master"
    "operation/worker"
    "operation/controller"
)

func generateCluster(){
    clientset:= login.GetClient()
    firstprocess := storage.Deploying(clientset)
    secondprocess := master.DeployingMaster(clientset,firstprocess)
    thirdprocess := worker.DeployingWorker(clientset,secondprocess )
    fourprocess := controller.DeployingController(clientset,thirdprocess )
    // Wait for the process to complete
    <- fourprocess
    close(fourprocess)
}

func deleteCluster(){
    clientset:= login.GetClient()
    firstdelete := controller.DeletingController(clientset)
    seconddelete := worker.DeletingWorker(clientset, firstdelete )
    <-seconddelete
    close(seconddelete)
}

func main() {
    
    //generateCluster()
    
    deleteCluster()

    fmt.Println("hello")
}
