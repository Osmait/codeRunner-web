package coderunner

import (
	"fmt"

	programinglanguages "github.com/Osmait/CodeRunner-web/internal/modules/programingLanguages"
	"github.com/Osmait/CodeRunner-web/internal/modules/runner"
)

type CodeRunner struct {
	Runner *runner.Runner
}

func NewCodeRunner(runner *runner.Runner) *CodeRunner {
	return &CodeRunner{
		Runner: runner,
	}
}

func (c *CodeRunner) RunCode() {
	lang := programinglanguages.NewPrograminLanguages("python", "py", "python")
	code := "2+2"
	cn := make(chan []byte)

	go func() {
		for msg := range cn {
			fmt.Println(string(msg))
		}
	}()

	c.Runner.Execute(code, lang.GetName(), cn)
}

func (c *CodeRunner) Run() error {
	return nil
}
