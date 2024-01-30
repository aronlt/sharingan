package fileparser

type Parser interface {
	// AnalysisInterface 分析为什么结构体没有实现接口,返回结构体没有实现的方法
	AnalysisInterface(interfaceName string, structName string) ([]string, []string, error)
	AllTokens() ([]string, []string, []string)
}

func NewParser(projectPath string) *NodeManager {
	return &NodeManager{
		projectPath: projectPath,
		packages:    make(map[string][]*FileNode, 0),
	}
}
