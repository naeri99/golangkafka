package main 

import (
	"context"
	"fmt"
	"log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"operation/login"
)

func main() {

	clientset := login.GetClient()

	// Use the clientset to list pods in the "default" namespace
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list pods: %v", err)
	}

	// Print the names of the pods
	fmt.Printf("There are %d pods in the default namespace:\n", len(pods.Items))
	for _, pod := range pods.Items {
		fmt.Printf("- %s\n", pod.Name)
	}
}
