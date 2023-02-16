package pretty

type StructuredDiff struct {
	FieldName string
	ValueA    string
	ValueB    string
}

type StructuredDiffer interface {
	Print(diff StructuredDiff)
	Results() []StructuredDiff
}

type differ struct {
	diffs []StructuredDiff
}

func (d *differ) Print(diff StructuredDiff) {
	d.diffs = append(d.diffs, diff)
}

func (d *differ) Results() []StructuredDiff {
	return d.diffs
}

func NewStructuredDiffer() StructuredDiffer {
	return &differ{}
}
