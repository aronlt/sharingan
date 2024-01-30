package fileparser

type FileNode struct {
	nodeManager    *NodeManager
	file           string
	packageName    string
	structNodes    map[string]*StructNode
	interfaceNodes map[string]*InterfaceNode
	functionNodes  map[string]*FunctionNode
	importers      map[string]string
}

func NewFileNode(nodeManager *NodeManager, file string, packageName string) *FileNode {
	return &FileNode{
		nodeManager:    nodeManager,
		packageName:    packageName,
		file:           file,
		structNodes:    make(map[string]*StructNode, 0),
		interfaceNodes: make(map[string]*InterfaceNode, 0),
		functionNodes:  make(map[string]*FunctionNode, 0),
		importers:      make(map[string]string, 0),
	}
}

func (f *FileNode) AllTokens(pkg string) ([]string, []string, []string) {
	allFunctionTokens := make([]string, 0)
	allStructTokens := make([]string, 0)
	allInterfaceTokens := make([]string, 0)
	for _, function := range f.functionNodes {
		allFunctionTokens = append(allFunctionTokens, pkg+"."+function.name)
	}
	for _, structNode := range f.structNodes {
		allStructTokens = append(allStructTokens, pkg+"."+structNode.name)
	}
	for _, interfaceNode := range f.interfaceNodes {
		allInterfaceTokens = append(allInterfaceTokens, pkg+"."+interfaceNode.name)
	}
	return allFunctionTokens, allStructTokens, allInterfaceTokens
}
