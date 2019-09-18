package cmd

import (
	"fmt"

	"github.com/ollybritton/monkey/object"

	"github.com/chzyer/readline"
	"github.com/ollybritton/monkey/evaluator"
	"github.com/ollybritton/monkey/lexer"
	"github.com/ollybritton/monkey/parser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// replCmd represents the repl command
var replCmd = &cobra.Command{
	Use:   "eval",
	Short: "Evaluate a given string",
	Long:  `Lex, parse and then evaluate a given string.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("monkey :: Evaluation\n\n")

		rl, err := readline.New("==> ")
		if err != nil {
			panic(errors.Wrap(err, "error creating repl"))
		}
		defer rl.Close()
		env := object.NewEnvironment()

		for {
			line, err := rl.Readline()
			if err != nil {
				break
			}

			l := lexer.New(line)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				for _, e := range p.Errors() {
					fmt.Println("\t", e)
				}

				fmt.Println("")
			}

			evaluated := evaluator.Eval(program, env)

			if evaluated != nil {
				fmt.Println(evaluated.Inspect())
			}

			fmt.Println("")
		}
	},
}

func init() {
	rootCmd.AddCommand(replCmd)
}
