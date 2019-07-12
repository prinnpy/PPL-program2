//Name: Prinn Prinyanut
//CSCI 3451
//Program 2

package main

import (
	"bufio"
	"errors"
	"math"
	"os"
	"fmt"
	"strconv"
	"strings"
)

// Cell is for defining each value
type Cell struct {
	value    string
	typ      string
	op       string
	priority int // + - x / ( ) 0 0 1 1 2 2
}

// Stack is for creating it
type Stack struct {
	values []*Cell
}

// Push is for pushing items into stack
func (s *Stack) Push(c *Cell) {
	s.values = append(s.values, c)
}

// Pop is for removing items from the stack
func (s *Stack) Pop() *Cell {
	if len(s.values) == 0 {
		return nil
	}
	top := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	return top
}

// Top is to get the value at the top of the stack
func (s Stack) Top() *Cell {
	if len(s.values) == 0 {
		return nil
	}
	return s.values[len(s.values)-1]
}

//Calculate is a struct for calculator
type Cal struct {
	stack        *Stack
	opStack      *Stack
	_queue       []*Cell
	postfixQueue []*Cell
}

// NewCal use to create operator and operand stacks
func NewCal() *Cal {
	return &Cal{
		stack:   &Stack{},
		opStack: &Stack{},
	}
}

func toNumber(n string) float64 {
	iN, err := strconv.ParseFloat(n, 64)
	if err != nil {
		panic(err)
	}
	return iN
}

func toString(iN float64) string {
	return strconv.FormatFloat(iN, 'f', -1, 64)
}

func (c Cal) isNumber(char string) bool {
	switch char {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ".":
		return true
	}
	return false
}

// Operations is use to figure out which operation are being used
var Operations = map[string]func(string, string) string{
	"ADD": func(a, b string) string {
		return toString(toNumber(a) + toNumber(b))
	},
	"MIN": func(a, b string) string {
		return toString(toNumber(a) - toNumber(b))
	},
	"MUL": func(a, b string) string {
		return toString(toNumber(a) * toNumber(b))
	},
	"DIV": func(a, b string) string {
		if b == "0" {
			if strings.Contains(a, ".") {
				fmt.Print("+Inf")
			} else {
				fmt.Print("illegal expression: runtime error: integer divide by zeros")
			}
			panic(a + " can not divided by 0.")
		}
		return toString(toNumber(a) / toNumber(b))
	},
	"POW": func(a, b string) string {
		if b == "0" {
			return "1"
		}
		return toString(math.Pow(toNumber(a), toNumber(b)))
	},
}

// OperatorMap determine which operation
var OperatorMap = map[string]string{
	"+": "ADD",
	"-": "MIN",
	"/": "DIV",
	"*": "MUL",
	"^": "POW",
}


func (c *Cal) prepare(expr string) {
	splits := strings.Split(expr, "")
	count := len(splits)
	num := ""
	//group := false
	//subExpr := ""
	for i := 0; i < count; i++ {
		char := splits[i]
		if char == "" {
			continue
		}
		switch char {
		case "^":
			num = ""
			c._queue = append(c._queue, &Cell{
				value:    char,
				typ:      "OP",
				op:       OperatorMap[char],
				priority: 2,
			})
		case "(", ")":
			num = ""
			c._queue = append(c._queue, &Cell{
				value:    char,
				typ:      "OP",
				priority: 2,
			})
		case "+", "-":
			num = ""
			c._queue = append(c._queue, &Cell{
				value:    char,
				typ:      "OP",
				op:       OperatorMap[char],
				priority: 0,
			})
		case "*", "/":
			num = ""
			c._queue = append(c._queue, &Cell{
				value:    char,
				typ:      "OP",
				op:       OperatorMap[char],
				priority: 1,
			})
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ".":
			num += char
			if (i+1 < count && !c.isNumber(splits[i+1])) || i == count-1 {
				c._queue = append(c._queue, &Cell{
					value:    num,
					typ:      "NUMBER",
					priority: 0,
				})
			}
		}
	}
}

func (c *Cal) postfixExpr() string {
	expr := ""
	for _, cell := range c.postfixQueue {
		expr += cell.value + " "
	}
	return expr
}

// ResetStack is for reseting the calculator to be zero every time after calculation
func (c *Cal) ResetStack() {
	c.opStack = &Stack{}
	c.stack = &Stack{}
	c._queue = []*Cell{}
	c.postfixQueue = []*Cell{}
}

func (c *Cal) postfix(expr string) string {
	c.prepare(expr)
	for _, cell := range c._queue {
		if cell.typ == "NUMBER" || cell.typ == "EXPR" {
			c.postfixQueue = append(c.postfixQueue, cell)
		} else if cell.typ == "OP" {
			if cell.value == "(" {
				c.opStack.Push(cell)
			} else if cell.value == ")" {
				for top := c.opStack.Pop(); top != nil && top.value != "("; {
					c.postfixQueue = append(c.postfixQueue, top)
					top = c.opStack.Pop()
				}
			} else {
				for top := c.opStack.Top(); top != nil && top.priority >= cell.priority && top.value != "("; {
					c.postfixQueue = append(c.postfixQueue, top)
					c.opStack.Pop() //remove
					top = c.opStack.Top()
				}
				c.opStack.Push(cell)
			}
		}
	}
	for top := c.opStack.Pop(); top != nil; {
		c.postfixQueue = append(c.postfixQueue, top)
		top = c.opStack.Pop()
	}
	return c.postfixExpr()
}

// GetPostfixExpr is to geting the postfix
func (c Cal) GetPostfixExpr(expr string) string {
	expr = strings.Trim(expr, " ")
	if len(expr) == 0 {
		return ""
	}
	postfixExpr := c.postfix(expr)
	c.ResetStack()
	return postfixExpr
}

// Calculate is the main driver for calculations
func (c Cal) Calculate(expr string) (string, error) {
	var Err error
	var res *Cell
	func() {
		defer func() {
			if err := recover(); err != nil {
				Err = errors.New(err.(string))
			}
		}()

		expr = strings.Trim(expr, " ")
		if len(expr) == 0 {
			fmt.Print("given expr is empty")
			return
		}
		c.postfix(expr)
		for _, cell := range c.postfixQueue {
			if cell.typ == "NUMBER" {
				c.stack.Push(cell)
			} else if cell.typ == "OP" {
				fn, ok := Operations[cell.op]
				if !ok {
					fmt.Print("illegal expression: operand stack underflow")
					panic("not support op " + cell.value)
				}
				b := c.stack.Pop()
				a := c.stack.Pop()
				if b == nil {
					fmt.Print("illegal expression: operand stack underflow")
					panic("Invalid number B")
				}
				if a == nil {
					fmt.Print("illegal expression: operand stack underflow")
					panic("Invalid number A")
				}
				c.stack.Push(&Cell{
					value: fn(a.value, b.value),
					typ:   "NUMBER",
				})
			}
		}
		res = c.stack.Pop()
		if res == nil {
			fmt.Print("Calculate fail!")
			panic("Calculate fail!")
		}
		c.ResetStack()
	}()
	result := ""
	if res != nil {
		result = res.value
	}
	return result, Err
}

// DoCalculation get the results
func (c Cal) DoCalculation(expr string) string {
	res, err := c.Calculate(expr)
	if err != nil {
		return ""
	}
	return res
}

func main() {
	// Make a scanner to read lines from standard input
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("> calculator")
	// Process each of the lines from standard input
	for scanner.Scan() {
		// Get the current line of text.
		line := scanner.Text()
		// Evaluate the expression and print the result
		cal := NewCal()
		fmt.Println(cal.DoCalculation(line))
	}
}