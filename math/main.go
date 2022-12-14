package main

import (
	"fmt"
	"math-parse/lexer"
	"math-parse/parser"
)

func main() {
	lexer := lexer.Lexer{
		Formula: "1+2+3+{FOUR}",
	}
	tokens, err := lexer.Lex()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, token := range tokens {
		fmt.Printf("%+v \n", token)
	}
	p := &parser.Parser{
		Tokens: tokens,
	}
	p.SetVar("FOUR", 4)
	ast, err := p.Parse()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%.15f \n", ast.Evaluate())
	}
}
