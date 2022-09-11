package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp/syntax"
	"time"
)

//	Maximum repitions to use for unlimited matches.
var maxRepititions = 32

func main() {

	//	Parse the regex argument.
	exp, err := syntax.Parse(os.Args[1], syntax.Perl)
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer

	rand.Seed(time.Now().UnixNano())
	if err := run(&b, exp); err != nil {
		panic(err)
	}

	fmt.Println(b.String())
}

func run(writer *bytes.Buffer, expression *syntax.Regexp) (err error) {

	//	Perform operations subject to operators.
	//
	//	I followed the list of operators from:
	//	https://pkg.go.dev/regexp/syntax
	switch expression.Op {

	//	Matches no string
	case syntax.OpNoMatch:
		return

	//	Matches empty string
	case syntax.OpEmptyMatch:
		return

	//	Matches runes sequence
	case syntax.OpLiteral:

		//	"expression.Rune" contains matched literals for OpCharClass and OpLiteral.
		writer.WriteString(string(expression.Rune))

	//	Matches runes interpreted as range pair list.
	case syntax.OpCharClass:

		sum := 0

		//	"expression.Rune" contains matched literals for OpCharClass and OpLiteral.
		for i := 0; i < len(expression.Rune); i += 2 {
			sum += 1 + int(expression.Rune[i+1]-expression.Rune[i])
		}

		nRune := rune(rand.Int63n(int64(sum)))
		for i := 0; i < len(expression.Sub); i += 2 {
			min, max := expression.Rune[i], expression.Rune[i+1]
			delta := max - min
			if nRune <= delta {
				writer.WriteRune(min + nRune)
				return nil
			}
			nRune -= 1 + delta
		}

	//	Matches any character except newline
	case syntax.OpAnyCharNotNL:

		//	New line rune literal `\n` is number 95.
		writer.WriteRune(rune(' ' + rand.Int63n(95)))

	//	Matches any character
	case syntax.OpAnyChar:
		i := rand.Int63n(96)
		char := rune(' ' + i)

		//	Replace the newline
		if i == 96 {
			char = '\n'
		}

		writer.WriteRune(char)

	//	Matches empty string at the beginning of a line
	case syntax.OpBeginLine:
		if writer.Len() != 0 {
			writer.WriteByte('\n')
		}

	//	Matches empty string at the end of a line
	case syntax.OpEndLine:
		if writer.Len() != 0 {
			writer.WriteByte('\n')
		} else {
			return io.EOF
		}

	//	Matches empty string at the beginning of a text
	case syntax.OpBeginText:
		return io.EOF

	//	Matches empty string at the end of a text
	case syntax.OpEndText:
		return io.EOF

	//	Matches word boundary '\b'
	case syntax.OpWordBoundary:
		//	TODO : Not supported yet
		return errors.New("word boundary not supported")

	//	Matches word non-boundary '\b'
	case syntax.OpNoWordBoundary:
		//	TODO : Not supported yet
		return errors.New("word non-boundary not supported")

	//	Capturing subexpression with index Cap, optional name Name
	case syntax.OpCapture:

		//	Loop over all the subexpressions.
		for _, subexp := range expression.Sub {
			if err := run(writer, subexp); err != nil {
				return err
			}
		}

	//	Matches concatenation of sub expressions.
	case syntax.OpConcat:

		//	Loop over all the subexpressions.
		for _, subexp := range expression.Sub {
			if err := run(writer, subexp); err != nil {
				return err
			}
		}

	//	Matches alternation of sub expressions.
	case syntax.OpAlternate:
		index := rand.Int63n(int64(len(expression.Sub)))
		return run(writer, expression.Sub[index])

	//	Matches Sub[0] zero or more times.
	case syntax.OpStar:
		min := 0
		max := maxRepititions

		//	Run the loop until value is lower than or equal to zero,
		//	and keep decrementing the value after every  iteration.
		for item := min + int(rand.Int63n(int64(max)-int64(min)+1)); item > 0; item-- {

			//	Loop over all the subexpressions.
			for _, subexp := range expression.Sub {
				if err := run(writer, subexp); err != nil {
					return err
				}
			}
		}

	//	Matches Sub[0] one or more times.
	case syntax.OpPlus:
		min := 1
		max := maxRepititions

		//	Run the loop until value is lower than or equal to zero,
		//	and keep decrementing the value after every  iteration.
		for item := min + int(rand.Int63n(int64(max)-int64(min)+1)); item > 0; item-- {

			//	Loop over all the subexpressions.
			for _, subexp := range expression.Sub {
				if err := run(writer, subexp); err != nil {
					return err
				}
			}
		}

	// 	Matches Sub[0] at least Min times, at most Max (Max == -1 is no limit).
	case syntax.OpRepeat:
		min := expression.Min
		max := expression.Max

		//	Replace no limit.
		if max == -1 {
			max = min + maxRepititions
		}

		//	Run the loop until value is lower than or equal to zero,
		//	and keep decrementing the value after every  iteration.
		for item := min + int(rand.Int63n(int64(max)-int64(min)+1)); item > 0; item-- {

			//	Loop over all the subexpressions.
			for _, subexp := range expression.Sub {
				if err := run(writer, subexp); err != nil {
					return err
				}
			}
		}

	// 	Matches Sub[0] zero or one times.
	case syntax.OpQuest:

		if rand.Int63n(0xFFFFFFFF) > rand.Int63n(0x7FFFFFFF) {

			//	Loop over all the subexpressions.
			for _, subexp := range expression.Sub {
				if err := run(writer, subexp); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
