package fileparser

import (
	"github.com/aronlt/toolkit/ds"
)

type FieldInfo struct {
	Name string
	Type string
	Tag  string
}

type StructNode struct {
	fileNode    *FileNode
	name        string
	fields      map[string]*FieldInfo
	embedFields ds.BuiltinSet[string]
	functions   []*FunctionNode
}

func NewStructNode(fileNode *FileNode, name string, embedFields ds.BuiltinSet[string], fields map[string]*FieldInfo) *StructNode {
	return &StructNode{
		fileNode:    fileNode,
		name:        name,
		fields:      fields,
		embedFields: embedFields,
		functions:   make([]*FunctionNode, 0),
	}
}

func (s *StructNode) getIdentity() string {
	return s.fileNode.packageName + "/" + s.name
}
