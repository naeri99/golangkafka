package kafka

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
	appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)


func clientSet()  (*kubernetes.Clientset, error) {
	userHomeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Printf("Error getting user home dir: %v\n", err)
        os.Exit(1)
    }

    // Build the kubeconfig path
    kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
    fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

    // Build the Kubernetes client config from the kubeconfig file
    kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
    if err != nil {
        fmt.Printf("Error getting kubernetes config: %v\n", err)
        os.Exit(1)
    }

    // Create the Kubernetes clientset
    clientset, err := kubernetes.NewForConfig(kubeConfig)
	
	return clientset , err
}


func DeleteKafka(names string){
	clientset, err :=clientSet()
	if err != nil {
        fmt.Printf("Error getting kubernetes config: %v\n", err)
        os.Exit(1)
    }
	deletePod(names, clientset)
	deleteNamespace(names, clientset)

}

func deletePod(name string, clientset *kubernetes.Clientset){
	namespace := name // Replace with your specific namespace

	// List StatefulSets in the specified namespace
	statefulSets, err := clientset.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	
	for _, sts := range statefulSets.Items {
		fmt.Printf("Deleting StatefulSet: %s\n", sts.Name)

		// Delete the StatefulSet
		err := clientset.AppsV1().StatefulSets(namespace).Delete(context.TODO(), sts.Name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Printf("Error deleting StatefulSet %s: %v\n", sts.Name, err)
		} else {
			fmt.Printf("Successfully deleted StatefulSet: %s\n", sts.Name)
		}
	}
}

func deleteNamespace(names string , clientset *kubernetes.Clientset){

	deletePolicy := metav1.DeletePropagationForeground

	err := clientset.CoreV1().Namespaces().Delete(context.TODO(), names, metav1.DeleteOptions{
        PropagationPolicy: &deletePolicy,
    })
    if err != nil {
        fmt.Printf("Failed to delete namespace: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Namespace %s deleted successfully\n", names)
}

func Deploying( name string ) {
 
    clientset, err := clientSet()
    if err != nil {
        fmt.Printf("Error creating Kubernetes clientset: %v\n", err)
        os.Exit(1)
    }

	namespace := &corev1.Namespace{
        ObjectMeta: metav1.ObjectMeta{
            Name: name,
        },
    }

    _, err = clientset.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
    if err != nil {
        fmt.Printf("Failed to create namespace: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Namespace 'kman' created successfully.")

    // Step 2: Create the Service
    service := &corev1.Service{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "my-headless-service",
            Namespace: "kman",
        },
        Spec: corev1.ServiceSpec{
            ClusterIP: "None",
            Selector: map[string]string{
                "app": "kcluster",
            },
        },
    }

    _, err = clientset.CoreV1().Services("kman").Create(context.TODO(), service, metav1.CreateOptions{})
    if err != nil {
        fmt.Printf("Failed to create service: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("Service 'my-headless-service' created successfully.")


	err = deployStatefulSet(clientset, "my-statefulset", "kman", "1", "kafkaman")
	if err != nil {
		fmt.Printf("Failed to deploy my-statefulset: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("StatefulSet 'my-statefulset' created successfully.")

	err = deployStatefulSet(clientset, "my-statefulset-sec", "kman", "2", "kafkamantwo")
	if err != nil {
		fmt.Printf("Failed to deploy my-statefulset-sec: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("StatefulSet 'my-statefulset-sec' created successfully.")

	err = deployStatefulSet(clientset, "my-statefulset-third", "kman", "3", "kafkaman")
	if err != nil {
		fmt.Printf("Failed to deploy my-statefulset-third: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("StatefulSet 'my-statefulset-third' created successfully.")


    // List the pods in the default namespace
    pods, err := clientset.CoreV1().Pods("kman").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        fmt.Printf("Error listing pods: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Pods in the default namespace:")
    for _, pod := range pods.Items {
        fmt.Printf("- %s\n", pod.Name)
    }
}


func deployStatefulSet(clientset *kubernetes.Clientset, name, namespace, nodeID, containerName string) error {
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "kcluster",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: "my-headless-service",
			Replicas:    int32Ptr(1), // Adjust the number of replicas as needed
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "kcluster",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "kcluster",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  containerName,
							Image: "bumory1987/kafkaman:0.006",
							Ports: []corev1.ContainerPort{
								{ContainerPort: 9092},
								{ContainerPort: 9093},
							},
							Env: []corev1.EnvVar{
								{Name: "NODE_ID", Value: nodeID},
								{Name: "CONNECTOR", Value: "1@my-statefulset-0.my-headless-service.kman.svc.cluster.local:9093,2@my-statefulset-sec-0.my-headless-service.kman.svc.cluster.local:9093,3@my-statefulset-third-0.my-headless-service.kman.svc.cluster.local:9093"},
								{Name: "IP_NAME", Value: fmt.Sprintf("INTERNAL://%s-0.my-headless-service.kman.svc.cluster.local:9092", name)},
								{Name: "SECURITYP_MAP", Value: "listener.security.protocol.map=INTERNAL:PLAINTEXT,CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT"},
								{Name: "INNER_CON", Value: "inter.broker.listener.name=INTERNAL"},
								{Name: "LISTENER", Value: "INTERNAL://0.0.0.0:9092,CONTROLLER://:9093"},
							},
						},
					},
				},
			},
		},
	}

	_, err := clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), statefulSet, metav1.CreateOptions{})
	return err
}

func int32Ptr(i int32) *int32 { return &i }

