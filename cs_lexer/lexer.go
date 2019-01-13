package cs_lexer

import (
	lex "github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
	"strings"
)

// Called at package initialization. Creates the lexer and populates token lists.
func init() {
	initTokens()
	var err error
	Lexer, err = initLexer()
	if err != nil {
		panic(err)
	}
}

var Literals []string               // The tokens representing literal strings
var Keywords []string               // The keyword tokens
var Tokens []string                 // All of the tokens (including literals and keywords)
var TokenIdentifiers map[string]int // A map from the token names to their int ids
var Lexer *lex.Lexer                // The lexer object. Use this to construct a Scanner

func initTokens() {
	Literals = []string{
		"[",
		"]",
		"{",
		"}",
		"=",
		",",
		".",
		":",
		";",
		"(",
		")",
		"-",
		"+",
		//"a",
		//"b",
	}
	Keywords = []string{
		"using",
		"namespace",
		"class",
		"static",
		"void",
		// TYPES
		"int",
		"string",
		"new",
	}
	Tokens = []string{
		"IDENTIFIER",
		"COMMENT",
		"NUMBER",
	}
	Tokens = append(Tokens, Keywords...)
	Tokens = append(Tokens, Literals...)
	TokenIdentifiers = make(map[string]int)
	for i, tok := range Tokens {
		TokenIdentifiers[tok] = i
	}
}

func TokToStr(id int) string {
	return Tokens[id]
}

func token(name string) lex.Action {
	return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
		return s.Token(TokenIdentifiers[name], string(m.Bytes), m), nil
	}
}
func skip(*lex.Scanner, *machines.Match) (interface{}, error) {
	return nil, nil
}

func initLexer() (*lex.Lexer, error) {
	lexer := lex.NewLexer()

	for _, literal := range Literals {
		r := "\\" + strings.Join(strings.Split(literal, ""), "\\")
		lexer.Add([]byte(r), token(literal))
	}
	for _, keyword := range Keywords {
		lexer.Add([]byte(strings.ToLower(keyword)), token(keyword))
	}

	lexer.Add([]byte("[0-9]+"), token("NUMBER"))
	lexer.Add([]byte(`//[^\n]*\n?`), skip)                             // SKIP COMMENTS
	lexer.Add([]byte(`/\*([^*]|\r|\n|(\*+([^*/]|\r|\n)))*\*+/`), skip) // SKIP COMMENTS
	lexer.Add([]byte(`([a-z]|[A-Z])([a-z]|[A-Z]|[0-9]|_)*`), token("ID"))
	lexer.Add([]byte(`"([^\\"]|(\\.))*"`), token("ID"))
	lexer.Add([]byte("( |\t|\n|\r)+"), skip)

	err := lexer.Compile()
	if err != nil {
		return nil, err
	}
	return lexer, nil
}
