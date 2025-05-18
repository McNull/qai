package throbber

import "encoding/json"

type Throb struct {
	Interval int      `json:"interval"`
	Frames   []string `json:"frames"`
}

var Throbs = map[string]Throb{}

func init() {
	json.Unmarshal([]byte(throbJSON), &Throbs)
}

func ThrobByName(name string) Throb {
	throb, ok := Throbs[name]
	if !ok {
		return Throb{}
	}
	return throb
}
