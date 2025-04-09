package gojimongo

import (
	// "fmt"
)

type Compiler struct {
	lexer 	*Lexer
	parser 	*Parser
}

func (c *Compiler) Compile(value string) (Query, error) {
	c.lexer = &Lexer{}
	err := c.lexer.Run(value)
	if err != nil {
		return nil, err
	}
	c.parser = &Parser{lexemes: c.lexer.lexemes, tokens: c.lexer.tokens}
	q, err := c.parser.Run()
	if err != nil {
		return nil, err 
	}
	return q, nil
}

// func (c *Compiler) Represent(q Query) {
// 	v := NewVisitorJSON()
// 	q.accept(v)
// 	fmt.Println(v.Result())
// }