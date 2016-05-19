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

func (w diffWriter) diff(av, bv reflect.Value) bool {
	if !av.IsValid() && bv.IsValid() {
		w.printf("nil != %#v", bv.Interface())
		return false
	}
	if av.IsValid() && !bv.IsValid() {
		w.printf("%#v != nil", av.Interface())
		return false
	}
	if !av.IsValid() && !bv.IsValid() {
		return true
	}

	at := av.Type()
	bt := bv.Type()
	if at != bt {
		w.printf("%v != %v", at, bt)
		return false
	}

	// numeric types, including bool
	if at.Kind() < reflect.Array {
		a, b := av.Interface(), bv.Interface()
		if a != b {
			w.printf("%#v != %#v", a, b)
			return false
		}
		return true
	}

	switch at.Kind() {
	case reflect.String:
		a, b := av.Interface(), bv.Interface()
		if a != b {
			w.printf("%q != %q", a, b)
			return false
		}
	case reflect.Ptr:
		switch {
		case av.IsNil() && !bv.IsNil():
			w.printf("nil != %v", bv.Interface())
			return false
		case !av.IsNil() && bv.IsNil():
			w.printf("%v != nil", av.Interface())
			return false
		case !av.IsNil() && !bv.IsNil():
			return w.diff(av.Elem(), bv.Elem())
		}
	case reflect.Struct:
		same := true
		hasUnexported := false
		for i := 0; i < av.NumField(); i++ {
			af := av.Field(i)
			if af.CanInterface() {
				if !w.relabel(at.Field(i).Name).diff(av.Field(i), bv.Field(i)) {
					same = false
				}
			} else {
				hasUnexported = true
			}
		}
		// We can't print the value or specific field name any differing fields without resorting
		// to some unsafe hackery, so we use reflect.DeepEqual to check if unexported fields don't
		// match and at least emit a general error.
		if same && hasUnexported && !reflect.DeepEqual(av.Interface(), bv.Interface()) {
			w.printf("unexported fields don't match")
			same = false
		}
		return same
	case reflect.Slice:
		lenA := av.Len()
		lenB := bv.Len()
		if lenA != lenB {
			w.printf("%s[%d] != %s[%d]", av.Type(), lenA, bv.Type(), lenB)
			return false
		}
		same := true
		for i := 0; i < lenA; i++ {
			if !w.relabel(fmt.Sprintf("[%d]", i)).diff(av.Index(i), bv.Index(i)) {
				same = false
			}
		}
		return same
	case reflect.Map:
		same := true
		ak, both, bk := keyDiff(av.MapKeys(), bv.MapKeys())
		for _, k := range ak {
			w := w.relabel(fmt.Sprintf("[%#v]", k.Interface()))
			w.printf("%q != (missing)", av.MapIndex(k))
			same = false
		}
		for _, k := range both {
			w := w.relabel(fmt.Sprintf("[%#v]", k.Interface()))
			if !w.diff(av.MapIndex(k), bv.MapIndex(k)) {
				same = false
			}
		}
		for _, k := range bk {
			w := w.relabel(fmt.Sprintf("[%#v]", k.Interface()))
			w.printf("(missing) != %q", bv.MapIndex(k))
			same = false
		}
		return same
	case reflect.Interface:
		return w.diff(reflect.ValueOf(av.Interface()), reflect.ValueOf(bv.Interface()))
	default:
		if !reflect.DeepEqual(av.Interface(), bv.Interface()) {
			w.printf("%# v != %# v", Formatter(av.Interface()), Formatter(bv.Interface()))
			return false
		}
	}
	return true
}

func (d diffWriter) relabel(name string) (d1 diffWriter) {
	d1 = d
	if d.l != "" && name[0] != '[' {
		d1.l += "."
	}
	d1.l += name
	return d1
}

func keyDiff(a, b []reflect.Value) (ak, both, bk []reflect.Value) {
	for _, av := range a {
		inBoth := false
		for _, bv := range b {
			if reflect.DeepEqual(av.Interface(), bv.Interface()) {
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
			if reflect.DeepEqual(av.Interface(), bv.Interface()) {
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
