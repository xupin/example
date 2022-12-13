package parser

import (
	"errors"
	"fmt"
	"math"
	"math-parse/enums"
	"math-parse/lexer"
	"strconv"
)

type Parser struct {
	Tokens   []*lexer.Token
	CurToken *lexer.Token
	Index    int
	Err      error
	Params   map[string]float64
}

type Node interface {
	Evaluate() float64
}

type Number struct {
	Val float64
}

type Stmt struct {
	Type  int
	Left  Node
	Right Node
}

type Func struct {
	Name string
	Args []Node
}

type Var struct {
	Key string
	Val float64
}

func init() {
	RegFunc()
}

// 解析token
func (r *Parser) Parse() (Node, error) {
	if len(r.Tokens) == 0 {
		return nil, errors.New("the token list is empty")
	}
	if r.CurToken == nil {
		r.CurToken = r.Tokens[0]
	}
	return r.Compile(), r.Err
}

// 设置变量
func (r *Parser) SetVar(key string, value float64) {
	if r.Params == nil {
		r.Params = make(map[string]float64, 0)
	}
	r.Params[key] = value
}

// 构建树
func (r *Parser) Compile() Node {
	left := r.ParseExpr()
	right := r.ParseRight(1, left)
	return right
}

// 从左开始处理
func (r *Parser) ParseExpr() Node {
	switch r.CurToken.Type {
	case enums.LPAREN:
		return r.ParseStmt()
	case enums.LBRACE:
		return r.ParseVar()
	case enums.NUMBER:
		return r.ParseNumber()
	case enums.ADD:
		return r.ParseNumber()
	case enums.SUB:
		if t := r.NextToken(); t.Type == enums.EOF {
			r.Err = errors.New("want '0-9' but get 'eof'")
			return nil
		}
		return &Stmt{
			Type:  enums.SUB,
			Left:  &Number{},
			Right: r.ParseExpr(),
		}
	case enums.MUL:
		return r.ParseNumber()
	case enums.QUO:
		return r.ParseNumber()
	case enums.REM:
		return r.ParseNumber()
	case enums.XOR:
		return r.ParseNumber()
	case enums.FUNC:
		return r.ParseFunc()
	default:
		r.Err = errors.New("want 'expr' but get 'exception'")
		return nil
	}
}

// 处理操作符右侧
func (r *Parser) ParseRight(precedence int, left Node) Node {
	for {
		curPrec := r.Precedence()
		if curPrec < precedence {
			return left
		}
		tokenType := r.CurToken.Type
		if t := r.NextToken(); t.Type == enums.EOF {
			r.Err = errors.New("want '0-9' but get 'eof'")
			return nil
		}
		right := r.ParseExpr()
		if right == nil {
			return nil
		}
		if curPrec < r.Precedence() {
			right = r.ParseRight(curPrec, right)
			if right == nil {
				return nil
			}
		}
		left = &Stmt{
			Type:  tokenType,
			Left:  left,
			Right: right,
		}
	}
}

// 变量
func (r *Parser) ParseVar() *Var {
	if t := r.NextToken(); t.Type == enums.EOF {
		r.Err = errors.New("want 'VAR' but get 'eof'")
		return nil
	}
	key := r.CurToken.Str
	if t := r.NextToken(); t.Type == enums.EOF {
		r.Err = errors.New("want '}' but get 'eof'")
		return nil
	}
	v, ok := r.Params[key]
	if !ok {
		fmt.Printf("变量[%s]未解析 \n", key)
		v = 0
	}
	node := &Var{
		Key: key,
		Val: v,
	}
	r.NextToken()
	return node
}

// 数字
func (r *Parser) ParseNumber() *Number {
	f, err := strconv.ParseFloat(r.CurToken.Str, 64)
	if err != nil {
		fmt.Printf("解析[%+v]异常 %+v \n", r.CurToken, err)
		return &Number{}
	}
	node := &Number{
		Val: f,
	}
	r.NextToken()
	return node
}

// 表达式
func (r *Parser) ParseStmt() Node {
	if t := r.NextToken(); t.Type == enums.EOF {
		r.Err = errors.New("want '0-9' but get 'eof'")
		return nil
	}
	node := r.Compile()
	if node == nil {
		r.Err = errors.New("want '0-9' but get ' '")
		return nil
	}
	if r.CurToken.Type != enums.RPAREN {
		r.Err = fmt.Errorf("want ')' but get '%s'", r.CurToken.Str)
		return nil
	}
	r.NextToken()
	return node
}

// 函数
func (r *Parser) ParseFunc() Node {
	funcName := r.CurToken.Str
	if t := r.NextToken(); t.Type == enums.LPAREN {
		nodes := make([]Node, 0)
		for r.CurToken.Type != enums.RPAREN {
			if t := r.NextToken(); t.Type == enums.EOF {
				break
			}
			if r.CurToken.Type == enums.COMMA {
				continue
			}
			nodes = append(nodes, r.Compile())
		}
		if r.CurToken.Type != enums.RPAREN {
			r.Err = fmt.Errorf("want ')' but get '%s'", r.CurToken.Str)
			return nil
		}
		if _, ok := CallFunc[funcName]; !ok {
			fmt.Printf("函数[%s]未定义 \n", funcName)
			return &Func{
				Name: DEF_FUNC,
				Args: nodes,
			}
		}
		return &Func{
			Name: funcName,
			Args: nodes,
		}
	}
	return nil
}

// 下一个字符
func (r *Parser) NextToken() *lexer.Token {
	r.Index++
	if r.Index >= len(r.Tokens) {
		r.CurToken = &lexer.Token{
			Str:  "eof",
			Type: enums.EOF,
		}
	} else {
		r.CurToken = r.Tokens[r.Index]
	}
	return r.CurToken
}

// 优先级
func (r *Parser) Precedence() int {
	switch r.CurToken.Type {
	case enums.ADD, enums.SUB:
		return 1
	case enums.MUL, enums.QUO, enums.REM:
		return 2
	case enums.XOR:
		return 3
	default:
		return 0
	}
}

func (r *Stmt) Evaluate() float64 {
	left := r.Left.Evaluate()
	right := r.Right.Evaluate()
	switch r.Type {
	case enums.ADD:
		return left + right
	case enums.SUB:
		return left - right
	case enums.MUL:
		return left * right
	case enums.QUO:
		if right == 0 {
			fmt.Printf("表达式[%g/%g]异常，除数不能为0 \n", left, right)
			return 0
		}
		return left / right
	case enums.REM:
		if right == 0 {
			fmt.Printf("表达式[%g%%%g]异常，除数不能为0 \n", left, right)
			return 0
		}
		return float64(int64(left) % int64(right))
	case enums.XOR:
		return math.Pow(left, right)
	default:
		return 0
	}
}

func (r *Number) Evaluate() float64 {
	return r.Val
}

func (r *Func) Evaluate() float64 {
	return CallFunc[r.Name](r.Args...)
}

func (r *Var) Evaluate() float64 {
	return r.Val
}

func (r *Number) String() string {
	return fmt.Sprintf("{Type: %d, Val: %f}", enums.NUMBER, r.Val)
}

func (r *Stmt) String() string {
	return fmt.Sprintf("{Type: %d, Left: %+v, Right: %+v}", r.Type, r.Left, r.Right)
}
