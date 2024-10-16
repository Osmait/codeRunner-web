package coderunner

import (
	"github.com/Osmait/CodeRunner-web/internal/modules/dispacher"
	programinglanguages "github.com/Osmait/CodeRunner-web/internal/modules/programingLanguages"
	"github.com/Osmait/CodeRunner-web/internal/modules/runner"
)

type CodeRunner struct {
	Runner                    runner.RunnerInterface
	outputs                   *dispacher.Notifier
	listOfProgramingLanguages *programinglanguages.AvailablePrograminLanguages
}

func NewCodeRunner(runner runner.RunnerInterface, outputs *dispacher.Notifier, availablePrograminLanguages *programinglanguages.AvailablePrograminLanguages) *CodeRunner {
	return &CodeRunner{
		Runner:                    runner,
		outputs:                   outputs,
		listOfProgramingLanguages: availablePrograminLanguages,
	}
}

func (c *CodeRunner) RunCode(codeRequest CodeRequest) {
	result, _ := c.listOfProgramingLanguages.SearchLanguage(codeRequest.Lang)
	c.Runner.Execute(codeRequest.Code, result, c.outputs)
}
