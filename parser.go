package gojimongo

import (
	"fmt"
	"strconv"
)

const (
	PARSER_ERROR_EXPECTED_CLOSING_BRACKET = "Expected closing bracket"
	PARSER_ERROR_EXPECTED_CLOSING_PARENT = "Expected closing parenthesis"
	PARSER_ERROR_SYNTAX_SLICE = "Slice syntax not correct"
	PARSER_ERROR_UNEXEXPECTED_TOKEN = "Unexpected token error"
	PARSER_ERROR_DOT_CHILD_ERROR = "Dot child segment error"
	PARSER_ERROR_UNKNOWN = "Unknown error"
	PARSER_ERROR_MISSING_OPENING_PAREN_FUNCTION = "missing opening parenthesis in function expression"
	PARSER_ERROR_MISSING_CLOSING_PAREN_FUNCTION = "missing closing parenthesis in function expression"
	PARSER_ERROR_UNTERMINATED_ERROR = "unterminated query"
	PARSER_ERROR_NOTANINTEGER = "failed to parse integer"
	PARSER_ERROR_MISSING_CLOSING_PAREN = "missing closing parenthesis"
	PARSER_ERROR_MISSING_CLOSING_BRACKET = "missing closing bracket"
	PARSER_ERROR_MISSING_OPENING_PAREN = "missing opening parenthesis"
	PARSER_ERROR_TERMINATING_SELECTORS_WITH_COMMA = "extra comma at the end of bracketed selectors"
	PARSER_ERROR_INCORRECT_DESCENDANT_SEGMENT_SYNTAX = "incorrect descendant segment syntax"
	PARSER_ERROR_EMPTY_BRACKETED_SELECTORS = "empty bracketed selectors"
)

type Visitor interface {
	// LITERAL EXPRESSIONS
	visitStringExpr(value *StringExpr)
	visitIntExpr(value *IntExpr)
	visitTrueExpr(value *TrueExpr)
	visitFalseExpr(value *FalseExpr)
	visitNullExpr(value *NullExpr)
	visitParExpr(value *ParExpr)
	visitFnExpr(value *FnExpr)

	// UNARY EXPRESSIONS
	visitNotExpr(value *NotExpr)
	visitMinusExpr(value *MinusExpr)

	// BINARY EXPRESSIONS
	visitAndExpr(value *AndExpr)
	visitOrExpr(value *OrExpr)
	visitGtExpr(value *GtExpr)
	visitLtExpr(value *LtExpr)
	visitLteExpr(value *LteExpr)
	visitGteExpr(value *GteExpr)
	visitEqeqExpr(value *EqeqExpr)
	visitNeqExpr(value *NeqExpr)

	visitSliceSelector(value *SliceSelector)
	visitIndexSelector(value *IndexSelector)
	visitNameSelector(value *NameSelector)
	visitFilterSelector(value *FilterSelector)
	visitWildcardSelector(value *WildCardSelector)

	visitDotChildSegment(value *DotChildSegment)
	visitChildSegment(value *ChildSegment)
	visitDescendantSegment(value *DescendantSegment)

	visitAbsQuery(value *AbsQuery)
	visitRelQuery(value *RelQuery)
}

type VisitorPrinter struct {

}

func (p *Parser) error(value string) error {
	return fmt.Errorf("[gojimongo][parser]: %s", value)
} 

// QUERIES
type Query interface {
	accept(visitor Visitor)
}

type RelQuery struct { 
	segments  []Segment 
}

type AbsQuery struct { segments  []Segment }

// SEGMENTS
type Segment interface{
	accept(visitor Visitor)
}

type ChildSegment 		struct { selectors []Selector }
type DotChildSegment 	struct { selector Selector }
type DescendantSegment 	struct { selectors []Selector }

// SELECTORS
type Selector interface{ 
	accept(visitor Visitor)
}
type NameSelector struct { name string }
type WildCardSelector struct {}
type IndexSelector struct { value Expr }
type SliceSelector struct { start Expr; stop Expr; step Expr }
type FilterSelector struct { cond Expr }
type FnExpr struct { name string; params []Expr }

// EXPRESSIONS
type Expr interface{
	accept(visitor Visitor)
}

type BinaryExpr 	interface{
	accept(visitor Visitor)
}
type ComparisonExpr interface{
	accept(Visitor)
}
type RelationalExpr interface{
	accept(Visitor)
}
type UnaryExpr  	interface{
	accept(Visitor)
}
type LiteralExpr    interface{
	accept(Visitor)
}

type AbsQueryExpr interface {
	accept(Visitor)
}

type RelQueryExpr interface {
	accept(Visitor)
}

type AndExpr 	struct { lhs Expr; rhs Expr }
type OrExpr 	struct { lhs Expr; rhs Expr }
type NeqExpr 	struct { lhs Expr; rhs Expr }
type EqeqExpr 	struct { lhs Expr; rhs Expr }
type GtExpr 	struct { lhs Expr; rhs Expr }
type GteExpr 	struct { lhs Expr; rhs Expr }
type LtExpr 	struct { lhs Expr; rhs Expr }
type LteExpr 	struct { lhs Expr; rhs Expr }

type NotExpr 	struct { expr Expr }
type MinusExpr 	struct { expr Expr }

type ParExpr 	struct { value Expr }
type StringExpr struct { value string }
type IntExpr 	struct { value int }
type FalseExpr 	struct { }
type TrueExpr 	struct { }
type NullExpr	struct { }

type Parser struct {
	tokens 		[]TokenType
	lexemes 	[]string
	currtok  	int
	currlex 	int
}

func (p *Parser) Run() (Query, error) {
	return p.parse()
}

func (p *Parser) parse() (Query, error) {
	return p.query()
}

func (p *Parser) query() (Query, error) {
	if p.matchCurr(DOLLAR) {
		p.advTok()
		return p.absQuery()
	} else if p.matchCurr(AT) {
		p.advTok()
		return p.relQuery()
	}
	return nil, p.error(PARSER_ERROR_UNEXEXPECTED_TOKEN)
}

func (p *Parser) relQuery() (Query, error) {
	segments, err := p.segments()
	if err != nil {
		return nil, err
	}
	query := &RelQuery{segments: segments}
	return query, nil
}

func (p *Parser) absQuery() (Query, error) {
	segments, err := p.segments()
	if err != nil {
		return nil, err
	}
	query := &AbsQuery{segments: segments}
	return query, nil
}

func (p *Parser) segments() ([]Segment, error) {
	segments := []Segment{}
	for p.notatend() {
		segment, err := p.segment()
		if err != nil {
			return nil, err
		}
		if segment == nil {
			break
		}
		segments = append(segments, segment)
	}
	return segments, nil
}

func(p *Parser) segment() (Segment, error) {
	if p.matchCurr(LBRACK) {
		p.advTok()
		segment, err := p.bracketedChildSegment()
		if err != nil {
			return nil, err
		}
		return segment, nil
	} else if p.matchCurr(DOT) {
		p.advTok()
		return p.dotChildSegment()
	} else if p.matchCurr(RECURSIVE_OP) {
		p.advTok()
		return p.descendantSegment()
	} else {
		return nil, nil
	}
}

// here we have 5 types of selectors we can choose from
// index
// slice
// name
// wildcard
// filter

// index and slice expect an integer, both positive and negative
// name and wildcard are strings, but wildcard corresponds to the '*' string
// filter expects a QUESTION_MARK

// TODO: this is buggy: we might also have a colon
// we also need to allow this type of selectors [0, 1, "name"]
// so basically a selector can be just any literal or slice
func (p *Parser) bracketedChildSegment() (Segment, error) {
	segment := &ChildSegment{}
	selectors, err := p.bracketedSelectors()
	if err != nil {
		return nil, err
	}
	segment.selectors = selectors
	return segment, nil
}

func (p *Parser) dotChildSegment() (Segment, error) {
	segment := &DotChildSegment{}
	selector, err := p.dotSelector()
	if err != nil {
		return nil, err
	}
	segment.selector = selector
	return segment, nil
}

// dot selectors only "one" string and wildcard
func (p *Parser) dotSelector() (Selector, error) {
	if p.matchCurr(IDENTIFIER) {
		lex := p.currLex()
		p.advTok()
		p.advLex()
		return &NameSelector{name: lex}, nil
	} else if p.matchCurr(STAR){
		return &WildCardSelector{}, nil
	} else {
		return nil, p.error(PARSER_ERROR_DOT_CHILD_ERROR)
	}
}

func (p *Parser) descendantSegment() (Segment, error) {
	selectors := []Selector{}
	if p.matchCurr(IDENTIFIER) {
		s, err := p.dotSelector()
		if err != nil {
			return nil, err
		}
		selectors = append(selectors, s)
	} else if p.matchCurr(LBRACK) {
		p.advTok()
		s, err := p.bracketedSelectors()
		if err != nil {
			return nil, err
		}
		selectors = s
	} else if p.matchCurr(STAR) {
		p.advTok()
		selectors = append(selectors, &WildCardSelector{})
	} else {
		return nil, p.error(PARSER_ERROR_INCORRECT_DESCENDANT_SEGMENT_SYNTAX)
	}
	return &DescendantSegment{selectors: selectors}, nil
}

func (p *Parser) bracketedSelectors() ([]Selector, error) {
	selectors := []Selector{}
	for !p.matchCurr(RBRACK) {
		if !p.notatend() {
			return nil, p.error(PARSER_ERROR_MISSING_CLOSING_BRACKET)
		}

		if p.matchCurr(STAR) {
			p.advTok()
			selectors = append(selectors, &WildCardSelector{})
		} else if p.matchCurr(MINUS) && p.matchNext(INTEGER) && p.matchOffset(2, COLON) || 
			p.matchCurr(COLON) || 
			p.matchCurr(INTEGER) && p.matchNext(COLON) {
			slice, err := p.slice()
			if err != nil {
				return nil, err
			}
			selectors = append(selectors, slice)
		} else if p.matchCurr(MINUS) || p.matchCurr(INTEGER) {
			index, err := p.index()
			if err != nil {
				return nil, err
			}
			selectors = append(selectors, index)
		} else if p.matchCurr(STRING) {
			name, err := p.name()
			if err != nil {
				return nil, err
			}
			selectors = append(selectors, name)
		} else if p.matchCurr(COMMA) {
			p.advTok()
			// check if we are not terminating with a comma
			if p.matchCurr(RBRACK) {
				return nil, p.error(PARSER_ERROR_TERMINATING_SELECTORS_WITH_COMMA)
			}
		} else if p.matchCurr(QUESTION_MARK) {
			p.advTok()
			filter, err := p.filter()
			if err != nil {
				return nil, err
			}
			selectors = append(selectors, filter)
		}
	}
	if len(selectors) == 0 {
		return nil, p.error(PARSER_ERROR_EMPTY_BRACKETED_SELECTORS)
	}

	p.advTok()
	return selectors, nil
}

// TODO: potentially check the character that are or are not valid
func (p *Parser) name() (Selector, error) {
	lex := p.currLex()
	p.advTok()
	p.advLex()
	return &StringExpr{value: lex}, nil
}

func (p *Parser) slice() (Selector, error) {
	selector := &SliceSelector{}
	sliceParts := []Expr{nil, nil, nil} // start, stop, step
	partIndex := 0

	expectValue := func() (Expr, error) {
		if p.matchCurr(MINUS) || p.matchCurr(INTEGER) {
			return p.unary()
		}
		return nil, nil
	}

	for p.notatend() {
		if partIndex >= 3 {
			return nil, p.error(PARSER_ERROR_SYNTAX_SLICE)
		}

		// Handle empty values like `:` or `::`
		if p.matchCurr(COLON) {
			sliceParts[partIndex] = nil
		} else {
			value, err := expectValue()
			if err != nil {
				return nil, err
			}
			sliceParts[partIndex] = value
		}

		partIndex++

		if !p.matchCurr(COLON) {
			break
		}
		p.advTok()
	}

	selector.start = sliceParts[0]
	selector.stop = sliceParts[1]
	selector.step = sliceParts[2]
	return selector, nil
}

func (p *Parser) index() (Selector, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	return &IndexSelector{value: expr}, nil
}

func (p *Parser) filter() (Selector, error) {
	cond, err := p.expr()
	if err != nil {
		return nil, err
	}
	return &FilterSelector{cond: cond}, nil
}

func (p *Parser) fn() (Expr, error) {
	expr := &FnExpr{}
	if p.matchCurr(IDENTIFIER) {
		lex := p.currLex()
		expr.name = lex
		p.advTok()
		p.advLex()
	}
	if !p.matchCurr(LPAREN) {
		return nil, p.error(PARSER_ERROR_MISSING_OPENING_PAREN_FUNCTION)
	}
	p.advTok()
	params, err := p.params()
	if err != nil {
		return nil, err
	}
	expr.params = params

	if !p.matchCurr(RPAREN) {
		return nil, p.error(PARSER_ERROR_MISSING_CLOSING_PAREN_FUNCTION)
	}
	p.advTok()
	return expr, nil
}

func (p *Parser) params() ([]Expr, error) {
	params := []Expr{}
	for p.notatend() {
		expr, err := p.expr()
		if err != nil {
			return nil, err
		}
		params = append(params, expr)
		if p.matchCurr(COMMA) {
			p.advTok()
		} else {
			break
		}
	}
	return params, nil
}

func (p *Parser) expr() (Expr, error) {
	return p.and()
}

// ===================================================
// BINARY EXPRESSIONS
// and -> or -> gt|lt|lte|gte -> eqeq|neq -> unary
// ===================================================
func (p *Parser) and() (BinaryExpr, error) {
	lhs, err := p.or()
	if err != nil {
		return nil, err
	}
	for p.matchCurr(AND) {
		p.advTok()
		rhs, err := p.or()
		if err != nil {
			return nil, err
		}
		lhs = &AndExpr{lhs: lhs, rhs: rhs}
	}
	return lhs, nil
}

func (p *Parser) or() (BinaryExpr, error) {
	lhs, err := p.relation()
	if err != nil {
		return nil, err		
	}
	for p.matchCurr(OR) {
		p.advTok()
		rhs, err := p.relation()
		if err != nil {
			return nil, err
		}
		lhs = &OrExpr{lhs: lhs, rhs: rhs}
	}
	return lhs, nil
}

func (p *Parser) relation() (BinaryExpr, error) {
	lhs, err := p.comparison()
	if err != nil {
		return nil, err
	}
	if p.matchCurr(LT) {
		p.advTok()
		rhs, err := p.comparison()
		if err != nil {
			return nil, err
		}
		return &LtExpr{lhs: lhs, rhs: rhs}, nil
	} else if p.matchCurr(LTE) {
		p.advTok()
		rhs, err := p.comparison()
		if err != nil {
			return nil, err
		}
		return &LteExpr{lhs: lhs, rhs: rhs}, nil
	} else if p.matchCurr(GT) {
		p.advTok()
		rhs, err := p.comparison()
		if err != nil {
			return nil, err
		}
		return &GtExpr{lhs: lhs, rhs: rhs}, nil
	} else if p.matchCurr(GTE) {
		p.advTok()
		rhs, err := p.comparison()
		if err != nil {
			return nil, err
		}
		return &GteExpr{lhs: lhs, rhs: rhs}, nil
	} else {
	 	return lhs, nil
	}
}

func (p *Parser) comparison() (BinaryExpr, error) {
	lhs, err := p.unary()
	if err != nil {
		return nil, err
	}
	if p.matchCurr(NEQ) {
		p.advTok()
		rhs, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &NeqExpr{lhs: lhs, rhs: rhs}, nil
	} else if p.matchCurr(EQEQ) {
		p.advTok()
		rhs, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &EqeqExpr{lhs: lhs, rhs: rhs}, nil
	} else {
		return lhs, nil
	}
}

// ===========================================================
// UNARY EXPRESSIONS
// not|minus -> literal
// ===========================================================

func (p *Parser) unary() (UnaryExpr, error) {
	if p.matchCurr(MINUS) {
		p.advTok()
		lhs, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &MinusExpr{expr: lhs}, nil
	} else if p.matchCurr(NOT) {
		p.advTok()
		lhs, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &NotExpr{expr: lhs}, nil
	}
	return p.literal()
}

// ===========================================================
// LITERAL EXPRESSIONS
// string | boolean | integer | null | identifier | par expr
// ===========================================================

func (p *Parser) literal() (LiteralExpr, error) {
	if !p.notatend() {
		return nil, p.error(PARSER_ERROR_UNTERMINATED_ERROR)
	}
	tok := p.currTok()
	switch tok {
	case STRING: return p.string()
	case INTEGER: return p.int()
	case FALSE: return p.false(), nil
	case TRUE: return p.true(), nil
	case NULL: return p.null(), nil
	case LPAREN: return p.par()
	case IDENTIFIER: return p.fn()
	case DOLLAR: p.advTok(); return p.absQuery()
	case AT: p.advTok(); return p.relQuery()
	}
	return nil, p.error(PARSER_ERROR_UNEXEXPECTED_TOKEN)
}

func (p *Parser) string() (LiteralExpr, error) {
	lexeme := p.currLex()
	literal := &StringExpr{value: lexeme}
	p.advTok()
	p.advLex()
	return literal, nil
}

func (p *Parser) int() (LiteralExpr, error) {
	lexeme := p.currLex()
	value, err := strconv.Atoi(lexeme)
	if err != nil {
		return nil, p.error(PARSER_ERROR_NOTANINTEGER)
	}
	p.advTok()
	p.advLex()
	return &IntExpr{value: value}, nil
}

func (p *Parser) true() LiteralExpr {
	p.advTok()
	p.advLex()
	return &TrueExpr{}
}

func (p *Parser) false() LiteralExpr {
	p.advTok()
	p.advLex()
	return &FalseExpr{}
}

func (p *Parser) null() LiteralExpr {
	p.advTok()
	p.advLex()
	return &NullExpr{}
}

func (p *Parser) par() (LiteralExpr, error) {
	p.advTok()
	e, err := p.expr()
	if err != nil {
		return nil, err
	}
	if !p.matchCurr(RPAREN) {
		return nil, p.error(PARSER_ERROR_MISSING_CLOSING_PAREN)
	}
	p.advTok()
	return &ParExpr{value: e}, nil
}

func (p *Parser) tok(offset int) TokenType {
	return p.tokens[p.currtok + offset]
}

func (p *Parser) currTok() TokenType {
	return p.tok(0)
}

func (p *Parser) matchCurr (t TokenType) bool {
	return p.matchOffset(0, t)
}

func (p *Parser) matchNext(t TokenType) bool {
	return p.matchOffset(1, t)
}

func (p *Parser) matchOffset(offset int, t TokenType) bool {
	if p.tokAccessible(offset) {
		return p.tok(offset) == t
	}
	return false
}

func (p *Parser) tokAccessible(offset int) bool {
	return p.currtok + offset < len(p.tokens)
}

func (p *Parser) lex(offset int) string {
	return p.lexemes[p.currlex + offset]
}

func (p *Parser) currLex() string {
	return p.lex(0)
}

func (p *Parser) advTok() {
	// fmt.Printf("[%s]\n", p.currTok())
	p.currtok++
}

func (p *Parser) advLex() {
	p.currlex++
}

func (p *Parser) notatend() bool {
	return p.currtok < len(p.tokens)
}