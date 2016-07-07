package pretty

import (
	"fmt"
	"io"
	"reflect"
)

type sbuf []string

func (s *sbuf) Write(b []byte) (int, error) {
	*s = append(*s, string(b))
	return len(b), nil
}

// Diff returns a slice where each element describes
// a difference between a and b.
func Diff(a, b interface{}) (desc []string) {
	Fdiff((*sbuf)(&desc), a, b)
	return desc
}

// Fdiff writes to w a description of the differences between a and b.
func Fdiff(w io.Writer, a, b interface{}) {
	diffWriter{w: w}.diff(reflect.ValueOf(a), reflect.ValueOf(b))
}

type diffWriter struct {
	w io.Writer
	l string // label
}

func (w diffWriter) printf(f string, a ...interface{}) {
	var l string
	if w.l != "" {
		l = w.l + ": "
	}
	fmt.Fprintf(w.w, l+f, a...)
}

func (w diffWriter) diff(av, bv reflect.Value) {
	if !av.CanInterface() || !bv.CanInterface() {
		return
	}
	if !av.IsValid() && bv.IsValid() {
		w.printf("nil != %# v", formatter{v: bv, quote: true})
		return
	}
	if av.IsValid() && !bv.IsValid() {
		w.printf("%# v != nil", formatter{v: av, quote: true})
		return
	}
	if !av.IsValid() && !bv.IsValid() {
		return
	}

	at := av.Type()
	bt := bv.Type()
	if at != bt {
		w.printf("%v != %v", at, bt)
		return
	}

	switch kind := at.Kind(); kind {
	case reflect.Bool:
		if a, b := av.Bool(), bv.Bool(); a != b {
			w.printf("%v != %v", a, b)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if a, b := av.Int(), bv.Int(); a != b {
			w.printf("%d != %d", a, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if a, b := av.Uint(), bv.Uint(); a != b {
			w.printf("%d != %d", a, b)
		}
	case reflect.Float32, reflect.Float64:
		if a, b := av.Float(), bv.Float(); a != b {
			w.printf("%v != %v", a, b)
		}
	case reflect.Complex64, reflect.Complex128:
		if a, b := av.Complex(), bv.Complex(); a != b {
			w.printf("%v != %v", a, b)
		}
	case reflect.Array:
		n := av.Len()
		for i := 0; i < n; i++ {
			w.relabel(fmt.Sprintf("[%d]", i)).diff(av.Index(i), bv.Index(i))
		}
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		if a, b := av.Pointer(), bv.Pointer(); a != b {
			w.printf("%#x != %#x", a, b)
		}
	case reflect.Interface:
		w.diff(av.Elem(), bv.Elem())
	case reflect.Map:
		ak, both, bk := keyDiff(av.MapKeys(), bv.MapKeys())
		for _, k := range ak {
			w := w.relabel(fmt.Sprintf("[%#v]", k))
			w.printf("%q != (missing)", av.MapIndex(k))
		}
		for _, k := range both {
			w := w.relabel(fmt.Sprintf("[%#v]", k))
			w.diff(av.MapIndex(k), bv.MapIndex(k))
		}
		for _, k := range bk {
			w := w.relabel(fmt.Sprintf("[%#v]", k))
			w.printf("(missing) != %q", bv.MapIndex(k))
		}
	case reflect.Ptr:
		switch {
		case av.IsNil() && !bv.IsNil():
			w.printf("nil != %# v", formatter{v: bv, quote: true})
		case !av.IsNil() && bv.IsNil():
			w.printf("%# v != nil", formatter{v: av, quote: true})
		case !av.IsNil() && !bv.IsNil():
			w.diff(av.Elem(), bv.Elem())
		}
	case reflect.Slice:
		lenA := av.Len()
		lenB := bv.Len()
		if lenA != lenB {
			w.printf("%s[%d] != %s[%d]", av.Type(), lenA, bv.Type(), lenB)
			break
		}
		for i := 0; i < lenA; i++ {
			w.relabel(fmt.Sprintf("[%d]", i)).diff(av.Index(i), bv.Index(i))
		}
	case reflect.String:
		if a, b := av.String(), bv.String(); a != b {
			w.printf("%q != %q", a, b)
		}
	case reflect.Struct:
		for i := 0; i < av.NumField(); i++ {
			w.relabel(at.Field(i).Name).diff(av.Field(i), bv.Field(i))
		}
	default:
		panic("unknown reflect Kind: " + kind.String())
	}
}

func (d diffWriter) relabel(name string) (d1 diffWriter) {
	d1 = d
	if d.l != "" && name[0] != '[' {
		d1.l += "."
	}
	d1.l += name
	return d1
}

// keyEqual compares a and b for equality.
// Both a and b must be valid map keys.
func keyEqual(a, b reflect.Value) bool {
	if a.Type() != b.Type() {
		return false
	}
	switch kind := a.Kind(); kind {
	case reflect.Int:
		a, b := a.Int(), b.Int()
		return a == b
	default:
		panic("invalid map reflect Kind: " + kind.String())
	}
}

func keyDiff(a, b []reflect.Value) (ak, both, bk []reflect.Value) {
	for _, av := range a {
		inBoth := false
		for _, bv := range b {
			if keyEqual(av, bv) {
				inBoth = true
				both = append(both, av)
				break
			}
		}
		if !inBoth {
			ak = append(ak, av)
		}
	}
	for _, bv := range b {
		inBoth := false
		for _, av := range a {
			if keyEqual(av, bv) {
				inBoth = true
				break
			}
		}
		if !inBoth {
			bk = append(bk, bv)
		}
	}
	return
}
