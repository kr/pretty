package pretty

import (
	"fmt"
	"testing"
	"time"
)

type test struct {
	v interface{}
	s string
}

type LongStructTypeName struct {
	LongFieldName      interface{}
	OtherLongFieldName interface{}
}

type T struct {
	x, y int
}

type F int

func (f F) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "F(%d)", int(f))
}

var _ fmt.Formatter = F(0)

var long = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var gosyntax = []test{
	{"", `""`},
	{"a", `"a"`},
	{1, "1"},
	{F(5), "F(5)"},
	{long, `"` + long[:50] + "\" +\n\"" + long[50:] + `"`},
	{
		LongStructTypeName{
			LongFieldName:      LongStructTypeName{},
			OtherLongFieldName: long,
		},
		`pretty.LongStructTypeName{
	LongFieldName:      pretty.LongStructTypeName{},
	OtherLongFieldName: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOP" +
	"QRSTUVWXYZ0123456789",
}`,
	},
	{
		&LongStructTypeName{
			LongFieldName:      &LongStructTypeName{},
			OtherLongFieldName: (*LongStructTypeName)(nil),
		},
		`&pretty.LongStructTypeName{
	LongFieldName:      &pretty.LongStructTypeName{},
	OtherLongFieldName: (*pretty.LongStructTypeName)(nil),
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
	{LongFieldName:3, OtherLongFieldName:3},
	{
		LongFieldName:      "abcdefghijklmnopqrstuvwxyzABCDEFGH" +
		"IJKLMNOPQRSTUVWXYZ0123456789",
		OtherLongFieldName: <nil>,
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
	[]byte{0x1, 0x2, 0x3},
	pretty.T{x:3, y:4},
	pretty.LongStructTypeName{
		LongFieldName:      "abcdefghijklmnopqrstuvwxyzABCDEFGH" +
		"IJKLMNOPQRSTUVWXYZ0123456789",
		OtherLongFieldName: <nil>,
	},
}`,
	},
	{
		time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		`time.Time{
	sec:  <internal>,
	nsec: <internal>,
	loc:  <internal>,
}`,
	},
}

func TestGoSyntax(t *testing.T) {
	for _, tt := range gosyntax {
		s := fmt.Sprintf("%# v", Formatter(tt.v))
		if tt.s != s {
			t.Fail()

			t.Logf("expected %q\n", tt.s)
			t.Logf("got      %q\n", s)
			t.Logf("expraw %s\n", tt.s)
			t.Logf("gotraw %s\n", s)
			t.Logf("----------------")
		}
	}
}
