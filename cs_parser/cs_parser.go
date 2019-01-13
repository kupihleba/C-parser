package cs_parser

import (
	"fmt"
	"parser/cs_lexer"
	"strconv"
)

type State int

const (
	GO State = iota
	BACKTRACK
	SUCCESS
	FAILURE
)

func (s State) String() string {
	switch s {
	case GO:
		return "Q"
	case BACKTRACK:
		return "B"
	case SUCCESS:
		return "SUCCESS"
	case FAILURE:
		return "FAILURE"
	default:
		return "???"
	}
}

var state State
var pos int // Next token pos
var l1 []Instance
var l2 []Action
var Tokens []int
var Grammar []Rule
var ExpToRule map[int]*Rule
var Expressions []Expression

func init() {
	initGrammar()
	//initSampleGrammar()
	initExpressions()
}

func reset() {
	state = GO
	pos = 0
	Tokens = Tokens[:0]
	l1 = l1[:0]
	l2 = l2[:0]
}

type Rule struct {
	name    string
	exprSet []Expression
}

func (r Rule) String() string {
	return r.name
}

type Expression []Instance

type Command int

const (
	READ_TOKEN Command = iota
	RULE_MATCHED
)

type Action struct {
	value int
	cmd   Command
}

func (a Action) String() string {
	if a.cmd == READ_TOKEN {
		return "-"
	} else {
		return strconv.Itoa(a.value)
	}
}

//func (this *Rule) matchSuffix(pattern []Instance) bool {
//	for _, expr := range this.exprSet {
//		if expr.matchSuffix(pattern) {
//			return true
//		}
//	}
//	return false
//}

func (this *Expression) matchSuffix(pattern []Instance) bool {
	if len(pattern) < len(*this) {
		return false
	}
	suffix := pattern[len(pattern)-len(*this):] // get suffix of expression length
	for i, s := range *this {
		if !equals(suffix[i], s) {
			//fmt.Printf("%v != %v\n", suffix[i], s)

			return false
		} else {
			//fmt.Printf("%v == %v\n", suffix[i], s)
		}
	}
	return true
}

func equals(a Instance, b Instance) bool {
	if a.t == b.t {
		if a.t == TERMINAL {
			return a.id_or_ref.(int) == b.id_or_ref.(int)
		} else if a.t == NONTERMINAL {
			return a.id_or_ref.(*Rule).name == b.id_or_ref.(*Rule).name
		}
	}
	return false
}

func initGrammar() {
	var CLASS_DECL, CLASS_BODY, USING, VAR_DECL, IDENT_CHAIN, ASSIGN, EXPRESSION, PROGRAM, METHOD, PREFIX, TYPE, NAMESPACE Rule
	CLASS_BODY.name = "CLASS BODY"
	CLASS_DECL.name = "CLASS DECLARATION"
	USING.name = "USING KEYWORD"
	VAR_DECL.name = "VARIABLE DECLARATION"
	IDENT_CHAIN.name = "CHAIN OF IDENTIFICATORS"
	ASSIGN.name = "ASSIGN EXPRESSION"
	EXPRESSION.name = "EXPRESSION"
	PROGRAM.name = "PROGRAM"
	METHOD.name = "METHOD DECLARATION"
	PREFIX.name = "PREFIX KEYWORD"
	NAMESPACE.name = "NAMESPACE"

	PROGRAM.exprSet = []Expression{
		{ref(&USING)},
		{ref(&CLASS_DECL)},
		{ref(&METHOD)},
		{ref(&NAMESPACE)},
		{ref(&VAR_DECL)},
		{ref(&PROGRAM), ref(&USING)},
		{ref(&PROGRAM), ref(&CLASS_DECL)},
		{ref(&PROGRAM), ref(&METHOD)},
		{ref(&PROGRAM), ref(&NAMESPACE)},
		{ref(&PROGRAM), ref(&VAR_DECL)},
	}
	NAMESPACE.exprSet = []Expression{
		{tok("namespace"), tok("IDENTIFIER"), tok("{"), tok("}")},
		{tok("namespace"), tok("IDENTIFIER"), tok("{"), ref(&PROGRAM), tok("}")},
	}
	USING.exprSet = []Expression{
		{tok("using"), tok("IDENTIFIER"), tok(";")},
	}
	CLASS_DECL.exprSet = []Expression{
		{tok("class"), tok("IDENTIFIER"), tok("{"), ref(&CLASS_BODY), tok("}")},
		{tok("class"), tok("IDENTIFIER"), tok("{"), tok("}")},
	}
	ASSIGN.exprSet = []Expression{
		{tok("IDENTIFIER"), tok("="), ref(&EXPRESSION)},
	}
	EXPRESSION.exprSet = []Expression{
		{tok("IDENTIFIER")},
		{tok("NUMBER")},
		{tok("("), ref(&EXPRESSION), tok(")")},
	}
	PREFIX.exprSet = []Expression{
		{tok("static")},
	}
	CLASS_BODY.exprSet = []Expression{
		{ref(&VAR_DECL)},
		{ref(&METHOD)},
		{ref(&METHOD), ref(&CLASS_BODY)},
		{ref(&VAR_DECL), ref(&CLASS_BODY)},
	}
	METHOD.exprSet = []Expression{
		{ref(&PREFIX), ref(&TYPE), tok("IDENTIFIER"), tok("("), ref(&VAR_DECL), tok(")"), tok("{"), tok("}")},
		{ref(&PREFIX), ref(&TYPE), tok("IDENTIFIER"), tok("("), ref(&VAR_DECL), tok(")"), tok(";")},

		{ref(&PREFIX), ref(&TYPE), tok("IDENTIFIER"), tok("("), tok(")"), tok("{"), tok("}")},
		{ref(&PREFIX), ref(&TYPE), tok("IDENTIFIER"), tok("("), tok(")"), tok(";")},

		{ref(&TYPE), tok("IDENTIFIER"), tok("("), ref(&VAR_DECL), tok(")"), tok("{"), tok("}")},
		{ref(&TYPE), tok("IDENTIFIER"), tok("("), ref(&VAR_DECL), tok(")"), tok(";")},

		{ref(&TYPE), tok("IDENTIFIER"), tok("("), tok(")"), tok("{"), tok("}")},
		{ref(&TYPE), tok("IDENTIFIER"), tok("("), tok(")"), tok(";")},
	}
	TYPE.exprSet = []Expression{
		{tok("string")},
		{tok("int")},
		{tok("void")},
		{tok("IDENTIFIER")},
	}
	VAR_DECL.exprSet = []Expression{
		{ref(&TYPE), ref(&IDENT_CHAIN), tok(";")},
	}
	IDENT_CHAIN.exprSet = []Expression{
		{tok("IDENTIFIER")},
		{tok("IDENTIFIER"), tok(","), ref(&IDENT_CHAIN)},
		{ref(&ASSIGN)},
		{ref(&ASSIGN), tok(","), ref(&IDENT_CHAIN)},
	}
	fmt.Printf("%v", Grammar)
	Grammar = append(Grammar, PROGRAM, CLASS_DECL, CLASS_BODY, USING, VAR_DECL, METHOD, PREFIX, TYPE, IDENT_CHAIN, ASSIGN, EXPRESSION, NAMESPACE)
}

func initSampleGrammar() {
	var S, A Rule
	S.name = "S"
	A.name = "A"
	S.exprSet = []Expression{
		{ref(&A), ref(&S)},
		{tok("a")},
	}
	A.exprSet = []Expression{
		{tok("b"), ref(&S), ref(&A)},
		{tok("b")},
	}
	Grammar = append(Grammar, S, A)
}

//func expression(context *Rule, instance ...Instance) {
//
//}

func initExpressions() {
	ExpToRule = make(map[int]*Rule)
	c := 0
	for i := 0; i < len(Grammar); i++ {
		Expressions = append(Expressions, Grammar[i].exprSet...)
		for j := 0; j < len(Grammar[i].exprSet); j++ {
			ExpToRule[c] = &Grammar[i]
			c++
		}
	}
}

func ref(ref *Rule) Instance {
	var i Instance
	i.t = NONTERMINAL
	i.id_or_ref = ref
	return i
}

func tok(token_type string) Instance {
	return tokID(cs_lexer.TokenIdentifiers[token_type])
}

func tokID(id int) Instance {
	var i Instance
	i.t = TERMINAL
	i.id_or_ref = id
	return i
}

type Type int

const (
	TERMINAL Type = iota
	NONTERMINAL
)

type Instance struct {
	id_or_ref interface{}
	//id int
	//ref *Rule
	t Type
}

func (i Instance) String() string {
	if i.nonTerminal() {
		return i.id_or_ref.(*Rule).name
	}
	return cs_lexer.TokToStr(i.id_or_ref.(int))
}

func (this *Instance) nonTerminal() bool {
	return this.t == NONTERMINAL
}
func (this *Instance) isTerminal() bool {
	return this.t == TERMINAL
}

func reachedTheEnd() bool {
	return pos == len(Tokens)
}

func findSuitableExpression() (int, error) {
	if len(l1) == 0 {
		return -1, fmt.Errorf("No expressions supplied")
	}

	for i, expr := range Expressions {
		if expr.matchSuffix(l1) {
			return i, nil
		}
	}
	return -1, fmt.Errorf("No matchSuffix")
}

func findSuitableExpressionAfter(index int) (int, error) {
	if len(l1) == 0 {
		return -1, fmt.Errorf("No expressions supplied")
	}

	for i := index + 1; i < len(Expressions); i++ {
		if Expressions[i].matchSuffix(l1) {
			return i, nil
		}
	}
	return -1, fmt.Errorf("No matchSuffix")
}

func printState(step int) {
	fmt.Printf("%d:\t(%s, %d, %v, %v)\n", step, state, pos, l1, l2)
}

func readToken() {
	l1 = append(l1, tokID(Tokens[pos]))
	l2 = append(l2, Action{pos, READ_TOKEN})
	pos++
}
func throwToken() {
	if pos <= 0 || l2[len(l2)-1].cmd != READ_TOKEN || l1[len(l1)-1].t != TERMINAL {
		panic("throw token failure")
	}
	l1 = l1[:len(l1)-1]
	l2 = l2[:len(l2)-1]
	pos--
}

func goExpression(expr_n int) {
	l1 = l1[:len(l1)-len(Expressions[expr_n])] // pop last elements, that matched rule
	l1 = append(l1, ref(ExpToRule[expr_n]))
	l2 = append(l2, Action{expr_n, RULE_MATCHED}) // l2 push
}

func rollback() {
	last_reduce := l2[len(l2)-1] // get previous reduce expression
	if last_reduce.cmd != RULE_MATCHED {
		panic("call for rollback on terminal")
	}

	l1 = l1[:len(l1)-1]                                // pop l1 rule
	l1 = append(l1, Expressions[last_reduce.value]...) // push old expressions back

	l2 = l2[:len(l2)-1]
}

func step() {
	switch state {
	case GO:
		if reachedTheEnd() &&
			len(l1) == 1 &&
			l1[0].nonTerminal() {
			state = SUCCESS
			return
		}

		expr, err := findSuitableExpression() // Trying to find suitable expr
		if err != nil {                       // if there was no matchSuffix
			//println("No matchSuffix!")
			if pos < len(Tokens) { // if we can read more tokens
				readToken()
			} else {
				state = BACKTRACK
			}
		} else { // if we found a suitable expr
			//println("MATCHED RULE:", ExpToRule[expr].name)
			goExpression(expr)
		}

		break
	case BACKTRACK:
		if len(l1) == 0 {
			state = FAILURE
			return
		}

		if l1[len(l1)-1].isTerminal() { // if last element is a terminal
			throwToken() // throw it away
			//println("Throw token!")
			return
		}

		// The last element is non terminal

		lastExpr := l2[len(l2)-1]
		if lastExpr.cmd == READ_TOKEN {
			panic("UNEXPECTED READ_TOKEN")
		}

		l1 = l1[:len(l1)-1]                             // pop l1 rule
		l1 = append(l1, Expressions[lastExpr.value]...) // push old expressions back

		l2 = l2[:len(l2)-1]
		expr, err := findSuitableExpressionAfter(lastExpr.value)

		if err == nil { // if alternative rule exists
			l1 = l1[:len(l1)-len(Expressions[expr])] // pop last elements, that matched rule
			l1 = append(l1, ref(ExpToRule[expr]))

			//l2 = l2[:len(l2)-1] // l2.pop()
			l2 = append(l2, Action{expr, RULE_MATCHED}) // l2 push
			state = GO
		} else {
			// ROLLBACK TEMP CHANGES
			if pos < len(Tokens) {
				readToken()
				state = GO
			}
		}

	case SUCCESS:
		println("FILE HAS BEEN SUCCESSFULLY PARSED")
		return

	case FAILURE:
		println("ERRORS WERE ENCOUNTERED WHILE PARSING")
		return

	}
}

const PARSE_LIMIT = 1000000

func Parse(tokens []int) {
	reset()
	Tokens = tokens
	for i := 1; state != SUCCESS && state != FAILURE; i++ {
		step()
		printState(i)
	}
}
