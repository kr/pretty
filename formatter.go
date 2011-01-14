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


func Formatter(x interface{}) fmt.Formatter {
	return formatter{x: x}
}


func (fo formatter) String() string {
	return fmt.Sprint(fo.x) // unwrap it
}


func (fo formatter) passThrough(f fmt.State, c int) {
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


func (fo formatter) Format(f fmt.State, c int) {
	if c == 'v' && f.Flag('#') && f.Flag(' ') {
		fo.format(f)
		return
	}
	fo.passThrough(f, c)
}


func (fo formatter) format(w io.Writer) {
	value := reflect.NewValue(fo.x)
	switch v := value.(type) {
	case *reflect.StringValue:
		lim := limit - 8*fo.d
		s := v.Get()
		z := len(s)
		n := (z + lim - 1) / lim
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
	case *reflect.SliceValue:
		s := fmt.Sprintf("%#v", fo.x)
		if len(s) < limit {
			io.WriteString(w, s)
			return
		}

		t := v.Type().(*reflect.SliceType)
		_, keep := t.Elem().(*reflect.InterfaceType)
		io.WriteString(w, reflect.Typeof(fo.x).String())
		w.Write(openCurlyLFBytes)
		for i := 0; i < v.Len(); i++ {
			for j := 0; j < fo.d+1; j++ {
				writeByte(w, '\t')
			}
			inner := formatter{d: fo.d + 1, x: v.Elem(i).Interface(), omit: !keep}
			fmt.Fprintf(w, "%# v", inner)
			w.Write(commaLFBytes)
		}
		for j := 0; j < fo.d; j++ {
			writeByte(w, '\t')
		}
		writeByte(w, '}')
	case *reflect.StructValue:
		t := v.Type().(*reflect.StructType)
		if reflect.DeepEqual(reflect.MakeZero(t).Interface(), fo.x) {
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
				if len(v.Name) + 2 > max {
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
				inner := formatter{d: fo.d + 1, x: v.Field(i).Interface()}
				io.WriteString(w, fmt.Sprintf("%# v", inner))
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

func writeByte(w io.Writer, b byte) {
	w.Write([]byte{b})
}
