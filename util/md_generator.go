package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/geekymedic/x-lite-cli/templates"
	"github.com/geekymedic/x-lite/xerrors"
	"go/ast"
	"go/parser"
	"go/token"
	"html"
	"os"
	"path/filepath"
	"strings"
	"github.com/tidwall/pretty"
)

type XType interface {
}

func String(xt XType) string {
	switch t := xt.(type) {
	case *XMap:
		return "map[" + String(t.Key) + "]" + String(t.Value)
	case *XArray:
		return "[]" + String(t.EleType)
	case *XPoint:
		return String(t.X)
	case *XBase:
		return t.BName
	default:
	}
	return ""
}

func ValueIg(xt XType) interface{} {
	switch t := xt.(type) {
	case *XMap:
		key := fmt.Sprintf("%v", ValueIg(t.Key))
		val := ValueIg(t.Value)
		return map[string]interface{}{key: val}
	case *XArray:
		return []interface{}{ValueIg(t.EleType)}
	case *XPoint:
		return ValueIg(t.X)
	case *XBase:
		if t.IsBasic() {
			return defValue(t.BName)
		} else {
			if _, ok := xTs[t.BName]; ok {
				return xTs[t.BName].ToMap()
			}
		}
	default:
		return nil
	}
	return nil
}

type XMap struct {
	Key   XType
	Value XType
}

type XArray struct {
	EleType XType
}

type XPoint struct {
	X XType
}

//buildin type (exclude map slice array etc.) or user define struct type
type XBase struct {
	BName string
}

func (x *XBase) IsBasic() bool {
	basic := []string{
		"bool", "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64",
		"string",
	}
	for _, b := range basic {
		if b == x.BName {
			return true
		}
	}
	return false
}

type XField struct {
	FName   string
	Type    XType
	Comment XFieldComment
}

type XFieldComment struct {
	Desc     string
	Required string
	DefaultV string
	Remark   string
}

func parseFieldComment(text string) XFieldComment {
	splits := strings.SplitN(text, "|", 4)
	splits = append(splits, []string{"", "", "", ""}...)
	return XFieldComment{
		Desc:     splits[0],
		Required: splits[1],
		DefaultV: splits[2],
		Remark:   splits[3],
	}
}

type XStruct struct {
	Fields []*XField
}

func (xs *XStruct) ToMap() map[string]interface{} {
	js := map[string]interface{}{}
	for _, field := range xs.Fields {
		if field.FName != "" {
			js[field.FName] = ValueIg(field.Type)
		} else {
			xt := field.Type
			vIg := ValueIg(xt)
			if sjs, ok := vIg.(map[string]interface{}); ok {
				for k, v := range sjs {
					js[k] = v
				}
			}
		}
	}
	return js
}

func (xs *XStruct) ToMarkdownTable() []MarkdownTableColumn {
	cls := []MarkdownTableColumn{}
	for _, field := range xs.Fields {
		if field.FName != "" {
			cls = append(cls, MarkdownTableColumn{
				FieldName:   field.FName,
				FieldType:   String(field.Type),
				FieldDesc:   field.Comment.Desc,
				FieldIgnore: field.Comment.Required,
				DefValue:    field.Comment.DefaultV,
				FieldRemark: field.Comment.Remark,
			})
		} else {
			xBaseName := String(field.Type)
			if _, ok := xTs[xBaseName]; ok {
				mt, succ := xTs[xBaseName].ToMarkdownTable()
				if succ {
					cls = append(cls, mt.Columns...)
				}
			}
		}
	}
	return cls
}

type XTypeSpec struct {
	TName string
	Type  XType
}

func (xSpec *XTypeSpec) ToMap() map[string]interface{} {
	switch t := xSpec.Type.(type) {
	case *XStruct:
		return t.ToMap()
	default:
	}
	return nil
}

func (xSpec *XTypeSpec) ToMarkdownTable() (MarkdownTable, bool) {
	switch t := xSpec.Type.(type) {
	case *XStruct:
		return MarkdownTable{Title: xSpec.TName, Columns: t.ToMarkdownTable()}, true
	default:
	}
	return MarkdownTable{}, false
}

type XInterfaceNode struct {
	ZhName   string
	Login    string
	Uri      string
	Page     string
	Desc     string
	Request  *XTypeSpec
	Response *XTypeSpec
}

func (in *XInterfaceNode) ToMarkDown() *MarkdownProperty {
	md := &MarkdownProperty{
		Login: in.Login,
		Zh:    in.ZhName,
		URI:   in.Uri,
		//RequestTable: []MarkdownTable,
		//RequestJson:  interface{},
		//ResponseTable: []MarkdownTable,
		//ResponseJson: interface{},
	}

	if mt, ok := in.Request.ToMarkdownTable(); ok {
		md.RequestTable = []MarkdownTable{mt}
	}

	if mt, ok := in.Response.ToMarkdownTable(); ok {
		md.ResponseTable = []MarkdownTable{mt}
	}

	b, _ := json.Marshal(in.Request.ToMap())

	md.RequestJson = fmt.Sprintf("```json\n%s```\n", pretty.Pretty(b))

	b, _ = json.Marshal(in.Response.ToMap())

	md.ResponseJson = fmt.Sprintf("```json\n%s```\n", pretty.Pretty(b))

	return md
}

func parseTypeExpr(expr ast.Expr) XType {
	switch t := expr.(type) {
	case *ast.Ident:
		return &XBase{BName: t.Name}
	case *ast.StructType:
		if t.Fields == nil {
			break
		}
		return parseStructType(t)
	case *ast.ArrayType:
		return &XArray{EleType: parseTypeExpr(t.Elt)}
	case *ast.StarExpr:
		return &XPoint{X: parseTypeExpr(t.X)}
	case *ast.MapType:
		return &XMap{Key: parseTypeExpr(t.Key), Value: parseTypeExpr(t.Value)}
	default:
	}
	return nil
}

func parseField(field *ast.Field) *XField {
	var name string
	if field.Names != nil && len(field.Names) > 0 {
		name = field.Names[0].Name
	}
	var comment XFieldComment
	if field.Comment != nil && len(field.Comment.List) > 0 {
		comment = parseFieldComment(field.Comment.Text())
	}
	switch t := field.Type.(type) {
	case *ast.Ident:
		return &XField{FName: name, Type: parseTypeExpr(t), Comment: comment}
	case *ast.ArrayType:
		return &XField{FName: name, Type: parseTypeExpr(t), Comment: comment}
	case *ast.StarExpr:
		return &XField{FName: name, Type: parseTypeExpr(t), Comment: comment}
	case *ast.MapType:
		return &XField{FName: name, Type: parseTypeExpr(t), Comment: comment}
	default:
	}
	return nil
}

func parseStructType(t *ast.StructType) *XStruct {
	if t.Fields == nil {
		return nil
	}
	xStruct := &XStruct{}
	for _, field := range t.Fields.List {
		f := parseField(field)
		if f == nil {
			continue
		}
		xStruct.Fields = append(xStruct.Fields, parseField(field))
	}
	return xStruct
}

func parseTypeSpec(spec *ast.TypeSpec) *XTypeSpec {
	name := spec.Name.Name
	t := parseTypeExpr(spec.Type)
	if name == "" || t == nil {
		return nil
	}
	return &XTypeSpec{TName: name, Type: parseTypeExpr(spec.Type)}
}

var xTs map[string]*XTypeSpec

func parseInterfaceComments(c *ast.CommentGroup, in *XInterfaceNode) (ok bool) {
	if c != nil {
		for _, m := range c.List {
			splits := strings.SplitN(m.Text, "@", 2)
			if len(splits) == 2 {
				subSplits := strings.SplitN(splits[1], ":", 2)
				if len(subSplits) == 2 {
					k := strings.TrimSpace(subSplits[0])
					v := strings.TrimSpace(subSplits[1])
					switch k {
					case "type":
						if v == "b.i" {
							ok = true
						}
					case "name":
						in.ZhName = v
					case "login":
						in.Login = v
					case "page":
						in.Page = v
					case "uri":
						in.Uri = v
					case "describe":
						in.Desc = v
					default:
					}
				}
			}
		}
	}
	return
}

func CreateMd(sysDir string, bffName string, impls []string) error {
	docBaseDir := filepath.Join(sysDir, "doc")
	err := os.MkdirAll(docBaseDir, os.ModePerm)
	if err != nil {
		return err
	}

	docBffDir := filepath.Join(docBaseDir, bffName)
	err = os.MkdirAll(docBffDir, os.ModePerm)
	if err != nil {
		return err
	}

	var filter = func(info os.FileInfo) bool {
		if strings.HasSuffix(info.Name(), ".go") {
			return true
		}
		return false
	}

	for _, impl := range impls {
		xTs = make(map[string]*XTypeSpec, 0)
		xIfNodes := make(map[string]*XInterfaceNode, 0)
		implDir := filepath.Join(sysDir, "bff", bffName, "impls", impl)
		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, implDir, filter, parser.ParseComments|parser.DeclarationErrors)
		if err != nil {
			panic(err)
		}

		for _, pkg := range pkgs {
			for fname, f := range pkg.Files {
				var (
					ok bool
					in XInterfaceNode
				)
				{
					ast.Inspect(f, func(n ast.Node) bool {
						switch t := n.(type) {
						case *ast.GenDecl:
							var (
								isRequestDecl  bool
								isResponseDecl bool
							)
							if t.Doc != nil {
								if strings.Contains(t.Doc.Text(), "b.i.rt") {
									isRequestDecl = true
								} else if strings.Contains(t.Doc.Text(), "b.i.re") {
									isResponseDecl = true
								}
							}
							if t.Tok == token.TYPE {
								for _, spec := range t.Specs {
									if specType, ok := spec.(*ast.TypeSpec); !ok {
										continue
									} else {
										ts := parseTypeSpec(specType)
										xTs[ts.TName] = ts
										if isRequestDecl {
											in.Request = ts
										}
										if isResponseDecl {
											in.Response = ts
										}
									}
								}
							}
						case *ast.FuncDecl:
							ok = parseInterfaceComments(t.Doc, &in)
						default:
						}
						return true
					})
				}
				if ok {
					xIfNodes[fname] = &in
				}
			}
		}

		for fname, xIf := range xIfNodes {
			_, goFile := filepath.Split(fname)
			mdFile := strings.Replace(goFile, ".go", ".md", -1)

			mp := xIf.ToMarkDown()

			var b bytes.Buffer
			err = templates.InterfaceMarkdownTxt.Execute(&b, &mp)
			c := html.UnescapeString(b.String())

			file := filepath.Join(docBffDir, mdFile)
			if _, err := os.Stat(file); os.IsNotExist(err) {
				f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
				if err != nil {
					return xerrors.By(err)
				}
				defer f.Close()
				f.WriteString(c)
			}
		}
	}
	return nil
}
