package pretty

import (
	"fmt"
	"testing"
)

type test struct {
	v interface{}
	s string
}

type LongStructTypeName struct {
	longFieldName      interface{}
	otherLongFieldName interface{}
}

type T struct {
	x, y int
}

type F int


func (f F) Format(s fmt.State, c int) {
	fmt.Fprintf(s, "F(%d)", int(f))
}

var long = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var gosyntax = []test{
	{"", `""`},
	{"a", `"a"`},
	{1, "1"},
	{F(5), "F(5)"},
	{long, `"` + long[:50] + "\" +\n\"" + long[50:] + `"`},
	{
		LongStructTypeName{
			longFieldName:      LongStructTypeName{},
			otherLongFieldName: long,
		},
		`pretty.LongStructTypeName{
	longFieldName:      pretty.LongStructTypeName{},
	otherLongFieldName: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOP" +
	"QRSTUVWXYZ0123456789",
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
	{longFieldName:3, otherLongFieldName:3},
	{
		longFieldName:      "abcdefghijklmnopqrstuvwxyzABCDEFGH" +
		"IJKLMNOPQRSTUVWXYZ0123456789",
		otherLongFieldName: <nil>,
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
		`[]interface { }{
	pretty.LongStructTypeName{},
	[]byte{0x1, 0x2, 0x3},
	pretty.T{x:3, y:4},
	pretty.LongStructTypeName{
		longFieldName:      "abcdefghijklmnopqrstuvwxyzABCDEFGH" +
		"IJKLMNOPQRSTUVWXYZ0123456789",
		otherLongFieldName: <nil>,
	},
}`,
	},
}


func TestGoSyntax(t *testing.T) {
	for _, tt := range gosyntax {
		s := fmt.Sprintf("%# v", Formatter(tt.v))
		if tt.s != s {
			t.Errorf("expected %q\n", tt.s)
			t.Errorf("got      %q\n", s)
			t.Errorf("expraw %s\n", tt.s)
			t.Errorf("gotraw %s\n", s)
		}
	}
}
