package programinglanguages

import "errors"

type AvailablePrograminLanguages struct {
	Availables []*ProgramingLanguages
}

func NewAvailablePrograminLanguages() *AvailablePrograminLanguages {
	return &AvailablePrograminLanguages{
		Availables: []*ProgramingLanguages{
			NewPrograminLanguages("python", "py", "python"),
			NewPrograminLanguages("javascript", "js", "node"),
		},
	}
}

func (a *AvailablePrograminLanguages) GetListOfAvailablesPrograminLanguages() []*ProgramingLanguages {
	return a.Availables
}

func (a *AvailablePrograminLanguages) IsAvaliable(name string) bool {
	for _, v := range a.GetListOfAvailablesPrograminLanguages() {
		if v.GetName() == name {
			return true
		}
	}
	return false
}

func (a *AvailablePrograminLanguages) SearchLanguage(name string) (*ProgramingLanguages, error) {
	for _, v := range a.GetListOfAvailablesPrograminLanguages() {
		if v.GetName() == name {
			return v, nil
		}
	}
	return nil, errors.New("ProgramingLanguages not Soporting")
}
