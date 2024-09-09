package schema 

import (
	"k8s.io/apimachinery/pkg/runtime"
	"fmt"
	)

func TestScheme(){
	Schema := runtime.NewScheme()
	fmt.Println(Schema)
}
