package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

//position is the token position in the source text,
//used for error tracing,line and col counts from 0.
type position struct {
	line int //line number
	col  int //colummn number
}

type TokenType int

//lexical token.
type token struct {
	lit string    //literal value
	typ TokenType //token type
	pos position  //postion
}

//for debug
func (t token) String() string {
	return fmt.Sprintf("<lit:\"%s\",typ:%s,pos line:%d,pos col:%d>", t.lit, tokens[t.typ], t.pos.line, t.pos.col)
}

type stateFn func(l *lexer) stateFn

//lexical parser.
type lexer struct {
	cur token //current scanned token

	src string //source

	pos   int //current scanning index
	start int //start scanning index
	width int //width of string scanned

	lineNum int //line counter,counts from 0
	colNum  int //columnn counter,counts from 1

	errors    []string   //errors stack
	state     stateFn    //state function
	tokenChan chan token //token channel
}

const (
	token_begin TokenType = iota
	tNUM                  // number -?digit*.digit*[E|e]-?digit*

	tPLUS  // +
	tMUNIS // -

	tMUTIL // *
	tDIV   // /

	tEOF // eof
	token_end
)

const eof = -1

var tokens = map[TokenType]string{
	tMUNIS: "[-]",
	tDIV:   "[/]",
	tPLUS:  "[+]",
	tMUTIL: "[*]",
	tEOF:   "[EndOfFile]",
}

//newLexer initiates token channel and go run lexer.run and return lexer.
func newLexer(src string) *lexer {
	l := &lexer{src: src,
		tokenChan: make(chan token),
	}
	go l.run()
	return l
}

//scan next token
func (l *lexer) next() rune {
	if l.pos >= len(l.src) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.src[l.pos:])
	l.width = w
	l.pos += l.width
	l.colNum += l.width
	return r
}

//push error message into tracing stack
func (l *lexer) err(e string) {
	l.errors = append(l.errors, e)
}

//format error
func (l *lexer) errf(f string, v ...interface{}) {
	l.errors = append(l.errors, fmt.Sprintf(f, v...))
}

//backup a token
func (l *lexer) backup() {
	l.pos -= l.width
	l.colNum -= l.width
}

//peek a token
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

//ignore a token
func (l *lexer) ignore() {
	l.start = l.pos
}

//emit token
func (l *lexer) emit(typ TokenType) {
	var t = token{
		lit: l.src[l.start:l.pos],
		typ: typ,
		pos: position{l.lineNum, l.colNum - (l.pos - l.start)},
	}
	l.tokenChan <- t
	if t.typ == tEOF {
		close(l.tokenChan)
	}

	l.start = l.pos

}

//read token
func (l *lexer) token() token {
	token := <-l.tokenChan
	l.cur = token
	return token
}

//main consumer routine
func (l *lexer) run() {
	for l.state = lexBegin; l.state != nil; {
		l.state = l.state(l)
	}
}

//------------------------------------sate function----------------------------------
//error handling
//TODO:sync errors.
func lexError(l *lexer) stateFn {
	//premature lexical scanning
	//不emit接收方就一直在等待
	l.emit(tEOF)
	return nil
}

func lexUnkown(l *lexer) stateFn {
	//premature lexical scanning
	l.emit(tEOF)
	return nil
}

func lexNum(l *lexer) stateFn {
	r := l.next()
	var hasMatissa = false
	var hasFraction = false
	//optional '-',must be fllowed by at least  digit
	if r == '-' {
		r = l.next()
		if !unicode.IsDigit(r) {
			l.errf("expect decimal after '-',found %c at line %d,column %d", r, l.lineNum+1, l.colNum)
			return lexError
		}
	}
	//scan matissa
	if r != '.' { //optional digit*
		for unicode.IsDigit(r) {
			hasMatissa = true
			r = l.next()
		}
	}
	//scan fraction
	if r == '.' { //optional .digit*
		r = l.next()
		for unicode.IsDigit(r) {
			hasFraction = true
			r = l.next()
		}
	}
	ln := l.lineNum
	cn := l.colNum + 1
	if r == 'e' || r == 'E' { //optional [E|e]-?digit*
		r = l.next()
		if r == '-' {
			r = l.next()
		}
		//expotiona must be fllowed by at least one decimal digit
		if !unicode.IsDigit(r) {
			l.errf("expect decimal after expotion c at line %d,column %d", r, l.lineNum+1, l.colNum)
			return lexError
		}
		for unicode.IsDigit(r) {
			r = l.next()
		}
	}

	if !hasMatissa && !hasFraction {
		l.errf("expect decimal at line %d,column %d", ln, cn)
	}
	if r != eof {
		l.backup()
	}
	l.emit(tNUM)
	return lexBegin
}

//end of scanning
func lexEOF(l *lexer) stateFn {
	l.emit(tEOF)
	return nil
}

//main lex entry
func lexBegin(l *lexer) stateFn {
	switch r := l.next(); {
	case unicode.IsDigit(r) || r == '.' || r == '-':
		r2 := l.peek()
		if r == '-' &&
			r2 != '.' && !unicode.IsDigit(r2) && r2 != 'E' && r2 != 'e' {
			goto FL
		}
		l.backup()
		return lexNum
	FL:
		fallthrough
	case r == '-':

	case r == ' ':
		l.ignore()
	case r == '\n':
		l.ignore()
		l.lineNum++
		l.colNum = 0
		//l.emit(tNEWLINE),currently not neccesary in parsing.
	case r == eof:
		return lexEOF
	default:
		l.errf("unkown char '%c' at line %d,column %d", r, l.lineNum+1, l.colNum)
		return lexUnkown
	}
	return lexBegin
}

func main() {
	println(`In this task we will focus on a lexer implementation,and it's concurrency part.
lexer is a lexical scanner that consumes source code and produce meaningful tokens.With this tokens we can then 
complete a small calculator.Our simple lexer just need to scann several tokens '+','-','*','\',and number.
Lexer is a typical produer-consumer pattern,so we need a channel to send token.
Instead of switch,we use sate function instead,in order to skip the case statements.
Now edit main.go and finish the task.`)
}
