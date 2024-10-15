package coderunner

import (
	"github.com/Osmait/CodeRunner-web/internal/modules/dispacher"
	programinglanguages "github.com/Osmait/CodeRunner-web/internal/modules/programingLanguages"
	"github.com/Osmait/CodeRunner-web/internal/modules/runner"
)

type CodeRunner struct {
	Runner  runner.RunnerInterface
	outputs *dispacher.Notifier
}

func NewCodeRunner(runner runner.RunnerInterface, outputs *dispacher.Notifier) *CodeRunner {
	return &CodeRunner{
		Runner:  runner,
		outputs: outputs,
	}
}

func (c *CodeRunner) RunCode() {
	lang := programinglanguages.NewPrograminLanguages("python", "py", "python")
	code := "2+2"
	c.Runner.Execute(code, lang.GetName(), c.outputs)
}
