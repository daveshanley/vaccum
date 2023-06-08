package owasp

import (
	"strings"

	"github.com/daveshanley/vacuum/model"
	"github.com/pb33f/libopenapi/utils"
	"gopkg.in/yaml.v3"
)

type DefineError struct {
}

// GetSchema returns a model.RuleFunctionSchema defining the schema of the DefineError rule.
func (cd DefineError) GetSchema() model.RuleFunctionSchema {
	return model.RuleFunctionSchema{Name: "define_error"}
}

// RunRule will execute the DefineError rule, based on supplied context and a supplied []*yaml.Node slice.
func (cd DefineError) RunRule(nodes []*yaml.Node, context model.RuleFunctionContext) []model.RuleFunctionResult {

	if len(nodes) <= 0 {
		return nil
	}

	var responseCode string
	for i, node := range nodes[0].Content {
		if i%2 == 0 {
			responseCode = node.Value
		} else if responseCode == "400" || responseCode == "422" || strings.ToUpper(responseCode) == "4XX" {
			return []model.RuleFunctionResult{}
		}
	}

	return []model.RuleFunctionResult{
		{
			Message:   "", // TODO
			StartNode: nodes[0],
			EndNode:   utils.FindLastChildNodeWithLevel(nodes[0], 0),
			Path:      responseCode, // TODO
			Rule:      context.Rule,
		},
	}
}