package gojimongo
import (
	"testing"
	"fmt"
)

func TestParser(t *testing.T) {
	queries := map[string]bool{
		"$.hello":                     true,
		"$['hello']":                  true,
		"$['hello'][0]":               true,
		"$.store.book[0].title":       true,
		"$..author":                   true,
		"$..book[?(@.price<10)]":      true,
		"$..book[?(@.price > 10)]":    true,
		"$..book[?(@.price > 10 && @.price < 30)]": true,
		"$..book[?(@.title)]":         true,
		"$..book[0:3]":                true,
		"$..book[:3]":                 true,
		"$..book[1:]":                 true,
		"$..book[::-1]":               true,
		"$[*]":                        true,
		"$..*":                        true,
		"@.price":                     true,
		"@.books[1].title":            true,
		"$..book[?(@.price==null)]":   true,
		"$..book[?(@.price!=null)]":   true,
		"$..book[?(@.available==true)]": true,
		"$..book[?(@.available==false)]": true,
		"$..book[?(@.price>=10)]":     true,
		"$..book[?(@.price<=10)]":     true,
		"$..[0]":                      true,
		"$..book[::]":                 true, // missing slice values
		"$.store..book":               true, // double dot in middle
		// "$..book[?(@.price - 1 < 9)]": true,
		"$.length()":                  false,
		"$.store.length()":            false,

		// Invalid queries
		"$..book[?(@.price >> 10)]":   false, // invalid operator
		"$..book[?(@.price = 10)]":    false, // single =
		"$..book[?(@.price >< 10)]":   false,
		"$..book[?(@.price <)]":       false,
		"$..book[?()]":                false,
		"$['unclosed":                 false,
		"$.store.[book]":              false,
		"$..book[?(@.price &&)]":      false,
		"$..book[?(&& @.price)]":      false,
		"$[?(@.price < 10)":           false, // unclosed filter
		"$[]": 						   false,
		"":                            false, // empty
	}
	c := &Compiler{}
	for query, shouldPass := range queries {
		_, err := c.Compile(query)
		if ((err == nil) && !shouldPass) || ((err != nil) && shouldPass) {
			fmt.Printf("Tokens: %s\n", c.lexer.tokens)
			t.Errorf("Parse(%q) = %v; expected pass: %v", query, err, shouldPass)
		}
	}
}


// func TestFilter(t *testing.T) {
// 	query := "$.book[?(@.price < 20)]"
// 	c := &Compiler{}
// 	_, err := c.Compile(query)	
// 	if err != nil {
// 		t.Errorf("Parse failed %v", err)
// 	}
// }