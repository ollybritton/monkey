package cmd

import (
	"fmt"

	"github.com/chzyer/readline"
	"github.com/ollybritton/monkey/lexer"
	"github.com/ollybritton/monkey/token"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// lexCmd represents the lex command
var lexCmd = &cobra.Command{
	Use:   "lex",
	Short: "Display lex output for a given input string.",
	Long:  `lex will tokenize an input string and display the output.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("monkey :: Lexical Analysis\n\n")

		rl, err := readline.New("==> ")
		if err != nil {
			panic(errors.Wrap(err, "error creating repl"))
		}
		defer rl.Close()

		for {
			line, err := rl.Readline()
			if err != nil {
				break
			}

			l := lexer.New(line)
			for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
				fmt.Printf("%+v\n", tok)
			}

			fmt.Println("")
		}
	},
}

func init() {
	rootCmd.AddCommand(lexCmd)
}
