package fileparser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/terror"
	"github.com/pkg/errors"
)

type NodeManager struct {
	projectPath string
	packages    map[string][]*FileNode
}

func GetTypeName(content string, vType ast.Expr) string {
	return content[vType.Pos()-1 : vType.End()-1]
}

func parseFuncType(content string, ft *ast.FuncType) ([]string, []string) {
	parameters := make([]string, 0)
	returns := make([]string, 0)
	// 设置对应的参数和返回值
	if ft.Params != nil && ft.Params.List != nil {
		for _, param := range ft.Params.List {
			typeName := GetTypeName(content, param.Type)
			parameters = append(parameters, typeName)
		}
	}

	if ft.Results != nil && ft.Results.List != nil {
		for _, result := range ft.Results.List {
			typeName := GetTypeName(content, result.Type)
			returns = append(returns, typeName)
		}
	}
	return parameters, returns
}

func (n *NodeManager) ParseStructFunctions() {
	for _, fileNodes := range n.packages {
		functionStructMap := make(map[string][]*FunctionNode)
		for _, fileNode := range fileNodes {
			for _, function := range fileNode.functionNodes {
				ds.MapOpAppendValue(functionStructMap, function.receiver, function)
			}
			for _, structNode := range fileNode.structNodes {
				if functions, ok := functionStructMap[structNode.name]; ok {
					structNode.functions = functions
				}
			}
		}
	}
}

// Inspect 解析源文件，解析出对应的结构体，接口，函数
func (n *NodeManager) Inspect(file string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return terror.Wrapf(err, "call ParseFile fail, file:%s", file)
	}

	bys, err := os.ReadFile(file)
	if err != nil {
		return terror.Wrapf(err, "can't read file:%s content", file)
	}
	content := string(bys)

	fileParser := NewFileNode(n, file, f.Name.Name)

	structParser := func(n ast.Node) bool {
		t, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if t.Type == nil {
			return true
		}

		x, ok := t.Type.(*ast.StructType)
		if !ok {
			return true
		}

		embedFields := ds.NewSet[string]()
		fields := make(map[string]*FieldInfo, 0)
		for _, v := range x.Fields.List {
			if v.Names == nil {
				embedFields.Insert(GetTypeName(content, v.Type))
				continue
			}
			var tag string
			if v.Tag != nil {
				tag = v.Tag.Value
			}

			typeName := GetTypeName(content, v.Type)
			fields[v.Names[0].Name] = &FieldInfo{
				Name: v.Names[0].Name,
				Type: typeName,
				Tag:  tag,
			}
		}

		structNode := NewStructNode(fileParser, t.Name.Name, embedFields, fields)

		fileParser.structNodes[structNode.name] = structNode
		return true
	}

	interfaceParser := func(n ast.Node) bool {
		t, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if t.Type == nil {
			return true
		}

		x, ok := t.Type.(*ast.InterfaceType)
		if !ok {
			return true
		}

		methods := make(map[string]*FunctionNode, 0)
		embedInterface := ds.NewSet[string]()
		if x.Methods != nil || x.Methods.List != nil {
			for _, v := range x.Methods.List {
				_, ok := v.Type.(*ast.Ident)
				if ok {
					embedInterface.Insert(v.Type.(*ast.Ident).Name)
					continue
				}

				ft, ok := v.Type.(*ast.FuncType)
				if ok {
					name := v.Names[0].Name
					parameters, returns := parseFuncType(content, ft)
					functionNode := NewFunctionNode(fileParser, name, "", "", parameters, returns)
					methods[name] = functionNode
					continue
				}
			}
		}
		interfaceNode := NewInterfaceNode(fileParser, t.Name.Name, content[t.Pos()-1:t.End()], embedInterface, methods)
		fileParser.interfaceNodes[t.Name.Name] = interfaceNode
		return true
	}

	// 函数节点解析器
	functionParser := func(n ast.Node) bool {
		x, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		receiver := ""

		// 解析方法接收者
		if x.Recv != nil && x.Recv.List != nil {
			switch t := x.Recv.List[0].Type.(type) {
			case *ast.Ident:
				receiver = t.Name
			case *ast.StarExpr:
				switch innerExpr := t.X.(type) {
				case *ast.Ident:
					receiver = innerExpr.Name
				case *ast.SelectorExpr:
					receiver = innerExpr.Sel.Name
				}
			}
		}
		parameters, returns := parseFuncType(content, x.Type)

		// 方法名，接受者，参数，返回值
		functionNode := NewFunctionNode(fileParser, x.Name.Name, receiver, content[x.Pos()-1:x.End()], parameters, returns)
		fileParser.functionNodes[x.Name.Name] = functionNode

		return true
	}

	importParser := func(n ast.Node) bool {
		x, ok := n.(*ast.ImportSpec)
		if !ok {
			return true
		}
		if x.Name != nil {
			fileParser.importers[x.Name.Name] = x.Path.Value
		} else {
			elems := strings.Split("/", x.Path.Value)
			var name string
			if len(elems) == 1 {
				name = x.Path.Value
			} else {
				name = elems[len(elems)-1]
			}
			fileParser.importers[name] = x.Path.Value
		}
		return true
	}

	ast.Inspect(f, interfaceParser)

	ast.Inspect(f, structParser)

	ast.Inspect(f, functionParser)

	ast.Inspect(f, importParser)
	ds.MapOpAppendValue(n.packages, f.Name.Name, fileParser)

	return nil
}

func (n *NodeManager) getStruct(pkg string, name string) (*StructNode, error) {
	structFileParsers, ok := n.packages[pkg]
	if !ok {
		return nil, errors.Errorf("invalid input, can't find struct package name")
	}
	for _, structParser := range structFileParsers {
		if v, ok := structParser.structNodes[name]; ok {
			return v, nil
		}
	}
	return nil, errors.Errorf("can't find struct")
}

func (n *NodeManager) getInterface(pkg string, name string) (*InterfaceNode, error) {
	interfaceFileParsers, ok := n.packages[pkg]
	if !ok {
		return nil, errors.Errorf("invalid input, can't find interface package name")
	}
	for _, interfaceParser := range interfaceFileParsers {
		if v, ok := interfaceParser.interfaceNodes[name]; ok {
			return v, nil
		}
	}
	return nil, errors.Errorf("can't find interface")
}

func (n *NodeManager) AllTokens() ([]string, []string, []string) {
	allFunctionTokens := make([]string, 0)
	allStructTokens := make([]string, 0)
	allInterfaceTokens := make([]string, 0)
	for pkg, fileParsers := range n.packages {
		for _, fileParser := range fileParsers {
			functionTokens, structTokens, interfaceTokens := fileParser.AllTokens(pkg)
			allFunctionTokens = append(allFunctionTokens, functionTokens...)
			allStructTokens = append(allStructTokens, structTokens...)
			allInterfaceTokens = append(allInterfaceTokens, interfaceTokens...)
		}
	}
	return allFunctionTokens, allStructTokens, allInterfaceTokens
}

func (n *NodeManager) AnalysisInterface(interfaceName string, structName string) ([]string, []string, error) {
	interfaceElems := strings.Split(interfaceName, ".")
	structElems := strings.Split(structName, ".")
	if len(interfaceElems) != 2 || len(structElems) != 2 {
		return nil, nil, errors.Errorf("invalid input, should include package name, like a.b")
	}
	interfaceNode, err := n.getInterface(interfaceElems[0], interfaceElems[1])
	if err != nil {
		return nil, nil, terror.Wrapf(err, "call getInterface fail, interfaceName:%s", interfaceName)
	}
	structNode, err := n.getStruct(structElems[0], structElems[1])
	if err != nil {
		return nil, nil, terror.Wrapf(err, "call getStruct fail, structName:%s", structName)
	}
	missing := make([]string, 0)
	wrong := make([]string, 0)
	for _, method := range interfaceNode.methods {
		found := false
		correct := false
		for _, function := range structNode.functions {
			if function.name == method.name {
				found = true
				correct = function.SameParameter(method) && function.SameReturns(method)
				break
			}
		}
		if !found {
			missing = append(missing, method.String())
		} else if !correct {
			wrong = append(wrong, method.String())
		}
	}
	return missing, wrong, nil

}
