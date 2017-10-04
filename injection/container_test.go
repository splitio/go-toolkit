package injection

import (
	"fmt"
	"testing"
)

// ObjectA bla bla bla
type ObjectA struct {
	*Context
	Value int
}

func TestDI(t *testing.T) {

	ctx := NewContext()
	ctx.AddDependency("numero", 12)
	ctx.AddDependency("string", "bla ble bli")

	obj := &ObjectA{Value: 15}

	ctx.Inject(obj)

	fmt.Println("asdasdas")
	fmt.Println(obj.Dependency("numero"))
	fmt.Println(obj.Dependency("string"))
	fmt.Println(obj.Value)
}
