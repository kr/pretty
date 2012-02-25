package pretty

import (
	"fmt"
	"io"
	"reflect"
)

const (
	limit = 50
)

var (
	commaLFBytes     = []byte(",\n")
	curlyBytes       = []byte("{}")
	openCurlyLFBytes = []byte("{\n")
	spacePlusLFBytes = []byte(" +\n")
)

type formatter struct {
	d int
	x interface{}

	omit bool
}

// Formatter makes a wrapper, f, that will format x as go source with line
// breaks and tabs. Object f responds to the "%v" formatting verb when both the
// "#" and " " (space) flags are set, for example:
//
//     fmt.Sprintf("%# v", Formatter(x))
//
// If one of these two flags is not set, or any other verb is used, f will
// format x according to the usual rules of package fmt.
// In particular, if x satisfies fmt.Formatter, then x.Format will be called.
func Formatter(x interface{}) (f fmt.Formatter) {
	return formatter{x: x}
}

func (fo formatter) String() string {
	return fmt.Sprint(fo.x) // unwrap it
}

func (fo formatter) passThrough(f fmt.State, c rune) {
	s := "%"
	for i := 0; i < 128; i++ {
		if f.Flag(i) {
			s += string(i)
		}
	}
	if w, ok := f.Width(); ok {
		s += fmt.Sprintf("%d", w)
	}
	if p, ok := f.Precision(); ok {
		s += fmt.Sprintf(".%d", p)
	}
	s += string(c)
	fmt.Fprintf(f, s, fo.x)
}

func (fo formatter) Format(f fmt.State, c rune) {
	if c == 'v' && f.Flag('#') && f.Flag(' ') {
		fo.format(f)
		return
	}
	fo.passThrough(f, c)
}

func (fo formatter) format(w io.Writer) {
	v := reflect.ValueOf(fo.x)
	switch v.Kind() {
	case reflect.String:
		lim := limit - 8*fo.d
		s := v.String()
		z := len(s)
		n := (z + lim - 1) / lim
		if n < 1 {
			n = 1 // empty string still produces output
		}
		sep := append(spacePlusLFBytes, make([]byte, fo.d)...)
		for j := 0; j < fo.d; j++ {
			sep[3+j] = '\t'
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				w.Write(sep)
			}
			l, h := i*lim, (i+1)*lim
			if h > z {
				h = z
			}
			fmt.Fprintf(w, "%#v", s[l:h])
		}
		return
	case reflect.Ptr:
		e := v.Elem()
		if !e.IsValid() {
			fmt.Fprintf(w, "%#v", fo.x)
		} else {
			writeByte(w, '&')
			if e.CanInterface() {
				fmt.Fprintf(w, "%# v", formatter{d: fo.d, x: e.Interface()})
			} else {
				fmt.Fprint(w, e.String())
			}
		}
	case reflect.Slice:
		s := fmt.Sprintf("%#v", fo.x)
		if len(s) < limit {
			io.WriteString(w, s)
			return
		}

		t := v.Type()

		io.WriteString(w, reflect.TypeOf(fo.x).String())
		w.Write(openCurlyLFBytes)
		for i := 0; i < v.Len(); i++ {
			for j := 0; j < fo.d+1; j++ {
				writeByte(w, '\t')
			}
			if v.Index(i).CanInterface() {
				inner := formatter{d: fo.d + 1, x: v.Index(i).Interface(), omit: t.Elem().Kind() != reflect.Interface}
				fmt.Fprintf(w, "%# v", inner)
			} else {
				fmt.Fprint(w, v.Index(i).String())
			}
			w.Write(commaLFBytes)
		}
		for j := 0; j < fo.d; j++ {
			writeByte(w, '\t')
		}
		writeByte(w, '}')
	case reflect.Struct:
		t := v.Type()
		if tryDeepEqual(reflect.Zero(t).Interface(), fo.x) {
			if !fo.omit {
				io.WriteString(w, t.String())
			}
			w.Write(curlyBytes)
			return
		}

		s := fmt.Sprintf("%#v", fo.x)
		if fo.omit {
			s = s[len(t.String()):]
		}
		if len(s) < limit {
			io.WriteString(w, s)
			return
		}

		if !fo.omit {
			io.WriteString(w, t.String())
		}
		w.Write(openCurlyLFBytes)
		var max int
		for i := 0; i < v.NumField(); i++ {
			if v := t.Field(i); v.Name != "" {
				if len(v.Name)+2 > max {
					max = len(v.Name) + 2
				}
			}
		}
		for i := 0; i < v.NumField(); i++ {
			if f := t.Field(i); f.Name != "" {
				for j := 0; j < fo.d+1; j++ {
					writeByte(w, '\t')
				}
				io.WriteString(w, f.Name)
				writeByte(w, ':')
				for j := len(f.Name) + 1; j < max; j++ {
					writeByte(w, ' ')
				}
				if v.Field(i).CanInterface() {
					inner := formatter{d: fo.d + 1, x: v.Field(i).Interface()}
					fmt.Fprintf(w, "%# v", inner)
				} else {
					io.WriteString(w, v.Field(i).String())
				}
				w.Write(commaLFBytes)
			}
		}
		for j := 0; j < fo.d; j++ {
			writeByte(w, '\t')
		}
		writeByte(w, '}')
	default:
		fmt.Fprintf(w, "%#v", fo.x)
	}
}

func tryDeepEqual(a, b interface{}) bool {
	defer func() { recover() }()
	return reflect.DeepEqual(a, b)
}

func writeByte(w io.Writer, b byte) {
	w.Write([]byte{b})
}
