package pretty

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_customStructuredDiffPrinter_Diff(t *testing.T) {
	type fields struct {
		opts []func(*Options)
	}
	type args struct {
		a interface{}
		b interface{}
	}

	type testStruct2 struct {
		str string
	}
	type testStruct struct {
		intField   int
		floatField float64
		child      testStruct2
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantDesc []StructuredDiff
		wantOk   bool
	}{
		{
			name: "equals",
			fields: fields{
				opts: nil,
			},
			args: args{
				a: testStruct{
					intField:   1,
					floatField: 2.3,
					child: testStruct2{
						str: "strValue",
					},
				},
				b: testStruct{
					intField:   1,
					floatField: 2.3,
					child: testStruct2{
						str: "strValue",
					},
				},
			},
			wantDesc: nil,
			wantOk:   true,
		},
		{
			name: "not equals",
			fields: fields{
				opts: nil,
			},
			args: args{
				a: testStruct{
					intField:   1,
					floatField: 2.3,
					child: testStruct2{
						str: "strValue",
					},
				},
				b: testStruct{
					intField:   2,
					floatField: 2.3,
					child: testStruct2{
						str: "strValue",
					},
				},
			},
			wantDesc: []StructuredDiff{
				{
					FieldName: "intField",
					ValueA:    "1",
					ValueB:    "2",
				},
			},
			wantOk: false,
		},
		{
			name: "numeric comparator",
			fields: fields{
				opts: []func(options *Options){
					WithNumericEpsilon(0.5),
				},
			},
			args: args{
				a: testStruct{
					intField:   1,
					floatField: 2.4,
					child: testStruct2{
						str: "strValue",
					},
				},
				b: testStruct{
					intField:   1,
					floatField: 2.3,
					child: testStruct2{
						str: "strValue",
					},
				},
			},
			wantDesc: nil,
			wantOk:   true,
		},
		{
			name: "numeric comparator 2",
			fields: fields{
				opts: []func(options *Options){
					WithNumericEpsilon(0.01),
				},
			},
			args: args{
				a: testStruct{
					intField:   1,
					floatField: 53.23,
					child: testStruct2{
						str: "strValue",
					},
				},
				b: testStruct{
					intField:   1,
					floatField: 53.24,
					child: testStruct2{
						str: "strValue",
					},
				},
			},
			wantDesc: nil,
			wantOk:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCustomDiff(tt.fields.opts...)
			gotDesc, gotOk := c.StructuredDiff(tt.args.a, tt.args.b)
			if !assert.Equal(t, tt.wantDesc, gotDesc) {
				t.Errorf("Diff() gotDesc = %v, want %v", gotDesc, tt.wantDesc)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Diff() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
