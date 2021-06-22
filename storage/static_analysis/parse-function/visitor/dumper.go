package visitor

import  (
	"io"
	"reflect"
	"fmt"
	"regexp"
	"github.com/z7zmey/php-parser/node"
	"github.com/z7zmey/php-parser/node/stmt"
	"github.com/z7zmey/php-parser/node/name"
	"github.com/z7zmey/php-parser/node/expr"
	"github.com/z7zmey/php-parser/node/expr/binary"
	"github.com/z7zmey/php-parser/node/scalar"
	"github.com/z7zmey/php-parser/node/expr/assign"
	"github.com/z7zmey/php-parser/parser"
	"github.com/z7zmey/php-parser/walker"
)

type Dumper struct {
	Writer     io.Writer
	Indent     string
	Comments   parser.Comments
	Positions  parser.Positions
	NsResolver *NamespaceResolver
}

var Path string
var className string
var extClassName string
var ClassSet = make(map[string]bool)
var ImplSet = make (map[string]bool)
var FuncSet = make(map[string]bool)
var methodName string
var functionName string
var metfunc int
var extclass int
var existClass int
var assigns  = make(map[string]string)
var rtn int
type ClassDef struct {
	ClsName string
	Extcls string
	Impcls []string
}

type Edge struct {
	Source string
	Dest string
}

var Graph []Edge
var Classes []ClassDef

func processStringExpr(n node.Node) string {
	switch v:= n.(type) {
	case *scalar.String:
		s := v.Value
		if len(s) > 0 && (s[0] == '"' || s[0] == '\'') {
			s = s[1:]
		}
		if len(s) > 0 && (s[len(s)-1] == '"' || s[len(s)-1] == '\'' ) {
			s = s[:len(s)-1]
		}
		result  := s
		return result
	case *binary.Concat:
		result :=processStringExpr(v.Left)
		result2 := processStringExpr(v.Right)
		return result+result2
	case *expr.ConstFetch:
		constIdentifier := v.Constant.(*name.Name).Parts[0].(*name.NamePart).Value
		return constIdentifier
	case *scalar.Lnumber:
		return v.Value
	case *expr.FunctionCall:
		return "function Call"
	case *expr.Variable:
		return ".*"
	}
	return ""
}
func (d Dumper) EnterNode(w walker.Walkable) bool {
	ClassSet["PDO"] = true
	ClassSet["PDOStatement"] = true
	ClassSet["mysqli"] = true

	n := w.(node.Node)
	switch reflect.TypeOf(n).String(){
	case "*stmt.Class":
		class,ok := n.(*stmt.Class)
		tmp, ok := class.ClassName.(*node.Identifier)
		if (!ok){
			break
		}
		className = tmp.Value
		extends := class.Extends
		var classinfo ClassDef
		classinfo.ClsName = className
		implms := class.Implements
		existClass = 1

		implIntName := ""
		if ( implms != nil) {
			interfaces := implms.InterfaceNames
			_ = interfaces
			//fmt.Printf("PATH:%s,%s,%d\n",Path,className, len(interfaces))
			for i := 0 ; i <len(interfaces); i++ {
				implName, ok1 := interfaces[i].(*name.Name)
				if ok1 {
					imlIntName := implName.Parts[0].(*name.NamePart).Value
//					fmt.Fprintf(d.Writer, "cls:%s:%s\n", className, imlIntName)
					classinfo.Impcls = append(classinfo.Impcls, implName.Parts[0].(*name.NamePart).Value)
					E := Edge{className,implName.Parts[0].(*name.NamePart).Value}
					Graph = append(Graph, E)
					for cls := range ClassSet {
						if ( imlIntName == cls ) {
							exist := ClassSet[imlIntName]
							if !exist {	
								ClassSet[imlIntName] = true
//								fmt.Fprintf(d.Writer,"added int %s\n", imlIntName )
							}
						}
					}
				}
			}
			implName, ok1  := implms.InterfaceNames[0].(*name.Name)
			if ok1 {
				implIntName = implName.Parts[0].(*name.NamePart).Value
			}
		}
		if extends != nil {
			extclass = 1
			extendsName, ok := extends.ClassName.(*name.Name)
			if ok {
				extClassName = extendsName.Parts[0].(*name.NamePart).Value
				classinfo.Extcls = extClassName
				E := Edge{className, extClassName}
				Graph = append(Graph, E)
				for  cls := range ClassSet {
					if ( extClassName == cls || implIntName == cls ){
					if  implIntName != "" {
						exist := ClassSet[implIntName]
						if !exist {
							ClassSet[implIntName] = true
						}
					}
					exists := ClassSet[extClassName]
					if !exists {
						ClassSet[extClassName] = true
					}
					exists = ClassSet[className]
					if !exists{
						ClassSet[className] = true

					}
//					fmt.Fprintf(d.Writer,"add classes:[%s][%s][%s] bc of %s\n", className, extClassName, implIntName, cls)
//					fmt.Fprintf(d.Writer, "%s:%s:%s\n", className, extClassName, implIntName)
					}
				}
			}
		} else {
			extclass = 0
		}
		Classes = append(Classes, classinfo)
//		_ = classes

		break
	case "*stmt.Interface":
		interf := n.(*stmt.Interface)
		intName := interf.InterfaceName.(*node.Identifier).Value
		intext := interf.Extends
		if intext != nil {
			intextNameArr := intext.InterfaceNames
			for i := 0; i < len(intextNameArr) ; i++{
				intextName, ok1 := intextNameArr[i].(*name.Name)
				if ok1{
					_ = intName
					_ = intextName
//					fmt.Fprintf(d.Writer, "interface: %s:%s\n", intName, intextName.Parts[0].(*name.NamePart).Value)
				}
			}

		}
	
	case "*stmt.ClassMethod":
		method := n.(*stmt.ClassMethod)
		methodName = method.MethodName.(*node.Identifier).Value
//		fmt.Fprintf(d.Writer, "class method: %s:%s:%s\n",Path , className , methodName)
		metfunc = 1
		if methodName == "__construct" {
			paramList := method.Params
			for i := 0 ; i < len(paramList); i++ {
				param := paramList[i].(*node.Parameter)
				paramType , ok := param.VariableType.(*name.Name)
				if ok {
					paramTypeName := paramType.Parts[0].(*name.NamePart).Value
				if ClassSet[paramTypeName] == true {
					ClassSet[className] = true
//					fmt.Fprintf(d.Writer,"add class bc of constructor [%s]\n", className)
				}

			}
		}
		}
		break
	case "*expr.FunctionCall":
		function := n.(*expr.FunctionCall)
		methodname, ok := function.Function.(*name.Name)
		if ok {
//			if (  methodname.Parts[0].(*name.NamePart).Value == `execute` || 
			if (	methodname.Parts[0].(*name.NamePart).Value == `mysqli_query` || methodname.Parts[0].(*name.NamePart).Value == `mysql_query`) {// || methodname.Parts[0].(*name.NamePart).Value == `_do_query`) {
				if metfunc == 0 {
//					fmt.Fprintf(d.Writer, "PATH: [%s] funcname:[%s]\n", Path, functionName)
				}
				if metfunc == 1 {
					if extclass == 1 {
//						fmt.Fprintf(d.Writer, "path: [%s] classname:[%s] extendedclass:[%s] methodname:[%s] \n", Path, className, extClassName, methodName)
//						fmt.Fprintf(d.Writer,"%s\n%s\n",className, extClassName)
						ClassSet[className] = true
				}
					if extclass == 0 {
						//fmt.Fprintf(d.Writer,"%s\n", className)
						ClassSet[className] = true
//						fmt.Fprintf(d.Writer, "path: [%s] classname:[%s] methodname:[%s] \n", Path, className, methodName)
					}
				}
			}
			if ( methodname.Parts[0].(*name.NamePart).Value == `query` ) {
//					fmt.Fprintf(d.Writer, "the query is: %s \n", function.ArgumentList.Arguments[0].(*node.Argument).Expr)
				}
		}
		break
	case "*stmt.Function":
		function := n.(*stmt.Function)
		functionName = function.FunctionName.(*node.Identifier).Value
//		fmt.Fprintf(d.Writer, "function: %s:%s\n", Path, functionName)
		metfunc = 0
		break
	case "*stmt.Expression":
		asgnStmt := n.(*stmt.Expression)
		s, ok := asgnStmt.Expr.(*assign.Assign)
		if  ok {
			if varname, ok := s.Variable.(*expr.Variable); ok {
				if varname, ok := varname.VarName.(*node.Identifier); ok {
					switch v:= s.Expression.(type){
					case *expr.New:
						classname ,ok := v.Class.(*expr.Variable)
						if ok {
							varNew := classname.VarName.(*node.Identifier).Value
//							fmt.Fprintf(d.Writer, "assign:%s:%s:%s:%s\n",className, varname.Value, assigns[varNew], varNew)
							assigns[varname.Value] = assigns[varNew]
						}
						break
					default:
						value := processStringExpr(s.Expression)
						assigns[varname.Value] = value
//						fmt.Fprintf(d.Writer,"assign: [%s]:[%s]:[%s]\n",varname.Value, value, Path)
					}
			}

		}
		if varname, ok := s.Variable.(*expr.PropertyFetch); ok {
			_ = varname
			switch v := s.Expression.(type) {
			case *expr.New:
				classname , ok := v.Class.(*name.Name)
				if ok {
					varNew := classname.Parts[0].(*name.NamePart).Value

					if ClassSet[varNew] == true {
						if metfunc == 1 {
							ClassSet[className] = true
//							fmt.Fprintf(d.Writer, "add class [%s] bc of fetching from [%s]\n", className, varNew)
						}

					}
				}
			}
		}
		}
		break
	case "*expr.New" :
		if rtn == 0 {
			break
		}
		expNew := n.(*expr.New)
		classname ,ok := expNew.Class.(*expr.Variable)
		if ok {
		for cls := range ClassSet {
			varNew := classname.VarName.(*node.Identifier).Value
			_, exist := assigns[varNew]
			if exist {
				r , _ := regexp.Compile("(\\.|\\*|\\W)*")
				match := r.MatchString(assigns[varNew])
				if match {
					lenfound := len(r.FindString(assigns[varNew]))
					if lenfound == len(assigns[varNew]){
						break
					}
				}
				if (assigns[varNew] == "" || assigns[varNew] == ".*"){
					break
				}

				r, err := regexp.Compile(assigns[varNew])
				if err != nil {
					break
				}
				matched := r.MatchString(cls)
			if ( matched ) {
//				fmt.Fprintf(d.Writer,"new op:[%s]:[%s]:[%s]:[%s]\n",assigns[varNew],varNew, cls, Path)
				exist := false
				if metfunc == 0 {
					exist = FuncSet[functionName]
				}
				if metfunc == 1 {
					exist = ClassSet[className] || ClassSet[methodName]
				}
				if !exist {
					if metfunc == 0 {
//						fmt.Fprintf(d.Writer, "new op:%s:%s:%s\n",functionName, cls, Path)
						FuncSet[functionName] = true
//						fmt.Fprintf(d.Writer,"function:[%s]\n",functionName)
					}
					if metfunc == 1 {
						fmt.Fprintf(d.Writer, "new op:%s:%s:%s\n", className, methodName, Path)
						ClassSet[className] = true
						//ClassSet[methodName] = true
						fmt.Fprintf(d.Writer,"class,method:[%s][%s]\n",className,methodName)
					}
				}
//				ClassSet[varNew] = true
			}
		}
		}
	}

	break
	case "*expr.StaticCall":
		call  := n.(*expr.StaticCall).Call
		//functionname, ok := call.(*node.Identifier)
		//_ = functionname
		if rtn == 0{
			break
		}
		if functionName == "" {
			break
		}
		class, ok := n.(*expr.StaticCall).Class.(*name.Name)
		if ok {
			classname := class.Parts[0].(*name.NamePart).Value
			for cls := range ClassSet {
				if ( cls == classname ){
					FuncSet[functionName] = true
					fmt.Fprintf(d.Writer, "scall:%s:%s:%s\n", classname, functionName, Path)
				}
			}
		if class.Parts[0].(*name.NamePart).Value == `PDO` {
		function := call.(*node.Identifier).Value
		_ = function
		if metfunc == 0 {
//			fmt.Fprintf(d.Writer, "called function:[%s]:[%s]:[%s]\n", Path, functionName, function)
		}
		if metfunc == 1 {
//			fmt.Fprintf(d.Writer, "called function:[%s]:[%s]:[%s]\n", Path, methodName,  function)
		}
	}
	}
	break
	case "*stmt.Return":
		rtn = 1
		ret := n.(*stmt.Return)
		vr, ok := ret.Expr.(*expr.Variable)
		if ok {
			varName := vr.VarName.(*node.Identifier).Value
			if metfunc == 0 {
//				fmt.Fprintf(d.Writer, "%s:%s\n", functionName, assigns[varName])
			}
			if metfunc == 1{
//				fmt.Fprintf(d.Writer, "%s:%s\n", className, assigns[varName])
			}
			_, exist := assigns[varName]
			if exist{
				// the value is in the list of assignments
				// check whether we can resolve it to a database class that extends database API in PHP
				r , _ := regexp.Compile("(\\.|\\*|\\W)*")
				match := r.MatchString(assigns[varName])
				if match {
					lenfound := len(r.FindString(assigns[varName]))
					if lenfound == len(assigns[varName]){
						break
					}
				}
				if (assigns[varName] == "" || assigns[varName] == ".*"){
					break
				}
//				fmt.Fprintf(d.Writer, "return value: %s:%s:[%s]\n", varName, assigns[varName], Path)
				for cls := range ClassSet {
//					fmt.Fprintf(d.Writer, "the class is [%s]\n", cls)
					r , err := regexp.Compile(assigns[varName])
					if err != nil {
						break
					}
//					fmt.Fprintf(d.Writer,"%s\n", err)
					matched := r.MatchString(cls)
					if matched { // found a match for the object type, means that the method or function is returning an object from sub-type of database API
						lenfound := len(r.FindString(cls))
						if len(cls) !=  lenfound {
							break
						}
						if metfunc == 0 {
//							fmt.Fprintf(d.Writer, "added %s:%s\n", functionName, cls)
							FuncSet[functionName] = true
						}
						if metfunc == 1 {
							if extclass == 1 {
//								fmt.Fprintf(d.Writer, "path: [%s] classname:[%s] extendedclass:[%s] methodname:[%s] \n", Path, className, extClassName, methodName)
//								fmt.Fprintf(d.Writer,"added: %s:%s:%s:%s\n",className, varName, assigns[varName],cls)
								ClassSet[className] = true
							}
							if extclass == 0 {
//								fmt.Fprintf(d.Writer,"%s\n were added bc [%s] return database object", className, methodName)
								ClassSet[className] = true
//								fmt.Fprintf(d.Writer, "added: %s:%s:%s:%s \n", className, varName,  assigns[varName],cls)
							}
						}
					}
				}

			}
		}
	}
	return true
}
// GetChildrenVisitor is invoked at every node parameter that contains children nodes
func (d Dumper) GetChildrenVisitor(key string) walker.Visitor {
	return Dumper{d.Writer, d.Indent + "    ", d.Comments, d.Positions, d.NsResolver}
}



// LeaveNode is invoked after node process
func (d Dumper) LeaveNode(w walker.Walkable) {
	// do nothing	
	n := w.(node.Node)
	switch reflect.TypeOf(n).String(){
	case "*stmt.Return":
		rtn = 0
	case "*stmt.Function":
		functionName = ""
	}
}
