package cmd

import (
	"fmt"

	"github.com/chzyer/readline"
	"github.com/ollybritton/monkey/lexer"
	"github.com/ollybritton/monkey/parser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("monkey :: Parser\n\n")

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
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				for _, msg := range p.Errors() {
					fmt.Println("\t" + msg)
				}
			}

			fmt.Println(program.String())
			fmt.Println("")
		}
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// parseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
