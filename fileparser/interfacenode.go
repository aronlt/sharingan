package fileparser

import (
	"strings"

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

func (i *InterfaceNode) AllMethods() []*FunctionNode {
	allMethods := ds.SliceGetCopy(ds.MapConvertValueToSlice(i.methods))
	i.embedInterface.ForEach(func(k string) {
		if strings.Contains(k, ".") {
			elems := strings.Split(k, ".")
			if len(elems) != 2 {
				return
			}
			pname := elems[0]
			iname := elems[1]
			fileNodes, ok := i.fileNode.nodeManager.packages[pname]
			if !ok {
				return
			}
			for _, filenode := range fileNodes {
				v, ok := filenode.interfaceNodes[iname]
				if ok {
					allMethods = append(allMethods, ds.MapConvertValueToSlice(v.methods)...)
					break
				}
			}
		} else {
			v, ok := i.fileNode.interfaceNodes[k]
			if !ok {
				return
			}
			allMethods = append(allMethods, ds.MapConvertValueToSlice(v.methods)...)
		}

	})
	return allMethods
}

func (i *InterfaceNode) getIdentity() string {
	return i.fileNode.packageName + "/" + i.name
}
