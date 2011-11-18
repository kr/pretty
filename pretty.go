// Package pretty provides pretty-printing for go values. This is useful during
// debugging, to avoid wrapping long output lines in the terminal.
//
// It provides a function, Formatter, that can be used with any function that
// accepts a format string. For example,
//
//    type LongTypeName struct {
//        longFieldName, otherLongFieldName int
//    }
//    func TestFoo(t *testing.T) {
//        var x []LongTypeName{{1, 2}, {3, 4}, {5, 6}}
//        t.Errorf("%# v", Formatter(x))
//    }
//
// This package also provides a convenience wrapper for each function in
// package fmt that takes a format string.
package pretty

import (
	"fmt"
	"io"
)

// Errorf is a convenience wrapper for fmt.Errorf.
//
// Calling Errorf(f, x, y) is equivalent to
// fmt.Errorf(f, Formatter(x), Formatter(y)).
func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, wrap(a)...)
}

// Fprintf is a convenience wrapper for fmt.Fprintf.
//
// Calling Fprintf(w, f, x, y) is equivalent to
// fmt.Fprintf(w, f, Formatter(x), Formatter(y)).
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, error error) {
	return fmt.Fprintf(w, format, wrap(a)...)
}

// Printf is a convenience wrapper for fmt.Printf.
//
// Calling Printf(f, x, y) is equivalent to
// fmt.Printf(f, Formatter(x), Formatter(y)).
func Printf(format string, a ...interface{}) (n int, errno error) {
	return fmt.Printf(format, wrap(a)...)
}

// Sprintf is a convenience wrapper for fmt.Sprintf.
//
// Calling Sprintf(f, x, y) is equivalent to
// fmt.Sprintf(f, Formatter(x), Formatter(y)).
func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, wrap(a)...)
}

func wrap(a []interface{}) []interface{} {
	w := make([]interface{}, len(a))
	for i, x := range a {
		w[i] = Formatter(x)
	}
	return w
}
