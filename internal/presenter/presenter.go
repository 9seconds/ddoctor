package presenter

import (
	"encoding/json"

	"github.com/9seconds/ddoctor/internal/checkers"
)

type results struct {
	Results []result `json:"results"`
}

type result struct {
	Ok       bool   `json:"ok"`
	Message  string `json:"message"`
	Producer string `json:"producer"`
}

func Serialize(data []*checkers.CheckResult, prettyPrint bool) ([]byte, error) {
	root := results{Results: make([]result, len(data))}

	for i, v := range data {
		var message string
		if v.Error != nil {
			message = v.Error.Error()
		}

		root.Results[i] = result{
			Ok:       v.Ok,
			Message:  message,
			Producer: v.Producer,
		}
	}

	if prettyPrint {
		return json.MarshalIndent(root, "", "  ")
	}

	return json.Marshal(root)
}
