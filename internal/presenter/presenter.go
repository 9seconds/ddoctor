package presenter

import (
	"encoding/json"
	"fmt"

	"github.com/9seconds/ddoctor/internal/checkers"
)

type results struct {
	results []result
}

type result struct {
	ok       bool
	message  string
	producer string
}

func Serialize(data []*checkers.CheckResult) ([]byte, error) {
	root := results{results: make([]result, len(data))}

	for _, v := range data {
		root.results = append(root.results, result{
			ok:       v.Ok,
			message:  v.Error.Error(),
			producer: v.Producer,
		})
	}

	return json.Marshal(&root)
}
