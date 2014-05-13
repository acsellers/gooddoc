package ruby

import "testing"

func TestClass(t *testing.T) {
	rc := Parse("class Foo\nend")
	if rc.Name != "Foo" {
		t.Fatal("Expected 'Foo' for class name, received '%s'", rc.Name)
	}
}

func TestModuleClass(t *testing.T) {
	rc := Parse("class Foo::Bar\nend")
	if rc.Name != "Bar" {
		t.Fatal("Expected 'Foo' for class name, received '%s'", rc.Name)
	}
	if len(rc.Modules) != 1 && rc.Modules[0] != "Foo" {
		t.Fatal("Expected []string{\"Foo\"} for class name, received '%q'", rc.Modules)
	}
}
