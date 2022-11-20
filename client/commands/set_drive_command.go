package commands

import "encoding/json"

type SetDriveData struct {
	Steer float64
	Speed float64
}

type SetDriveCommand SetDriveData

func (s SetDriveCommand) GetName() string {
	return "set_drive"
}

func (s SetDriveCommand) GetData() (json.RawMessage, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

var _ Command = (*SetDriveCommand)(nil)
