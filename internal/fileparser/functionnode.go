package fileparser

import (
	"strings"

	"github.com/aronlt/toolkit/ds"
)

type FunctionNode struct {
	fileNode   *FileNode
	name       string
	receiver   string
	parameters []string
	returns    []string
	content    string
}

func NewFunctionNode(fileNode *FileNode, name string, receiver string, content string, params []string, returns []string) *FunctionNode {
	return &FunctionNode{
		fileNode:   fileNode,
		name:       name,
		receiver:   receiver,
		content:    content,
		parameters: params,
		returns:    returns,
	}
}

func (s *FunctionNode) getIdentity() string {
	return s.fileNode.packageName + "/" + s.name
}

func (s *FunctionNode) String() string {
	data := strings.Builder{}
	data.WriteString("func ")
	data.WriteString(s.name)
	data.WriteString("(")
	for i, parameter := range s.parameters {
		data.WriteString(parameter)
		if i != len(s.parameters)-1 {
			data.WriteString(",")
		}
	}
	data.WriteString(") ")

	if len(s.returns) > 1 {
		data.WriteString("(")
	}
	for i, ret := range s.returns {
		data.WriteString(ret)
		if i != len(s.returns)-1 {
			data.WriteString(",")
		}
	}
	if len(s.returns) > 1 {
		data.WriteString(")")
	}
	return data.String()
}

func (s *FunctionNode) SameParameter(node *FunctionNode) bool {
	return ds.SliceCmpAbsEqual(s.parameters, node.parameters)
}

func (s *FunctionNode) SameReturns(node *FunctionNode) bool {
	return ds.SliceCmpAbsEqual(s.returns, node.returns)
}
