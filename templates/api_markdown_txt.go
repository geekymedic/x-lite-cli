package templates

const (
	interfaceMarkdownTxt = `### {{.Zh}}
#### 请求方法

> POST

#### 是否需要登录

> {{.Login}}

#### 请求路径

> {{.URI}}

#### 请求格式
{{range $i,$c := .RequestTable}}
**{{$c.Title}}**
{{$fieldsSize := len $c.Columns}}
{{ if (gt $fieldsSize 0)}}
| 参数名称 |类型| 参数含义 |必填|默认值|备注|
| ------ | ------ |------ |------ |------ |------ |
{{range $i,$e := $c.Columns}}| {{$e.FieldName}} | {{$e.FieldType}} | {{$e.FieldDesc}} | {{$e.FieldIgnore}} | {{$e.DefValue}} | {{$e.FieldRemark}} |
{{end}}{{end}}
{{end}}
***Example***:
{{.RequestJson}}

#### 返回格式
{{range $i, $c := .ResponseTable}}
**{{$c.Title}}**
{{$fileSize := len $c.Columns}}
{{ if (gt $fileSize 0)}}
| 参数名称 |类型| 参数含义 |备注|
| ------ | ------ |------ |------ |
{{range $i,$e := $c.Columns}}| {{$e.FieldName}} | {{$e.FieldType}} | {{$e.FieldDesc}} | {{$e.FieldRemark}} |
{{end}}{{end}}
{{end}}
***Example***:
{{.ResponseJson}}
`
)
