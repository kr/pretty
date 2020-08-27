package pretty_test

import (
	"fmt"
	"github.com/kr/pretty"
)

func Example() {
	type myType struct {
		a, b int
	}
	var x = []myType{{1, 2}, {3, 4}, {5, 6}}
	fmt.Printf("%# v\n", pretty.Formatter(x))

	var zeroedFields = []myType{{33, 0}, {a: 0, b: 34}}
	// Note the '+' in the format
	fmt.Printf("%# +v", pretty.Formatter(zeroedFields))
	// output:
	// []pretty_test.myType{
	//     {a:1, b:2},
	//     {a:3, b:4},
	//     {a:5, b:6},
	// }
	// []pretty_test.myType{
	//     {a:33},
	//     {b:34},
	// }
}
