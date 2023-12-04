package fileparser

import (
	"github.com/aronlt/toolkit/ds"
)

type InterfaceNode struct {
	fileNode       *FileNode
	name           string
	content        string
	methods        map[string]*FunctionNode
	embedInterface ds.BuiltinSet[string]
}

func NewInterfaceNode(fileNode *FileNode, name string, content string, embedInterface ds.BuiltinSet[string], methods map[string]*FunctionNode) *InterfaceNode {
	return &InterfaceNode{
		fileNode:       fileNode,
		name:           name,
		content:        content,
		embedInterface: embedInterface,
		methods:        methods,
	}
}

func (i *InterfaceNode) getIdentity() string {
	return i.fileNode.packageName + "/" + i.name
}
