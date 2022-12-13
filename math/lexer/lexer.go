package lexer

import (
	"errors"
	"fmt"
	"math-parse/enums"
	"unicode"
)

type Lexer struct {
	Formula string
	Char    byte
	Pos     int
}

type Token struct {
	Str   string
	Type  int
	Start int
	End   int
}

func (r *Lexer) Lex() ([]*Token, error) {
	tokens := make([]*Token, 0)
	if len(r.Formula) == 0 {
		return tokens, errors.New("the token list is empty")
	}
	if r.IsChinese() {
		return tokens, errors.New("the token list contains Chinese characters")
	}
	r.Char = r.Formula[0]
	for r.Pos < len(r.Formula) {
		token := r.Scan()
		if token.Type == enums.ILLEGAL {
			return []*Token{}, fmt.Errorf("'%s' is not supported", token.Str)
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

// 是否包含中文
func (r *Lexer) IsChinese() bool {
	for _, c := range r.Formula {
		if (c >= 65281 && c <= 65374) || c == 12288 {
			return true
		}
	}
	return false
}

// 获取下一个token
func (r *Lexer) Scan() *Token {
	var token *Token
	pos := r.Pos
	switch r.Char {
	case '(':
		token = &Token{
			Str:   string(r.Char),
			Type:  enums.LPAREN,
			Start: pos,
			End:   pos,
		}
		r.NextChar()
	case ')':
		token = &Token{
			Str:   string(r.Char),
			Type:  enums.RPAREN,
			Start: pos,
			End:   pos,
		}
		r.NextChar()
	case '{':
		token = &Token{
			Str:   string(r.Char),
			Type:  enums.LBRACE,
			Start: pos,
			End:   r.Pos,
		}
		r.NextChar()
	case '}':
		token = &Token{
			Str:   string(r.Char),
			Type:  enums.RBRACE,
			Start: pos,
			End:   r.Pos,
		}
		r.NextChar()
	case ',':
		token = &Token{
			Str:   string(r.Char),
			Type:  enums.COMMA,
			Start: pos,
			End:   pos,
		}
		r.NextChar()
	case
		'0',
		'1',
		'2',
		'3',
		'4',
		'5',
		'6',
		'7',
		'8',
		'9':
		for r.IsDigit() {
			if !r.NextChar() {
				break
			}
		}
		token = &Token{
			Str:   string(r.Formula[pos:r.Pos]),
			Type:  enums.NUMBER,
			Start: pos,
			End:   r.Pos,
		}
	case '+':
		token = &Token{
			Str:   string(r.Char),
			Type:  enums.ADD,
			Start: pos,
			End:   pos,
		}
		r.NextChar()
	case '-':
		token = &Token{
			Str:   string(r.Char),
			Type:  enums.SUB,
			Start: pos,
			End:   pos,
		}
		r.NextChar()
	case '*':
		token = &Token{
			Str:   string(r.Char),
			Type:  enums.MUL,
			Start: pos,
			End:   pos,
		}
		r.NextChar()
		if r.Char == '*' {
			token = &Token{
				Str:   "**",
				Type:  enums.XOR,
				Start: pos,
				End:   pos,
			}
			r.NextChar()
		}
	case '/':
		token = &Token{
			Str:   string(r.Char),
			Type:  enums.QUO,
			Start: pos,
			End:   pos,
		}
		r.NextChar()
	case '%':
		token = &Token{
			Str:   string(r.Char),
			Type:  enums.REM,
			Start: pos,
			End:   pos,
		}
		r.NextChar()
	case '^':
		token = &Token{
			Str:   string(r.Char),
			Type:  enums.XOR,
			Start: pos,
			End:   pos,
		}
		r.NextChar()
	default:
		if r.IsWs() { // 跳过空白字符
			for r.NextChar() {
				if !r.IsWs() {
					break
				}
			}
			return r.Scan()
		} else if r.IsLetter() { // 判断是不是字母（函数）
			for r.IsLetter() {
				if !r.NextChar() {
					break
				}
			}
			token = &Token{
				Str:   string(r.Formula[pos:r.Pos]),
				Type:  enums.FUNC,
				Start: pos,
				End:   r.Pos,
			}
		} else if r.IsVar() { // 判断是不是变量
			for r.IsVar() {
				if !r.NextChar() {
					break
				}
			}
			token = &Token{
				Str:   string(r.Formula[pos:r.Pos]),
				Type:  enums.VAR,
				Start: pos,
				End:   r.Pos,
			}
		} else {
			token = &Token{
				Str:   string(r.Char),
				Type:  enums.ILLEGAL,
				Start: pos,
				End:   r.Pos,
			}
		}
	}
	return token
}

// 下一个字符
func (r *Lexer) NextChar() bool {
	// 判断是否越界
	r.Pos++
	if r.Pos >= len(r.Formula) {
		return false // eof
	}
	r.Char = r.Formula[r.Pos] // 移动[当前字符位置]
	return true
}

// 空白字符
func (r *Lexer) IsWs() bool {
	return unicode.IsSpace(rune(r.Char))
}

// 数字（小数）
func (r *Lexer) IsDigit() bool {
	return unicode.IsNumber(rune(r.Char)) || r.Char == '.'
}

// 字母
func (r *Lexer) IsLetter() bool {
	return r.Char >= 'a' && r.Char <= 'z'
}

// 变量
func (r *Lexer) IsVar() bool {
	return r.Char >= 'A' && r.Char <= 'Z'
}
