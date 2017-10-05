package injection

import (
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

	if obj.Dependency("numero").(int) != 12 {
		t.Error("Invalid number")
	}

	if obj.Dependency("string").(string) != "bla ble bli" {
		t.Error("Invalid string")
	}

	if obj.Value != 15 {
		t.Error("Invalid object value")
	}

}
