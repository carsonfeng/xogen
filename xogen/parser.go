package xogen

import (
	"fmt"
	"github.com/carsonfeng/xogen/common/utils"
	"go/ast"
	"go/token"
	"strings"
)

type (
	PropertyInfo struct {
		Name     string
		Return   string
		Source   string
		GetValue string
	}
	ModelInfo struct {
		Package     string
		Name        string
		ModelSpecs  []XormModelSpec
		ModelImport string
	}

	SpecField struct {
		FieldName string
		VarName   string // FieldName的小写字母驼峰命名
		DBField   string // 数据库字段名
		FieldType string
		EndMark   string
	}

	XormModelSpec struct {
		ModelName                       string
		TableName                       string
		TypeSpec                        *ast.TypeSpec
		StructType                      *ast.StructType
		Fields                          []*ast.Field
		AllFields                       []SpecField // 所有spec字段
		AllFieldsExceptIdCreatedUpdated []SpecField // 所有spec字段(除updated)
		UniFields                       []SpecField // 唯一spec字段
		KeyFields                       []SpecField // 重要的spec字段，默认id的下一个，例如uid
		CreatedFields                   []SpecField
		UpdatedFields                   []SpecField
	}
)

func getDaoNameFromOutputPath(outputPath string) (r string) {
	r = outputPath
	s := strings.Split(outputPath, "/")
	if len(s) == 0 {
		return
	}
	daoFileName := s[len(s)-1]

	name := strings.Replace(daoFileName, "_daos.go", "", 1)
	name = strings.Replace(name, "_dao.go", "", 1)

	r = ""
	for i, c := range name {
		if c == '_' {
			continue
		}
		if c == '.' {
			break
		}
		if 0 == i || (i > 0 && daoFileName[i-1] == '_') {
			if c >= 'a' && c <= 'z' {
				r += string(c - 32)
			} else {
				r += string(c)
			}
		} else {
			r += string(c)
		}
	}
	return r
}

func Parse(filename string, modelNames []string, outputPath string, modelImport string) (daoModelInfo *ModelInfo, err error) {
	cf, e := NewCodeFile(filename)
	if e != nil {
		err = e
		return
	}
	daoModelInfo = &ModelInfo{
		Package:     "daos",
		Name:        getDaoNameFromOutputPath(outputPath),
		ModelSpecs:  []XormModelSpec{},
		ModelImport: modelImport,
	}

	for _, d := range cf.file.Decls {
		specs := filterXormTypeDecl(d, modelNames)
		if 0 == len(specs) {
			continue
		}
		for _, spec := range specs {
			daoModelInfo.ModelSpecs = append(daoModelInfo.ModelSpecs, spec)
		}
	}

	return
}

// UniqueFields returns the unique fields of the model.
func (s *XormModelSpec) UniqueFields() (r []string) {
	return s.XFields("unique")
}

// PkFields returns the primary key fields of the model.
func (s *XormModelSpec) PkFields() (r string) {
	t := s.XFields(" pk")
	if len(t) > 0 {
		r = t[0]
	}
	return
}

// XFields returns the fields of the model, which's tag contains X substring.
func (s *XormModelSpec) XFields(x string) (r []string) {
	r = make([]string, 0, len(s.Fields))
	for _, f := range s.Fields {
		if len(f.Names) < 1 {
			continue
		}
		if strings.Contains(strings.ToLower(f.Tag.Value), x) {
			r = append(r, f.Names[0].Name)
		}
	}
	return
}

func (s *XormModelSpec) XormTagFields() (r []string) {
	for _, f := range s.Fields {
		if len(f.Names) < 1 {
			continue
		}
		lowTag := strings.ToLower(f.Tag.Value)
		if strings.Contains(lowTag, "xorm") && "xorm:\"-\"" != lowTag {
			r = append(r, f.Names[0].Name)
		}
	}
	return
}

func (s *XormModelSpec) FilterXormTagFields() {
	var xormFields []*ast.Field = make([]*ast.Field, 0, len(s.Fields))

	var allFields = make([]SpecField, 0, len(s.Fields))
	var allFieldsExceptIdCreatedUpdated = make([]SpecField, 0, len(s.Fields))
	var uniFields = make([]SpecField, 0, len(s.Fields))
	var keyFields = make([]SpecField, 0, len(s.Fields))
	var createdFields = make([]SpecField, 0, len(s.Fields))
	var updatedFields = make([]SpecField, 0, len(s.Fields))

	for i, f := range s.Fields {
		lowTag := strings.ToLower(f.Tag.Value)
		if strings.Contains(lowTag, "xorm") && "xorm:\"-\"" != lowTag {
			xormFields = append(xormFields, f)
			fieldType := ""
			if ident, ok := f.Type.(*ast.Ident); ok {
				fieldType = ident.Name
			} else if sel, ok2 := f.Type.(*ast.SelectorExpr); ok2 && sel.Sel.Name == "Time" {
				fieldType = "time.Time"
			} else {
				continue
			}
			endMark := ","

			allFields = append(allFields, SpecField{
				FieldName: f.Names[0].Name,
				VarName:   utils.FieldToVarName(f.Names[0].Name),
				DBField:   utils.CamelToUnderline(f.Names[0].Name),
				FieldType: fieldType,
				EndMark:   endMark,
			})

			if "Id" != f.Names[0].Name && !strings.Contains(lowTag, "updated") && !strings.Contains(lowTag, "created") {
				allFieldsExceptIdCreatedUpdated = append(allFieldsExceptIdCreatedUpdated, SpecField{
					FieldName: f.Names[0].Name,
					VarName:   utils.FieldToVarName(f.Names[0].Name),
					DBField:   utils.CamelToUnderline(f.Names[0].Name),
					FieldType: fieldType,
					EndMark:   endMark,
				})
			}

			if strings.Contains(lowTag, "unique") {
				uniFields = append(uniFields, SpecField{
					FieldName: f.Names[0].Name,
					VarName:   utils.FieldToVarName(f.Names[0].Name),
					DBField:   utils.CamelToUnderline(f.Names[0].Name),
					FieldType: fieldType,
					EndMark:   endMark,
				})
			}

			if strings.Contains(lowTag, "created") {
				createdFields = append(createdFields, SpecField{
					FieldName: f.Names[0].Name,
					VarName:   utils.FieldToVarName(f.Names[0].Name),
					DBField:   utils.CamelToUnderline(f.Names[0].Name),
					FieldType: fieldType,
					EndMark:   endMark,
				})
			}

			if strings.Contains(lowTag, "updated") {
				updatedFields = append(updatedFields, SpecField{
					FieldName: f.Names[0].Name,
					VarName:   utils.FieldToVarName(f.Names[0].Name),
					DBField:   utils.CamelToUnderline(f.Names[0].Name),
					FieldType: fieldType,
					EndMark:   endMark,
				})
			}

			if 1 == i && f.Names[0].Name != "Id" {
				keyFields = append(keyFields, SpecField{
					FieldName: f.Names[0].Name,
					VarName:   utils.FieldToVarName(f.Names[0].Name),
					DBField:   utils.CamelToUnderline(f.Names[0].Name),
					FieldType: fieldType,
					EndMark:   endMark,
				})
			}
		}
	}

	if len(allFields) > 0 {
		allFields[len(allFields)-1].EndMark = ""
	}
	if len(uniFields) > 0 {
		uniFields[len(uniFields)-1].EndMark = ""
	}
	if len(keyFields) > 0 {
		keyFields[len(keyFields)-1].EndMark = ""
	}

	s.Fields = xormFields
	s.AllFields = allFields
	s.AllFieldsExceptIdCreatedUpdated = allFieldsExceptIdCreatedUpdated
	if len(uniFields) > 0 {
		s.UniFields = uniFields
	} else {
		s.UniFields = keyFields
	}
	s.KeyFields = keyFields
	s.CreatedFields = createdFields
	s.UpdatedFields = updatedFields
}

func filterXormTypeDecl(d ast.Decl, modelNames []string) (r []XormModelSpec) {
	decl, isGen := d.(*ast.GenDecl)
	if !isGen || decl.Tok != token.TYPE || len(decl.Specs) < 1 {
		return nil
	}

	r = make([]XormModelSpec, 0, len(decl.Specs))
	for _, spec := range decl.Specs {
		ts, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			continue
		}
		if !utils.StrInList(ts.Name.Name, modelNames) {
			continue
		}
		xormSpec := XormModelSpec{
			ModelName:  ts.Name.Name,
			TypeSpec:   ts,
			StructType: st,
			Fields:     st.Fields.List,
		}
		xormSpec.TableName = utils.CamelToUnderline(xormSpec.ModelName)
		if len(xormSpec.XormTagFields()) < 1 {
			// 完全没有XORM tag的field，忽略
			continue
		}
		xormSpec.FilterXormTagFields()
		r = append(r, xormSpec)
	}
	return
}

//func filterInterTypeDecl(d ast.Decl) (ts *ast.TypeSpec, interf *ast.InterfaceType) {
//	decl, isGen := d.(*ast.GenDecl)
//	if !isGen || decl.Tok != token.TYPE || len(decl.Specs) < 1 {
//		return nil, nil
//	}
//	ts, ok := decl.Specs[0].(*ast.TypeSpec)
//	if !ok {
//		return nil, nil
//	}
//	inter, ok := ts.Type.(*ast.InterfaceType)
//	if !ok {
//		return nil, nil
//	}
//	if inter.Methods == nil {
//		return nil, nil
//	}
//	return ts, inter
//}

func getMethodInfo(cf *CodeFile, typeName string, method *ast.Field) (*PropertyInfo, error) {
	if len(method.Names) != 1 {
		return nil, fmt.Errorf("%s.%s: 方法名数量不是1", typeName, cf.GetText(method))
	}
	methodName := method.Names[0].Name
	ft, ok := method.Type.(*ast.FuncType)
	if !ok {
		return nil, fmt.Errorf("%s.%s: 不是一个方法", typeName, methodName)
	}
	if ft.Params != nil && len(ft.Params.List) > 0 {
		return nil, fmt.Errorf("%s.%s: 参数表必须为空", typeName, methodName)
	}
	if ft.Results == nil || len(ft.Results.List) != 1 || len(ft.Results.List[0].Names) > 1 {
		return nil, fmt.Errorf("%s.%s: 返回值数量必须为1", typeName, methodName)
	}
	pi := &PropertyInfo{
		Name:   methodName,
		Return: cf.GetText(ft.Results.List[0].Type),
	}
	if c := method.Comment; c == nil || len(c.List) != 1 {
		return nil, fmt.Errorf("%s.%s: 解析注释失败", typeName, methodName)
	} else {
		comment := getCommentLineContent(c.List[0].Text)
		if sp := strings.Index(comment, "."); sp != -1 {
			pi.Source = comment[:sp]
			pi.GetValue = comment[sp+1:]
		} else {
			pi.Source = comment
			pi.GetValue = methodName
		}
	}

	return pi, nil
}

func getCommentLineContent(s string) string {
	if strings.HasPrefix(s, "//") {
		s = s[2:]
	} else if strings.HasPrefix(s, "/*") && strings.HasSuffix(s, "*/") {
		s = s[2 : len(s)-2]
	}
	return strings.TrimSpace(s)
}
