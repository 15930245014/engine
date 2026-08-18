package dsl

import (
	"errors"
	"fmt"
	"gitlab.shudieds.com/mec/lib/utils/convert"
	"gitlab.shudieds.com/zxh/engine/conf"
	"strings"
)

type Token struct {
	// raw characters
	Tok string
	// type with Literal/Operator
	Type,
	Flag int

	Offset int

	//输入变量值
	Val interface{}

	//附加
	GblExtra []string
}

func (t *Token) ToString() string {
	typeStr := ""
	switch t.Type {
	case conf.Identifier:
		typeStr = "func"
		break
	case conf.COMMA:
		typeStr = "com"
		break
	case conf.Operator:
		typeStr = "operator"
		break
	case conf.Literal:
		typeStr = "number"
		break
	case conf.VARIABLE:
		typeStr = "variable"
		break
	case conf.FIELD:
		typeStr = "field"
		break
	case conf.GLOBAL:
		typeStr = "global"
	case conf.INPUT:
		typeStr = "input"
	default:
		typeStr = "undefined"
	}
	return fmt.Sprintf("type:%s,val:%s", typeStr, t.Tok)
}

type tokenParser struct {
	Source string
	ch     byte
	offset int //当前下标
	err    error
	//转义状态机
	escapeDFA map[byte]bool

	//操作符
	op map[byte]bool
}

func GetTokens(s string) ([]Token, error) {
	p := &tokenParser{
		Source: s, //元字符串
		err:    nil,
		ch:     s[0],
		escapeDFA: map[byte]bool{
			't':  true,
			'n':  true,
			'\\': true,
			'"':  true,
			'r':  true,
		},
		op: map[byte]bool{
			'+': true,
			'-': true,
			'*': true,
			'/': true,
			'^': true,
			'%': true,
		},
	}
	toks := p.parse()
	if p.err != nil {
		return nil, p.err
	}
	return toks, nil
}

// 生成token
func (p *tokenParser) parse() []Token {
	//声明返回值
	toks := make([]Token, 0)
	for {
		tok := p.nextTok()
		if tok.Type == -1 {
			break
		}
		toks = append(toks, tok)
	}
	return toks
}

func (p *tokenParser) nextTok() Token {
	if p.offset >= len(p.Source) || p.err != nil {
		return Token{Type: -1}
	}

	//去重前面制表符
	var err error
	for p.isWhitespace(p.ch) && err == nil {
		err = p.nextCh()
	}

	//记录截取的开始
	start := p.offset
	var tok Token
	switch p.ch {
	case
		'(',
		')',
		'+',
		'*',
		'/',
		'^':
		tok = Token{
			Tok:  string(p.ch),
			Type: conf.Operator,
		}
		tok.Offset = start
		err = p.nextCh()
	case '-':
		if p.nextCh() == nil && p.ch == '>' {
			tok = Token{
				Tok:  "->",
				Type: conf.Operator,
			}
			tok.Offset = start
			err = p.nextCh()
		} else {
			tok = Token{
				Tok:  "-",
				Type: conf.Operator,
			}
			tok.Offset = start
		}
	case '%':
		if p.nextCh() == nil && p.ch == 's' {
			tok = Token{
				Tok:  "self",
				Type: conf.Identifier,
			}
			tok.Offset = start
			err = p.nextCh()
		} else {
			tok = Token{
				Tok:  "%",
				Type: conf.Operator,
			}
			tok.Offset = start

		}
	case '>':
		if p.nextCh() != nil {
			s := fmt.Sprintf("symbol error: unknown '%v', pos [%v:]\n%s",
				string(p.ch),
				start,
				ErrPos(p.Source, start))
			p.err = errors.New(s)
		} else if p.ch == '=' {
			tok = Token{
				Tok:    ">=",
				Type:   conf.Operator,
				Offset: start,
			}
			err = p.nextCh()
		} else {
			tok = Token{
				Tok:    ">",
				Type:   conf.Operator,
				Offset: start,
			}
		}
	case '<':
		if p.nextCh() != nil {
			s := fmt.Sprintf("symbol error: unknown '%v', pos [%v:]\n%s",
				string(p.ch),
				start,
				ErrPos(p.Source, start))
			p.err = errors.New(s)
		} else if p.ch == '=' {
			tok = Token{
				Tok:    "<=",
				Type:   conf.Operator,
				Offset: start,
			}
			err = p.nextCh()
		} else {
			tok = Token{
				Tok:    "<",
				Type:   conf.Operator,
				Offset: start,
			}
		}
	case '=':
		if p.nextCh() != nil || p.ch != '=' {
			s := fmt.Sprintf("symbol error: unknown '%v', pos [%v:]\n%s",
				string(p.ch),
				start,
				ErrPos(p.Source, start))
			p.err = errors.New(s)
			break
		}
		tok = Token{
			Tok:    "==",
			Type:   conf.Operator,
			Offset: start,
		}
		err = p.nextCh()
	case '&':
		if p.nextCh() == nil && p.ch == '&' {
			tok = Token{
				Tok:    "&&",
				Type:   conf.Operator,
				Offset: start,
			}
			err = p.nextCh()
		} else {
			s := fmt.Sprintf("symbol error: unknown '%v', pos [%v:]\n%s",
				string(p.ch),
				start,
				ErrPos(p.Source, start))
			p.err = errors.New(s)
		}
	case '|':
		if p.nextCh() == nil && p.ch == '|' {
			tok = Token{
				Tok:    "||",
				Type:   conf.Operator,
				Offset: start,
			}
			err = p.nextCh()
		} else {
			s := fmt.Sprintf("symbol error: unknown '%v', pos [%v:]\n%s",
				string(p.ch),
				start,
				ErrPos(p.Source, start))
			p.err = errors.New(s)
		}
	case '!':
		if p.nextCh() != nil {
			s := fmt.Sprintf("symbol error: unknown '%v', pos [%v:]\n%s",
				string(p.ch),
				start,
				ErrPos(p.Source, start))
			p.err = errors.New(s)
		} else if p.ch == '=' {
			tok = Token{
				Tok:    "!=",
				Type:   conf.Operator,
				Offset: start,
			}
			err = p.nextCh()
		} else {
			tok = Token{
				Tok:    "!",
				Type:   conf.Operator,
				Offset: start,
			}
		}
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
		for p.isDigitNum(p.ch) && p.nextCh() == nil {
			if (p.ch == '-' || p.ch == '+') && p.Source[p.offset-1] != 'e' {
				break
			}
		}
		tok = Token{
			Tok:  strings.ReplaceAll(p.Source[start:p.offset], "_", ""),
			Type: conf.Literal,
		}
		tok.Offset = start
	case ',':
		tok = Token{
			Tok:  string(p.ch),
			Type: conf.COMMA,
		}
		tok.Offset = start
		err = p.nextCh()
	case '$': //变量
		if p.nextCh() != nil || !p.isChar(p.ch) {
			s := fmt.Sprintf("symbol error: unknown '%v', pos [%v:]\n%s",
				string(p.ch),
				start,
				ErrPos(p.Source, start))
			p.err = errors.New(s)
		} else {
			for p.nextCh() == nil && (p.isWordChar(p.ch) || p.ch == '.') {
			}
			eName := p.Source[start+1 : p.offset]
			tArr := strings.Split(eName, ".")
			tok = Token{
				Tok:      tArr[0],
				GblExtra: tArr,
				Type:     conf.VARIABLE,
				Offset:   start + 1,
			}
		}
	case '@': //field
		if p.nextCh() != nil || !p.isChar(p.ch) {
			s := fmt.Sprintf("symbol error: unknown '%v', pos [%v:]\n%s",
				string(p.ch),
				start,
				ErrPos(p.Source, start))
			p.err = errors.New(s)
		} else {
			for p.nextCh() == nil && (p.isWordChar(p.ch) || p.ch == '.') {
			}
			eName := p.Source[start+1 : p.offset]
			tArr := strings.Split(eName, ".")
			tok = Token{
				Tok:      tArr[0],
				GblExtra: tArr[1:],
				Type:     conf.FIELD,
				Offset:   start + 1,
			}
		}
	case '"':
		var pre byte
		pre = ' '
		outPut := make([]byte, 0)
		for p.nextCh() == nil && (p.ch != '"' && pre != '\\') {
			if p.ch == '\\' {
				if p.nextCh() == nil && p.escapeDFA[p.ch] == true {
					outPut = append(outPut, p.ch)
				} else {
					outPut = append(outPut, '\\')
				}
			} else {
				outPut = append(outPut, p.ch)
			}
			pre = p.ch
		}

		//只有一个字符
		if p.offset >= len(p.Source) {
			s := fmt.Sprintf("symbol error: unknown '%v', pos [%v:]\n%s",
				string(p.ch),
				start,
				ErrPos(p.Source, start))
			p.err = errors.New(s)
		}

		tok = Token{
			Tok:    "",
			Type:   conf.INPUT,
			Val:    convert.BytesToString(outPut),
			Offset: start + 1,
		}
		err = p.nextCh()
	case '`': //也是字符串 不会处理转义
		for p.nextCh() == nil && p.ch != '`' {
		}

		//只有一个字符
		if p.offset >= len(p.Source) {
			s := fmt.Sprintf("symbol error: unknown '%v', pos [%v:]\n%s",
				string(p.ch),
				start,
				ErrPos(p.Source, start))
			p.err = errors.New(s)
		}

		tok = Token{
			Tok:    "",
			Type:   conf.INPUT,
			Val:    p.Source[start+1 : p.offset],
			Offset: start + 1,
		}
		err = p.nextCh()
	//case '_':
	//	if p.nextCh() != nil || !p.isChar(p.ch) {
	//		s := fmt.Sprintf("symbol error: unknown '%v', pos [%v:]\n%s",
	//			string(p.ch),
	//			start,
	//			ErrPos(p.Source, start))
	//		p.err = errors.New(s)
	//	} else {
	//		pre := p.ch
	//		for p.nextCh() == nil && (p.isWordChar(p.ch) || p.ch == '.') {
	//			if p.ch == '.' && pre == '.' {
	//				s := fmt.Sprintf("symbol error: unknown '%v', pos [%v:]\n%s",
	//					string(p.ch),
	//					p.offset,
	//					ErrPos(p.Source, start))
	//				p.err = errors.New(s)
	//			}
	//			pre = p.ch
	//		}
	//		if p.err == nil {
	//			tName := p.Source[start:p.offset]
	//			tArr := strings.Split(tName, ".")
	//
	//			if conf.SysSet[tArr[0]] == false {
	//				s := fmt.Sprintf("symbol error: unknown 全局变量'%v', pos [%v:]\n%s",
	//					tName,
	//					p.offset,
	//					ErrPos(p.Source, start))
	//				p.err = errors.New(s)
	//			} else {
	//				tok = Token{
	//					Tok:      tArr[0],
	//					GblExtra: tArr,
	//					Type:     conf.GLOBAL,
	//					Offset:   start,
	//				}
	//			}
	//		}
	//	}
	//	break
	default:
		//函数
		if p.isChar(p.ch) {
			for p.isWordChar(p.ch) && p.nextCh() == nil {
			}
			fName := p.Source[start:p.offset]

			//check 函数名称
			if tokenCheckDef[fName] == false {
				s := fmt.Sprintf("symbol error: unknown func %s '%v', pos [%v:]\n%s",
					fName,
					string(p.ch),
					start,
					ErrPos(p.Source, start))
				p.err = errors.New(s)
			} else {
				tok = Token{
					Tok:  fName,
					Type: conf.Identifier,
				}
				tok.Offset = start
			}

		} else if p.ch != ' ' { //其他的错误符号
			s := fmt.Sprintf("symbol error: unknown '%v', pos [%v:]\n%s",
				string(p.ch),
				start,
				ErrPos(p.Source, start))
			p.err = errors.New(s)
		}
	}
	return tok
}

/*
*
是否继续
*/
func (p *tokenParser) nextCh() error {
	p.offset++
	if p.offset < len(p.Source) {
		p.ch = p.Source[p.offset]
		return nil
	}
	return errors.New("EOF")
}
func (p *tokenParser) isHasNextCh() bool {
	if p.offset < len(p.Source) {
		return true
	}
	return false
}

/*
*
是否为制表符
*/
func (p *tokenParser) isWhitespace(c byte) bool {
	return c == ' ' ||
		c == '\t' ||
		c == '\n' ||
		c == '\v' ||
		c == '\f' ||
		c == '\r'
}

/*
*
是否为数字
*/
func (p *tokenParser) isDigitNum(c byte) bool {
	return '0' <= c && c <= '9' || c == '.' || c == '_' || c == 'e' || c == '-' || c == '+'
}

/*
*
是否为字母
*/
func (p *tokenParser) isChar(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

/*
*
是否为单词
*/
func (p *tokenParser) isWordChar(c byte) bool {
	return p.isChar(c) || '0' <= c && c <= '9' || c == '_'
}
