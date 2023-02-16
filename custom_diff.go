package pretty

import (
	"math"
	"reflect"
)

const defaultPrecision = 0.0000001

type customDiff struct {
}

type Comparator interface {
	Diff(a, b interface{}) (desc []string, ok bool)
	StructuredDiff(a, b interface{}) (desc []StructuredDiff, ok bool)
}

type Equals func(a, b interface{}) bool
type Float64Equals func(a, b float64) bool

type Options struct {
	customComparators map[reflect.Type]Equals
	numericComparator Float64Equals
}

func WithCustomComparators(customComparators map[reflect.Type]Equals) func(*Options) {
	return func(s *Options) {
		s.customComparators = customComparators
	}
}

var precision = defaultPrecision

func newMustAbsoluteDeltaLessThan(e float64) func(a, b float64) bool {
	return func(a, b float64) bool {
		return math.Abs(a-b) <= e+precision
	}
}

// WithNumericEpsilon - sets the maximum tolerance of absolute difference of all numeric types
func WithNumericEpsilon(epsilon float64) func(*Options) {
	return func(s *Options) {
		s.numericComparator = newMustAbsoluteDeltaLessThan(epsilon)
	}
}

func NewCustomDiff(options ...func(*Options)) Comparator {
	opts := Options{
		customComparators: make(map[reflect.Type]Equals),
	}

	for _, o := range options {
		o(&opts)
	}
	return &customDiffPrinter{
		customComparators: opts.customComparators,
		numericComparator: opts.numericComparator,
	}
}

type customDiffPrinter struct {
	customComparators map[reflect.Type]Equals
	numericComparator Float64Equals
}

func (c customDiffPrinter) Diff(a, b interface{}) (desc []string, ok bool) {
	diffPrinter{
		w:                 (*sbuf)(&desc),
		customComparators: c.customComparators,
		numericComparator: c.numericComparator,
		aVisited:          make(map[visit]visit),
		bVisited:          make(map[visit]visit),
	}.diff(reflect.ValueOf(a), reflect.ValueOf(b))
	return desc, len(desc) == 0
}

func (c customDiffPrinter) StructuredDiff(a, b interface{}) (desc []StructuredDiff, ok bool) {
	descStr := make([]string, 0)
	structuredOut := NewStructuredDiffer()
	diffPrinter{
		w:                 (*sbuf)(&descStr),
		structuredOutput:  structuredOut,
		customComparators: c.customComparators,
		numericComparator: c.numericComparator,
		aVisited:          make(map[visit]visit),
		bVisited:          make(map[visit]visit),
	}.diff(reflect.ValueOf(a), reflect.ValueOf(b))
	return structuredOut.Results(), len(structuredOut.Results()) == 0
}
