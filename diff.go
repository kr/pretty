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
	if !av.IsValid() && bv.IsValid() {
		w.printf("nil != %#v", bv.Interface())
		return
	}
	if av.IsValid() && !bv.IsValid() {
		w.printf("%#v != nil", av.Interface())
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

	// numeric types, including bool
	if at.Kind() < reflect.Array {
		a, b := av.Interface(), bv.Interface()
		if a != b {
			w.printf("%#v != %#v", a, b)
		}
		return
	}

	switch at.Kind() {
	case reflect.String:
		a, b := av.Interface(), bv.Interface()
		if a != b {
			w.printf("%q != %q", a, b)
		}
	case reflect.Ptr:
		switch {
		case av.IsNil() && !bv.IsNil():
			w.printf("nil != %v", bv.Interface())
		case !av.IsNil() && bv.IsNil():
			w.printf("%v != nil", av.Interface())
		case !av.IsNil() && !bv.IsNil():
			w.diff(av.Elem(), bv.Elem())
		}
	case reflect.Struct:
		for i := 0; i < av.NumField(); i++ {
			w.relabel(at.Field(i).Name).diff(av.Field(i), bv.Field(i))
		}
	case reflect.Interface:
		w.diff(reflect.ValueOf(av.Interface()), reflect.ValueOf(bv.Interface()))
	default:
		if !reflect.DeepEqual(av.Interface(), bv.Interface()) {
			w.printf("%#v != %#v", av.Interface(), bv.Interface())
		}
	}
}

func (d diffWriter) relabel(name string) (d1 diffWriter) {
	d1 = d
	if d.l != "" {
		d1.l = d.l + "." + name
	} else {
		d1.l = name
	}
	return d1
}
