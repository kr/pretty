package pretty

import (
	"reflect"
	"testing"
)

func Test_customDiffPrinter_Diff(t *testing.T) {
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
		wantDesc []string
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
			wantDesc: []string{"intField: 1 != 2"},
			wantOk:   false,
		},

		{
			name: "numeric comparator",
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
			wantDesc: []string{"intField: 1 != 2"},
			wantOk:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCustomDiff()
			gotDesc, gotOk := c.Diff(tt.args.a, tt.args.b)
			if !reflect.DeepEqual(gotDesc, tt.wantDesc) {
				t.Errorf("Diff() gotDesc = %v, want %v", gotDesc, tt.wantDesc)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Diff() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
