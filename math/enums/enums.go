package enums

const (
	ILLEGAL = iota
	// symbols
	EOF
	WS     // whitespace
	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }
	COMMA  // ,
	PERIOD // .

	// literals
	NUMBER // TODO 可以细分归类

	// operators
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %
	XOR // ^

	FUNC // 函数
	VAR  // 变量
)
