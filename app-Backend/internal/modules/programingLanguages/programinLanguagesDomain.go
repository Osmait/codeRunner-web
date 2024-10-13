package programinglanguages

type ProgramingLanguages struct {
	Name      string
	Extension string
	Runner    string
}

func NewPrograminLanguages(name, ext, runner string) *ProgramingLanguages {
	return &ProgramingLanguages{
		Name:      name,
		Extension: ext,
		Runner:    runner,
	}
}

func (p *ProgramingLanguages) GetExtension() string {
	return p.Extension
}

func (p *ProgramingLanguages) GetName() string {
	return p.Name
}

func (p *ProgramingLanguages) GetRunner() string {
	return p.Runner
}
