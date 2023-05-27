package core

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"io/ioutil"
	"text/template"
)

// dao CRUD templates

const (
	IMPORT_BEGIN = "// BEGIN_IMPORT, DO NOT MODIFY OR DELETE THIS LINE\n"
	IMPORT_END   = "// END_IMPORT, DO NOT MODIFY OR DELETE THIS LINE\n"
)

func getImports(outfile string) (string, error) {
	return "", nil
	//orig, err := ioutil.ReadFile(outfile)
	//if err != nil {
	//	return "", err
	//}
	//code := string(orig)
	//begin := strings.Index(code, IMPORT_BEGIN)
	//if begin == -1 {
	//	return "", nil
	//}
	//begin += len(IMPORT_BEGIN)
	//end := strings.Index(code[begin:], IMPORT_END)
	//if end == -1 {
	//	return "", nil
	//}
	//end += begin
	//return code[begin:end], nil
}

type renderInfo struct {
	ModelInfo
	Imports     string
	ImportBegin string
	ImportEnd   string
}

func RenderModelImpl(model *ModelInfo, useGeneric bool, outfile string, outMode fs.FileMode, daoTpl string) error {
	var out bytes.Buffer
	var tpl string
	tpl = daoTpl
	imports, err := getImports(outfile)
	fmt.Println(outfile, imports, err)
	if err != nil {
		return err
	}
	//if imports == "" {
	//	imports = "import (\n\"sync\"\n)\n"
	//}
	info := renderInfo{
		ModelInfo:   *model,
		Imports:     imports,
		ImportBegin: IMPORT_BEGIN,
		ImportEnd:   IMPORT_END,
	}
	if err := template.Must(template.New("model").Funcs(getFuncs(model)).Parse(tpl)).Execute(&out, info); err != nil {
		return err
	}
	formatted, err := format.Source(out.Bytes())
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outfile, formatted, outMode)
}

func getFuncs(model *ModelInfo) template.FuncMap {
	return template.FuncMap{
		//"capitalize": func(s string) string {
		//	if len(s) == 0 {
		//		return s
		//	}
		//	return strings.ToUpper(s[0:1]) + s[1:]
		//},
		//"getType": func(sourceName string) (string, error) {
		//	if n, ok := model.SourceTypes[sourceName]; ok {
		//		return n, nil
		//	} else {
		//		return "", fmt.Errorf("未知数据源名: %s", sourceName)
		//	}
		//},
	}
}
