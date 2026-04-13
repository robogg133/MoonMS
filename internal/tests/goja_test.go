package tests

import (
	"testing"

	"github.com/dop251/goja"
)

const CODE = `function sum(a, b) {
	return a+b;
}
	
console.log(sum(12,2))

`

func TestTest(t *testing.T) {

	program, err := goja.Compile("test", CODE, false)
	if err != nil {
		t.Fatal(err)
	}

	r := goja.New()

	_, err = r.RunProgram(program)
	if err != nil {
		t.Fatal(err)
	}

}
