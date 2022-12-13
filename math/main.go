package main

import (
	"fmt"
	"math-parse/lexer"
	"math-parse/parser"
)

func main() {
	hurt := 10.0
	relive := 2.0
	level := 49.0
	qinmi := 8000.00
	lexer := lexer.Lexer{
		Formula: "{HURT} + ({HURT}*0.6*({RELIVE}*0.01+1)) + ({LEVEL}**0.1/100*{HURT}) + ({QINMI}**0.01666667*0.010248/(2+{RELIVE}*0.001))",
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
	p.SetVar("HURT", hurt)
	p.SetVar("RELIVE", relive)
	p.SetVar("LEVEL", level)
	p.SetVar("QINMI", qinmi)
	ast, _ := p.Parse()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%.15f \n", ast.Evaluate())
	}
}
