package gojimongo

import (
	"fmt"
	"strings"
)

type VisitorJSON struct {
	builder *strings.Builder
}

func NewVisitorJSON() *VisitorJSON {
	return &VisitorJSON{
		builder: &strings.Builder{},
	}
}

func (v *VisitorJSON) Result() string {
	return v.builder.String()
}

func (v *VisitorJSON) appendLine(line string) {
	v.builder.WriteString(line)
	v.builder.WriteString("\n")
}

func (v *VisitorJSON) visitStringExpr(e *StringExpr) {
	v.appendLine(fmt.Sprintf(`{ "type": "StringExpr", "value": "%s" }`, e.value))
}

func (v *VisitorJSON) visitIntExpr(e *IntExpr) {
	v.appendLine(fmt.Sprintf(`{ "type": "IntExpr", "value": %d }`, e.value))
}

func (v *VisitorJSON) visitTrueExpr(e *TrueExpr) {
	v.appendLine(`{ "type": "TrueExpr", "value": true }`)
}

func (v *VisitorJSON) visitFalseExpr(e *FalseExpr) {
	v.appendLine(`{ "type": "FalseExpr", "value": false }`)
}

func (v *VisitorJSON) visitNullExpr(e *NullExpr) {
	v.appendLine(`{ "type": "NullExpr", "value": null }`)
}

func (v *VisitorJSON) visitNotExpr(e *NotExpr) {
	v.appendLine(`{ "type": "NotExpr", "expr":`)
	e.expr.accept(v)
	v.appendLine(`}`)
}

func (v *VisitorJSON) visitMinusExpr(e *MinusExpr) {
	v.appendLine(`{ "type": "MinusExpr", "expr":`)
	e.expr.accept(v)
	v.appendLine(`}`)
}

func (v *VisitorJSON) visitAndExpr(e *AndExpr) {
	v.appendLine(`{ "type": "AndExpr", "lhs":`)
	e.lhs.accept(v)
	v.appendLine(`, "rhs":`)
	e.rhs.accept(v)
	v.appendLine(`}`)
}

func (v *VisitorJSON) visitOrExpr(e *OrExpr) {
	v.appendLine(`{ "type": "OrExpr", "lhs":`)
	e.lhs.accept(v)
	v.appendLine(`, "rhs":`)
	e.rhs.accept(v)
	v.appendLine(`}`)
}
func (v *VisitorJSON) visitEqeqExpr(e *EqeqExpr) {
	v.appendLine(`{ "type": "EqeqExpr", "lhs":`)
	e.lhs.accept(v)
	v.appendLine(`, "rhs":`)
	e.rhs.accept(v)
	v.appendLine(`}`)
}

func (v *VisitorJSON) visitNeqExpr(e *NeqExpr) {
	v.appendLine(`{ "type": "NeqExpr", "lhs":`)
	e.lhs.accept(v)
	v.appendLine(`, "rhs":`)
	e.rhs.accept(v)
	v.appendLine(`}`)
}

func (v *VisitorJSON) visitGtExpr(e *GtExpr) {
	v.appendLine(`{ "type": "GtExpr", "lhs":`)
	e.lhs.accept(v)
	v.appendLine(`, "rhs":`)
	e.rhs.accept(v)
	v.appendLine(`}`)
}

func (v *VisitorJSON) visitGteExpr(e *GteExpr) {
	v.appendLine(`{ "type": "GteExpr", "lhs":`)
	e.lhs.accept(v)
	v.appendLine(`, "rhs":`)
	e.rhs.accept(v)
	v.appendLine(`}`)
}

func (v *VisitorJSON) visitLtExpr(e *LtExpr) {
	v.appendLine(`{ "type": "LtExpr", "lhs":`)
	e.lhs.accept(v)
	v.appendLine(`, "rhs":`)
	e.rhs.accept(v)
	v.appendLine(`}`)
}

func (v *VisitorJSON) visitLteExpr(e *LteExpr) {
	v.appendLine(`{ "type": "LteExpr", "lhs":`)
	e.lhs.accept(v)
	v.appendLine(`, "rhs":`)
	e.rhs.accept(v)
	v.appendLine(`}`)
}

func (v *VisitorJSON) visitParExpr(e *ParExpr) {
	v.appendLine(`{ "type": "ParExpr", "value":`)
	e.value.accept(v)
	v.appendLine(`}`)
}

func (v *VisitorJSON) visitFnExpr(e *FnExpr) {
	v.appendLine(fmt.Sprintf(`{ "type": "FnExpr", "name": "%s", "params": [`, e.name))
	for i, param := range e.params {
		if i > 0 {
			v.appendLine(`,`)
		}
		param.accept(v)
	}
	v.appendLine(`]}`)
}

func (v *VisitorJSON) visitSliceSelector(s *SliceSelector) {
	v.appendLine(`{ "type": "SliceSelector", "start":`)
	s.start.accept(v)
	v.appendLine(`, "stop":`)
	s.stop.accept(v)
	v.appendLine(`, "step":`)
	s.step.accept(v)
	v.appendLine(`}`)
}

func (v *VisitorJSON) visitIndexSelector(s *IndexSelector) {
	v.appendLine(`{ "type": "IndexSelector", "value":`)
	s.value.accept(v)
	v.appendLine(`}`)
}

func (v *VisitorJSON) visitNameSelector(s *NameSelector) {
	v.appendLine(fmt.Sprintf(`{ "type": "NameSelector", "name": "%s" }`, s.name))
}

func (v *VisitorJSON) visitFilterSelector(s *FilterSelector) {
	v.appendLine(`{ "type": "FilterSelector", "cond":`)
	s.cond.accept(v)
	v.appendLine(`}`)
}

func (v *VisitorJSON) visitWildcardSelector(s *WildCardSelector) {
	v.appendLine(`{ "type": "WildCardSelector" }`)
}

func (v *VisitorJSON) visitDotChildSegment(s *DotChildSegment) {
	v.appendLine(`{ "type": "DotChildSegment", "selector":`)
	s.selector.accept(v)
	v.appendLine(`}`)
}

func (v *VisitorJSON) visitChildSegment(s *ChildSegment) {
	v.appendLine(`{ "type": "ChildSegment", "selectors": [`)
	for i, sel := range s.selectors {
		if i > 0 {
			v.appendLine(`,`)
		}
		sel.accept(v)
	}
	v.appendLine(`]}`)
}

func (v *VisitorJSON) visitDescendantSegment(s *DescendantSegment) {
	v.appendLine(`{ "type": "DescendantSegment" }`)
}

func (v *VisitorJSON) visitAbsQuery(q *AbsQuery) {
	v.appendLine(`{ "type": "AbsQuery", "segments": [`)
	for i, seg := range q.segments {
		if i > 0 {
			v.appendLine(`,`)
		}
		seg.accept(v)
	}
	v.appendLine(`]}`)
}

func (v *VisitorJSON) visitRelQuery(q *RelQuery) {
	v.appendLine(`{ "type": "RelQuery", "segments": [`)
	for i, seg := range q.segments {
		if i > 0 {
			v.appendLine(`,`)
		}
		seg.accept(v)
	}
	v.appendLine(`]}`)
}
