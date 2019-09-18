panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x30 pc=0x1100974]

goroutine 1 [running]:
github.com/ollybritton/monkey/evaluator.Eval(0x12392e0, 0x0, 0xc000045bb8, 0x12392e0, 0x0)
	/Users/olly/code/go/src/github.com/ollybritton/monkey/evaluator/evaluator.go:39 +0x534
github.com/ollybritton/monkey/evaluator.evalProgram(0xc000048780, 0x1, 0x1, 0xc000045bb8, 0xc000045ab0, 0x10bcc8e)
	/Users/olly/code/go/src/github.com/ollybritton/monkey/evaluator/evaluator.go:83 +0x91
github.com/ollybritton/monkey/evaluator.Eval(0x1239360, 0xc00000e100, 0xc000045bb8, 0x1, 0x1)
	/Users/olly/code/go/src/github.com/ollybritton/monkey/evaluator/evaluator.go:26 +0x12a
github.com/ollybritton/monkey/cmd/monkey/cmd.glob..func1(0x139b6e0, 0x13ba4c0, 0x0, 0x0)
	/Users/olly/code/go/src/github.com/ollybritton/monkey/cmd/monkey/cmd/eval.go:49 +0x2cb
github.com/spf13/cobra.(*Command).execute(0x139b6e0, 0x13ba4c0, 0x0, 0x0, 0x139b6e0, 0x13ba4c0)
	/Users/olly/code/go/src/github.com/spf13/cobra/command.go:833 +0x2ae
github.com/spf13/cobra.(*Command).ExecuteC(0x139be60, 0xc000045f68, 0x1185fde, 0x139be60)
	/Users/olly/code/go/src/github.com/spf13/cobra/command.go:917 +0x2fc
github.com/spf13/cobra.(*Command).Execute(...)
	/Users/olly/code/go/src/github.com/spf13/cobra/command.go:867
github.com/ollybritton/monkey/cmd/monkey/cmd.Execute()
	/Users/olly/code/go/src/github.com/ollybritton/monkey/cmd/monkey/cmd/root.go:23 +0x32
main.main()
	/Users/olly/code/go/src/github.com/ollybritton/monkey/cmd/monkey/main.go:21 +0x20
exit status 2
