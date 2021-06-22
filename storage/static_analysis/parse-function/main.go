package main
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"github.com/z7zmey/php-parser/php7"
//	"github.com/twmb/algoimpl/go/graph"
//	"github.com/yourbasic/graph"
//	"github.com/z7zmey/php-parser/visitor"
	visit "parse-function/visitor"
//	visit "php-api-deps/scan-project/visitor"
//	l "php-api-deps/logger"
	"strings"
)

var fileList []string
var classSet = make(map[string]bool)

func outputPath(path string, basePath string) string {
	relativePath := strings.TrimPrefix(path, basePath)
	return relativePath
}


func processFile(path string, basepath string, funcOutputFile *os.File) {
	fileContents, _ := ioutil.ReadFile(path)
//	l.Log(l.Info, "Processing file: %s", path)
	parser := php7.NewParser(bytes.NewBufferString(string(fileContents)), path)
	parser.Parse()
	rootNode := parser.GetRootNode()
	nsResolver := visit.NewNamespaceResolver()
	rootNode.Walk(nsResolver)
	visit.Path = outputPath(path, basepath)
	visitor := visit.Dumper{
		Writer:     os.Stdout,
		Indent:     "",
		NsResolver: nsResolver,
	}
	rootNode.Walk(visitor)
}


func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

var poi []string
func gettargetedclass () {
	poi = append(poi,"PDO","PDOStatement")
	for {
	size := len(visit.ClassSet)
	for _, item := range visit.Graph {
		foundsrc :=  visit.ClassSet[item.Source]
		founddest := visit.ClassSet[item.Dest]
		if foundsrc && !founddest{
//			fmt.Printf("not found dest:%s\n", item.Dest)
			visit.ClassSet[item.Dest] = true
		}
		if founddest && !foundsrc {
//			fmt.Printf("not found source:%s\n", item.Source)
			visit.ClassSet[item.Source] = true
		}
		if foundsrc && founddest {
//			fmt.Printf("found both:%s:%s\n", item.Source, item.Dest)
		}
		if !foundsrc && !founddest {
//			fmt.Printf("not found any:%s:%s\n", item.Source, item.Dest)
		}
	}	
	newsize := len(visit.ClassSet)
	if size == newsize {
		break
	}
	}

}

func main() {
//	l.Level = l.Debug
	project_path := os.Args[1]
	fileOutput, er := os.Create(os.Args[2])
	err := filepath.Walk(project_path, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ".php" ||
			filepath.Ext(path) == ".install" ||
			filepath.Ext(path) == ".engine" ||
			filepath.Ext(path) == ".module" ||
			filepath.Ext(path) == ".theme" ||
			filepath.Ext(path) == ".inc" {
			fileList = append(fileList, path)
		}
		return nil
	})
	if er != nil {
		log.Fatal(er)
	}
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range fileList {
		processFile(file, project_path, fileOutput)
	}

	gettargetedclass()
	for item,_ := range visit.ClassSet {
		fmt.Fprintf(fileOutput,"%s\n",item)
	}
	for item,_ := range visit.FuncSet {
		fmt.Fprintf(fileOutput,"%s\n",item)
	}
}
