package gojimongo

import (
	"fmt"
	"strings"
)

type Lexer struct {
	lexemes []string
	tokens 	[]TokenType
	curr 	int
	start 	int
}

func (l *Lexer) Run(value string) error {
	for l.curr < len(value) {
		l.start = l.curr
		char := value[l.curr]
		switch char {
		case ':': l.tokens = append(l.tokens, COLON)
		case '(': l.tokens = append(l.tokens, LPAREN)
		case ')': l.tokens = append(l.tokens, RPAREN)
		case ',': l.tokens = append(l.tokens, COMMA)
		case '{': l.tokens = append(l.tokens, LBRACE)
		case '}': l.tokens = append(l.tokens, RBRACE)
		case '[': l.tokens = append(l.tokens, LBRACK)
		case ']': l.tokens = append(l.tokens, RBRACK)
		case '\n', '\t', '\r', ' ': break
		case '*': l.tokens = append(l.tokens, STAR)
		case '"', '\'':
			err := l.handleStringLiteral(value, char)
			if err != nil {
				return err
			}
			continue
		case '$': l.tokens = append(l.tokens, DOLLAR)
		case '?': l.tokens = append(l.tokens, QUESTION_MARK)
		case '-': l.tokens = append(l.tokens, MINUS)
		case '@': l.tokens = append(l.tokens, AT)
		case '&':
			if l.curr+1 < len(value) && value[l.curr+1] == '&' {
				l.tokens = append(l.tokens, AND)
				l.curr++
			} else {
				return fmt.Errorf("unexpected character: %c at position %d", char, l.curr+1)
			}
		case '|':
			if l.curr+1 < len(value) && value[l.curr+1] == '|' {
				l.tokens = append(l.tokens, OR)
				l.curr++
			} else {
				return fmt.Errorf("unexpected character: %c at position %d", char, l.curr+1)
			}
		case '=':
			if l.curr+1 < len(value) && value[l.curr+1] == '=' {
				l.tokens = append(l.tokens, EQEQ)
				l.curr++
			} else {
				return fmt.Errorf("unexpected character: %c at position %d", char, l.curr+1)
			}
		case '!':
			if l.curr+1 < len(value) && value[l.curr+1] == '=' {
				l.tokens = append(l.tokens, NEQ)
				l.curr++
			} else {
				l.tokens = append(l.tokens, NOT)
			}
		case '<':
			if l.curr+1 < len(value) && value[l.curr+1] == '=' {
				l.tokens = append(l.tokens, LTE)
				l.curr++
			} else {
				l.tokens = append(l.tokens, LT)
			}
		case '>':
			if l.curr+1 < len(value) && value[l.curr+1] == '=' {
				l.tokens = append(l.tokens, GTE)
				l.curr++
			} else {
				l.tokens = append(l.tokens, GT)
			}
		case '.':
			if l.curr+1 < len(value) && value[l.curr+1] == '.' {
				l.tokens = append(l.tokens, RECURSIVE_OP)
				l.curr++
			} else {
				l.tokens = append(l.tokens, DOT)
			}
		default:
			if isDigit(char) {
				for l.curr < len(value) && isDigit(value[l.curr]) {
					l.curr++
				}
				l.lexemes = append(l.lexemes, value[l.start:l.curr])
				l.tokens = append(l.tokens, INTEGER)
				continue
			}
			if isIdentifier(char) {
				for l.curr < len(value) && isIdentifier(value[l.curr]) {
					l.curr++
				}
				identifier := value[l.start:l.curr]
				switch strings.ToLower(identifier) {
				case "null":
					l.lexemes = append(l.lexemes, identifier)
					l.tokens = append(l.tokens, NULL)
				case "true":
					l.lexemes = append(l.lexemes, identifier)
					l.tokens = append(l.tokens, TRUE)
				case "false":
					l.lexemes = append(l.lexemes, identifier)
					l.tokens = append(l.tokens, FALSE)
				case "match":
					l.lexemes = append(l.lexemes, identifier)
					l.tokens = append(l.tokens, MATCH)
				case "count":
					l.lexemes = append(l.lexemes, identifier)
					l.tokens = append(l.tokens, COUNT)
				case "length":
					l.lexemes = append(l.lexemes, identifier)
					l.tokens = append(l.tokens, LENGTH)
				default:
					l.lexemes = append(l.lexemes, identifier)
					l.tokens = append(l.tokens, IDENTIFIER)
				}
				continue
			}
		}
		l.curr++
	}
	return nil
}

func (l *Lexer) handleStringLiteral(value string, quote byte) error {
	// Start scanning the string literal
	l.curr++ // Move past the opening quote
	for l.curr < len(value) && value[l.curr] != quote {
		if value[l.curr] == '\\' { // Handle escaped characters
			l.curr++ // Skip the backslash
			if l.curr < len(value) {
				// Allow escaping the quote itself
				if value[l.curr] == quote || value[l.curr] == '\\' {
					l.curr++ // Move past the escaped character
				}
			}
		} else {
			l.curr++ // Move to the next character in the string
		}
	}
	// If we exit the loop, we should have reached the closing quote
	if l.curr < len(value) && value[l.curr] == quote {
		l.curr++
		lexeme := value[l.start:l.curr]
		l.lexemes = append(l.lexemes, lexeme)
		l.tokens = append(l.tokens, STRING)
	} else {
		return fmt.Errorf("unterminated string literal starting at position %d", l.start)
	}
	return nil
}

func isDigit(char byte) bool {
	return char <= '9' && char >= '0'
}

func isIdentifier(char byte) bool {
	return char <= 'z' && char >= 'a' ||
		char <= 'Z' && char >= 'A' ||
		char == '_'
}