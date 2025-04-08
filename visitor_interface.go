package gojimongo

// LITERAL EXPRESSIONS
func (e *StringExpr) accept(visitor Visitor) {
	visitor.visitStringExpr(e)
}

func (e *IntExpr) accept(visitor Visitor) {
	visitor.visitIntExpr(e)
}

func (e *TrueExpr) accept(visitor Visitor) {
	visitor.visitTrueExpr(e)
}

func (e *FalseExpr) accept(visitor Visitor) {
	visitor.visitFalseExpr(e)
}

func (e *NullExpr) accept(visitor Visitor) {
	visitor.visitNullExpr(e)
}

func (e *ParExpr) accept(visitor Visitor) {
	visitor.visitParExpr(e)
}

func (e *FnExpr) accept(visitor Visitor) {
	visitor.visitFnExpr(e)
}

// UNARY EXPRESSIONS
func (e *NotExpr) accept(visitor Visitor) {
	visitor.visitNotExpr(e)
}

func (e *MinusExpr) accept(visitor Visitor) {
	visitor.visitMinusExpr(e)
}

// BINARY EXPRESSIONS
func (e *AndExpr) accept(visitor Visitor) {
	visitor.visitAndExpr(e)
}

func (e *OrExpr) accept(visitor Visitor) {
	visitor.visitOrExpr(e)
}

func (e *GtExpr) accept(visitor Visitor) {
	visitor.visitGtExpr(e)
}

func (e *LtExpr) accept(visitor Visitor) {
	visitor.visitLtExpr(e)
}

func (e *LteExpr) accept(visitor Visitor) {
	visitor.visitLteExpr(e)
}

func (e *GteExpr) accept(visitor Visitor) {
	visitor.visitGteExpr(e)
}

func (e *EqeqExpr) accept(visitor Visitor) {
	visitor.visitEqeqExpr(e)
}

func (e *NeqExpr) accept(visitor Visitor) {
	visitor.visitNeqExpr(e)
}

// SELECTORS
func (s *SliceSelector) accept(visitor Visitor) {
	visitor.visitSliceSelector(s)
}

func (s *IndexSelector) accept(visitor Visitor) {
	visitor.visitIndexSelector(s)
}

func (s *NameSelector) accept(visitor Visitor) {
	visitor.visitNameSelector(s)
}

func (s *FilterSelector) accept(visitor Visitor) {
	visitor.visitFilterSelector(s)
}

func (s *WildCardSelector) accept(visitor Visitor) {
	visitor.visitWildcardSelector(s)
}

// SEGMENTS
func (s *DotChildSegment) accept(visitor Visitor) {
	visitor.visitDotChildSegment(s)
}

func (s *ChildSegment) accept(visitor Visitor) {
	visitor.visitChildSegment(s)
}

func (s *DescendantSegment) accept(visitor Visitor) {
	visitor.visitDescendantSegment(s)
}

func (q *AbsQuery) accept(visitor Visitor) {
	visitor.visitAbsQuery(q)
}

func (q *RelQuery) accept(visitor Visitor) {
	visitor.visitRelQuery(q)
}