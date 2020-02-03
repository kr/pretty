package pretty

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"unsafe"
)

var (
	_ Logfer   = (*testing.T)(nil)
	_ Logfer   = (*testing.B)(nil)
	_ Printfer = (*log.Logger)(nil)
)


var unexpdiffs = []difftest{
	{a: nil, b: nil},
	{a: S{A: 1}, b: S{A: 1}},

	{0, "", []string{`int != string`}},
	{0, 1, []string{`0 != 1`}},
	{S{}, new(S), []string{`pretty.S != *pretty.S`}},
	{"a", "b", []string{`"a" != "b"`}},
	{S{}, S{A: 1}, []string{`A: 0 != 1`}},
	{new(S), &S{A: 1}, []string{`A: 0 != 1`}},
	{S{S: new(S)}, S{S: &S{A: 1}}, []string{`S.A: 0 != 1`}},
	{S{}, S{I: 0}, []string{`I: nil != int(0)`}},
	{S{I: 1}, S{I: "X"}, []string{`I: int != string`}},
	{S{}, S{C: []int{1}}, []string{`C: []int[0] != []int[1]`}},
	{S{C: []int{}}, S{C: []int{1}}, []string{`C: []int[0] != []int[1]`}},
	{S{C: []int{1, 2, 3}}, S{C: []int{1, 2, 4}}, []string{`C[2]: 3 != 4`}},
	{S{}, S{A: 1, S: new(S)}, []string{`A: 0 != 1`, `S: nil != &pretty.S{}`}},

	// unexported fields of every reflect.Kind (both equal and unequal)
	{struct{ X bool }{false}, struct{ X bool }{false}, nil},
	{struct{ X bool }{false}, struct{ X bool }{true}, []string{`X: false != true`}},
	{struct{ X int }{0}, struct{ X int }{0}, nil},
	{struct{ X int }{0}, struct{ X int }{1}, []string{`X: 0 != 1`}},
	{struct{ X int8 }{0}, struct{ X int8 }{0}, nil},
	{struct{ X int8 }{0}, struct{ X int8 }{1}, []string{`X: 0 != 1`}},
	{struct{ X int16 }{0}, struct{ X int16 }{0}, nil},
	{struct{ X int16 }{0}, struct{ X int16 }{1}, []string{`X: 0 != 1`}},
	{struct{ X int32 }{0}, struct{ X int32 }{0}, nil},
	{struct{ X int32 }{0}, struct{ X int32 }{1}, []string{`X: 0 != 1`}},
	{struct{ X int64 }{0}, struct{ X int64 }{0}, nil},
	{struct{ X int64 }{0}, struct{ X int64 }{1}, []string{`X: 0 != 1`}},
	{struct{ X uint }{0}, struct{ X uint }{0}, nil},
	{struct{ X uint }{0}, struct{ X uint }{1}, []string{`X: 0 != 1`}},
	{struct{ X uint8 }{0}, struct{ X uint8 }{0}, nil},
	{struct{ X uint8 }{0}, struct{ X uint8 }{1}, []string{`X: 0 != 1`}},
	{struct{ X uint16 }{0}, struct{ X uint16 }{0}, nil},
	{struct{ X uint16 }{0}, struct{ X uint16 }{1}, []string{`X: 0 != 1`}},
	{struct{ X uint32 }{0}, struct{ X uint32 }{0}, nil},
	{struct{ X uint32 }{0}, struct{ X uint32 }{1}, []string{`X: 0 != 1`}},
	{struct{ X uint64 }{0}, struct{ X uint64 }{0}, nil},
	{struct{ X uint64 }{0}, struct{ X uint64 }{1}, []string{`X: 0 != 1`}},
	{struct{ X uintptr }{0}, struct{ X uintptr }{0}, nil},
	{struct{ X uintptr }{0}, struct{ X uintptr }{1}, []string{`X: 0 != 1`}},
	{struct{ X float32 }{0}, struct{ X float32 }{0}, nil},
	{struct{ X float32 }{0}, struct{ X float32 }{1}, []string{`X: 0 != 1`}},
	{struct{ X float64 }{0}, struct{ X float64 }{0}, nil},
	{struct{ X float64 }{0}, struct{ X float64 }{1}, []string{`X: 0 != 1`}},
	{struct{ X complex64 }{0}, struct{ X complex64 }{0}, nil},
	{struct{ X complex64 }{0}, struct{ X complex64 }{1}, []string{`X: (0+0i) != (1+0i)`}},
	{struct{ X complex128 }{0}, struct{ X complex128 }{0}, nil},
	{struct{ X complex128 }{0}, struct{ X complex128 }{1}, []string{`X: (0+0i) != (1+0i)`}},
	{struct{ X [1]int }{[1]int{0}}, struct{ X [1]int }{[1]int{0}}, nil},
	{struct{ X [1]int }{[1]int{0}}, struct{ X [1]int }{[1]int{1}}, []string{`X[0]: 0 != 1`}},
	{struct{ X chan int }{c0}, struct{ X chan int }{c0}, nil},
	{struct{ X chan int }{c0}, struct{ X chan int }{c1}, []string{fmt.Sprintf("X: %p != %p", c0, c1)}},
	{struct{ X func() }{f0}, struct{ X func() }{f0}, nil},
	{struct{ X func() }{f0}, struct{ X func() }{f1}, []string{fmt.Sprintf("X: %p != %p", f0, f1)}},
	{struct{ X interface{} }{0}, struct{ X interface{} }{0}, nil},
	{struct{ X interface{} }{0}, struct{ X interface{} }{1}, []string{`X: 0 != 1`}},
	{struct{ X interface{} }{0}, struct{ X interface{} }{""}, []string{`X: int != string`}},
	{struct{ X interface{} }{0}, struct{ X interface{} }{nil}, []string{`X: int(0) != nil`}},
	{struct{ X interface{} }{nil}, struct{ X interface{} }{0}, []string{`X: nil != int(0)`}},
	{struct{ X map[int]int }{map[int]int{0: 0}}, struct{ X map[int]int }{map[int]int{0: 0}}, nil},
	{struct{ X map[int]int }{map[int]int{0: 0}}, struct{ X map[int]int }{map[int]int{0: 1}}, []string{`X[0]: 0 != 1`}},
	{struct{ X *int }{new(int)}, struct{ X *int }{new(int)}, nil},
	{struct{ X *int }{&i0}, struct{ X *int }{&i1}, []string{`X: 0 != 1`}},
	{struct{ X *int }{nil}, struct{ X *int }{&i0}, []string{`X: nil != &int(0)`}},
	{struct{ X *int }{&i0}, struct{ X *int }{nil}, []string{`X: &int(0) != nil`}},
	{struct{ X []int }{[]int{0}}, struct{ X []int }{[]int{0}}, nil},
	{struct{ X []int }{[]int{0}}, struct{ X []int }{[]int{1}}, []string{`X[0]: 0 != 1`}},
	{struct{ X string }{"a"}, struct{ X string }{"a"}, nil},
	{struct{ X string }{"a"}, struct{ X string }{"b"}, []string{`X: "a" != "b"`}},
	{struct{ X N }{N{0}}, struct{ X N }{N{0}}, nil},
	{struct{ X N }{N{0}}, struct{ X N }{N{1}}, []string{`X.N: 0 != 1`}},
	{
		struct{ X unsafe.Pointer }{unsafe.Pointer(uintptr(0))},
		struct{ X unsafe.Pointer }{unsafe.Pointer(uintptr(0))},
		nil,
	},
	{
		struct{ X unsafe.Pointer }{unsafe.Pointer(uintptr(0))},
		struct{ X unsafe.Pointer }{unsafe.Pointer(uintptr(1))},
		[]string{`X: 0x0 != 0x1`},
	},
}

func TestUnexpDiff(t *testing.T) {
	for _, tt := range unexpdiffs {
		got := UnexpDiff(tt.a, tt.b)
		eq := len(got) == len(tt.exp)
		if eq {
			for i := range got {
				eq = eq && got[i] == tt.exp[i]
			}
		}
		if !eq {
			t.Errorf("unexported field diffing % #v", tt.a)
			t.Errorf("with    % #v", tt.b)
			diffdiff(t, got, tt.exp)
			continue
		}
	}
}

func TestUnexpFdiff(t *testing.T) {
	var buf bytes.Buffer
	UnexpFdiff(&buf, 0, 1)
	want := "0 != 1\n"
	if got := buf.String(); got != want {
		t.Errorf("UnexpFdiff(0, 1) = %q want %q", got, want)
	}
}

