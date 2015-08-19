package pretty

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

type test struct {
	v interface{}
	s string
}

type passtest struct {
	v    interface{}
	f, s string
}

type LongStructTypeName struct {
	longFieldName      interface{}
	otherLongFieldName interface{}
}

type SA struct {
	t *T
	v T
}

type T struct {
	x, y int
}

type F int

func (f F) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "F(%d)", int(f))
}

type Stringer struct{ i int }

func (s *Stringer) String() string { return "foo" }

var long = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var passthrough = []passtest{
	{1, "%d", "1"},
	{"a", "%s", "a"},
	{&Stringer{}, "%s", "foo"},
}

func TestPassthrough(t *testing.T) {
	for _, tt := range passthrough {
		s := fmt.Sprintf(tt.f, Formatter(tt.v))
		if tt.s != s {
			t.Errorf("expected %q", tt.s)
			t.Errorf("got      %q", s)
			t.Errorf("expraw\n%s", tt.s)
			t.Errorf("gotraw\n%s", s)
		}
	}
}

var gosyntax = []test{
	{nil, `nil`},
	{"", `""`},
	{"a", `"a"`},
	{1, "int(1)"},
	{1.0, "float64(1)"},
	{[]int(nil), "[]int(nil)"},
	{[0]int{}, "[0]int{}"},
	{complex(1, 0), "(1+0i)"},
	//{make(chan int), "(chan int)(0x1234)"},
	{unsafe.Pointer(uintptr(unsafe.Pointer(&long))), fmt.Sprintf("unsafe.Pointer(0x%02x)", uintptr(unsafe.Pointer(&long)))},
	{func(int) {}, "func(int) {...}"},
	{map[int]int{1: 1}, "map[int]int{1:1}"},
	{int32(1), "int32(1)"},
	{io.EOF, `&errors.errorString{s:"EOF"}`},
	{[]string{"a"}, `[]string{"a"}`},
	{
		[]string{long},
		`[]string{"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"}`,
	},
	{F(5), "pretty.F(5)"},
	{
		SA{&T{1, 2}, T{3, 4}},
		`pretty.SA{
    t:  &pretty.T{x:1, y:2},
    v:  pretty.T{x:3, y:4},
}`,
	},
	{
		map[int]string{1: "a", 2: "b"},
		`map[int]string{1:"a", 2:"b"}`,
	},
	{
		map[int]string{1: "a", 2: "b", 3: "c"},
		`map[int]string{
    1:  "a",
    2:  "b",
    3:  "c",
}`,
	},
	{
		map[int][]byte{1: {}, 2: {}},
		`map[int][]uint8{
    1:  {},
    2:  {},
}`,
	},
	{
		map[int]T{1: {}, 2: {}},
		`map[int]pretty.T{
    1:  {},
    2:  {},
}`,
	},
	{
		map[string]T{"a": {}, "b": {}},
		`map[string]pretty.T{
    "a": {},
    "b": {},
}`,
	},
	{
		long,
		`"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"`,
	},
	{
		LongStructTypeName{
			longFieldName:      LongStructTypeName{},
			otherLongFieldName: long,
		},
		`pretty.LongStructTypeName{
    longFieldName:      pretty.LongStructTypeName{},
    otherLongFieldName: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
}`,
	},
	{
		&LongStructTypeName{
			longFieldName:      &LongStructTypeName{},
			otherLongFieldName: (*LongStructTypeName)(nil),
		},
		`&pretty.LongStructTypeName{
    longFieldName:      &pretty.LongStructTypeName{},
    otherLongFieldName: (*pretty.LongStructTypeName)(nil),
}`,
	},
	{
		[]LongStructTypeName{
			{nil, nil},
			{3, 3},
			{long, nil},
		},
		`[]pretty.LongStructTypeName{
    {},
    {
        longFieldName:      int(3),
        otherLongFieldName: int(3),
    },
    {
        longFieldName:      "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
        otherLongFieldName: nil,
    },
}`,
	},
	{
		[]interface{}{
			LongStructTypeName{nil, nil},
			[]byte{1, 2, 3},
			T{3, 4},
			LongStructTypeName{long, nil},
		},
		`[]interface {}{
    pretty.LongStructTypeName{},
    []uint8{0x1, 0x2, 0x3},
    pretty.T{x:3, y:4},
    pretty.LongStructTypeName{
        longFieldName:      "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
        otherLongFieldName: nil,
    },
}`,
	},
}

func TestGoSyntax(t *testing.T) {
	for _, tt := range gosyntax {
		s := fmt.Sprintf("%# v", Formatter(tt.v))
		if tt.s != s {
			t.Errorf("expected %q", tt.s)
			t.Errorf("got      %q", s)
			t.Errorf("expraw\n%s", tt.s)
			t.Errorf("gotraw\n%s", s)
		}
	}
}

type I struct {
	i int
	R interface{}
}

func (i *I) I() *I { return i.R.(*I) }

func TestCycle(t *testing.T) {
	type A struct{ *A }
	v := &A{}
	v.A = v

	// panics from stack overflow without cycle detection
	t.Logf("Example cycle:\n%# v", Formatter(v))

	p := &A{}
	s := fmt.Sprintf("%# v", Formatter([]*A{p, p}))
	if strings.Contains(s, "CYCLIC") {
		t.Errorf("Repeated address detected as cyclic reference:\n%s", s)
	}

	type R struct {
		i int
		*R
	}
	r := &R{
		i: 1,
		R: &R{
			i: 2,
			R: &R{
				i: 3,
			},
		},
	}
	r.R.R.R = r
	t.Logf("Example longer cycle:\n%# v", Formatter(r))

	r = &R{
		i: 1,
		R: &R{
			i: 2,
			R: &R{
				i: 3,
				R: &R{
					i: 4,
					R: &R{
						i: 5,
						R: &R{
							i: 6,
							R: &R{
								i: 7,
								R: &R{
									i: 8,
									R: &R{
										i: 9,
										R: &R{
											i: 10,
											R: &R{
												i: 11,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	// here be pirates
	r.R.R.R.R.R.R.R.R.R.R.R = r
	t.Logf("Example very long cycle:\n%# v", Formatter(r))

	i := &I{
		i: 1,
		R: &I{
			i: 2,
			R: &I{
				i: 3,
				R: &I{
					i: 4,
					R: &I{
						i: 5,
						R: &I{
							i: 6,
							R: &I{
								i: 7,
								R: &I{
									i: 8,
									R: &I{
										i: 9,
										R: &I{
											i: 10,
											R: &I{
												i: 11,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	iv := i.I().I().I().I().I().I().I().I().I().I()
	*iv = *i
	t.Logf("Example long interface cycle:\n%# v", Formatter(i))
}

func TestLess(t *testing.T) {
	zero := new(int)
	one := new(int)
	*one = 1
	data := []struct {
		i, j interface{}
		less bool
	}{
		{true, true, false},
		{false, true, true},
		{uint(0), uint(0), false},
		{uint(0), uint(1), true},
		{0., 0., false},
		{0., 1., true},
		{complex(0., 0.), complex(0., 0.), false},
		{complex(0., 0.), complex(1., 0.), true},
		{complex(1., 0.), complex(0., 0.), false},
		{[]int{0, 0}, []int{0}, false},
		{[]int{0, 0}, []int{1}, true},
		{[]int{1}, []int{0}, false},
		{map[int]int{}, map[int]int{0: 0}, true},
		//{interface{}(0), interface{}(0), false},
		{TestLess, TestCycle, false},
		{TestCycle, TestLess, true},
		{zero, zero, false},
		{zero, one, true},
		{zero, (*int)(nil), false},
		{(*int)(nil), (*int)(nil), false},
		{(*int)(nil), zero, true},
	}
	for _, line := range data {
		if less(reflect.ValueOf(line.i), reflect.ValueOf(line.j)) != line.less {
			t.Errorf("less(%q, %q) != %t", line.i, line.j, line.less)
		}
	}
	// Return value is non-deterministic. Just ensure it doesn't crash.
	less(reflect.ValueOf(struct{}{}), reflect.ValueOf(struct{}{}))
}
