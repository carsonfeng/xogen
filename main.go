package main

import (
	"flag"
	"fmt"
	"github.com/carsonfeng/xogen/core"
	"github.com/carsonfeng/xogen/testcases"
	"os"
	"strings"
)

/**

xopilot is a code generator for xorm. It generates common CRUD code based on xorm models to dao/service level codes.

usage:
go run main.go -model_path models.go -models struct_name -output_dao output_path

For example:
go run main.go -model_path ../../src/go_pk/app/common/models/models.go -models UserIpRecord,ABC -output_dao ../../src/go_pk/app/common/daos/user_ip_record.go

*/

type Flag struct {
	ModelPath   string
	StructNames []string
	OutputPath  string
	ModelImport string
}

var input Flag

func main() {
	flag.StringVar(&input.ModelPath, "model_path", "testcases/models/models.go", "path of model")
	structNames := ""
	flag.StringVar(&structNames, "models", "DemoTeamModel1,DemoUserModel2", "xorm model struct names, separated by comma ',' ")
	flag.StringVar(&input.OutputPath, "output_dao", "testcases/daos/demo_daos.go", "path of dao output")
	flag.StringVar(&input.ModelImport, "model_import", "xogen/testcases/models", "imported package of model")
	flag.Parse()

	input.StructNames = strings.Split(structNames, ",")

	fmt.Printf("Input: %+v\n", input)

	srcStat, err := os.Stat(input.ModelPath)
	if err != nil {
		panic(err)
	}

	model, err := core.Parse(input.ModelPath, input.StructNames, input.OutputPath, input.ModelImport)

	if nil != err {
		panic(err)
	}

	if err = core.RenderModelImpl(model, false, input.OutputPath, srcStat.Mode().Perm(), testcases.DaoTpl); nil != err {
		panic(err)
	}
	fmt.Printf("Xorm generation completed. \ndao file path: %s\n", input.OutputPath)
}
