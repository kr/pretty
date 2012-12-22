# pretty

Package pretty provides pretty-printing for Go values. This is useful during
debugging, to avoid wrapping long output lines in the terminal.

It provides a function, Formatter, that can be used with any function that
accepts a format string. For example,

    type LongTypeName struct {
        longFieldName, otherLongFieldName int
    }
    func TestFoo(t *testing.T) {
        var x []LongTypeName{{1, 2}, {3, 4}, {5, 6}}
        t.Errorf("%# v", Formatter(x))
    }

This package also provides a convenience wrapper for each function in
package fmt that takes a format string.


## Documentation

See [GoDoc](http://godoc.org/github.com/kr/pretty) for automatic documentation.


## Installation

    $ go get github.com/kr/pretty

then

    import "github.com/kr/pretty"
