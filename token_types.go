package gojimongo
import (
	"fmt"
)
type TokenType int

const (
	COLON TokenType = iota + 1

	// Delimiters
	LPAREN
	RPAREN
	COMMA
	LBRACE
	RBRACE
	LBRACK
	RBRACK
	DOLLAR
	RECURSIVE_OP		

	// Logical operators
	AND
	OR
	NEQ	
	LT
	GT
	LTE
	GTE
	EQEQ

	// Misc operators
	AT
	DOT
	QUESTION_MARK
	NOT
	PLUS
	MINUS
	SLASH
	STAR

	// Functions
	LENGTH
	COUNT
	MATCH

	// Literals and identifiers
	INTEGER
	STRING
	IDENTIFIER
	TRUE
	FALSE
	NULL
)
var TokenNames = map[TokenType]string{
	COLON:         "COLON",
	LPAREN:        "LPAREN",
	RPAREN:        "RPAREN",
	COMMA:         "COMMA",
	LBRACE:        "LBRACE",
	RBRACE:        "RBRACE",
	LBRACK:        "LBRACK",
	RBRACK:        "RBRACK",
	DOLLAR:        "DOLLAR",
	RECURSIVE_OP:  "RECURSIVE_OP",
	AND:           "AND",
	OR:            "OR",
	NEQ:           "NEQ",
	LT:            "LT",
	GT:            "GT",
	LTE:           "LTE",
	GTE:           "GTE",
	EQEQ:          "EQEQ",
	AT:            "AT",
	DOT:           "DOT",
	QUESTION_MARK: "QUESTION_MARK",
	NOT:           "NOT",
	PLUS:          "PLUS",
	MINUS:         "MINUS",
	SLASH:         "SLASH",
	STAR:          "STAR",
	LENGTH:        "LENGTH",
	COUNT:         "COUNT",
	MATCH:         "MATCH",
	INTEGER:       "INTEGER",
	STRING:        "STRING",
	IDENTIFIER:    "IDENTIFIER",
	NULL: 	   	   "NULL",
	TRUE: 		   "TRUE",
	FALSE: 			"FALSE",
}

func (t TokenType) String() string {
	if name, exists := TokenNames[t]; exists {
		return name
	}
	return fmt.Sprintf("UNKNOWN(%d)", t)
}
