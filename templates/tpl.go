package templates

import "html/template"

var (
	//bff
	BffMainTxt, _   = template.New("").Parse(bffMainTxt)
	BffImplTxt, _   = template.New("").Parse(bffImplTxt)
	BffRouterTxt, _ = template.New("").Parse(bffRouterTxt)
	ErrCodeTxt, _   = template.New("").Parse(errCodeTxt)

	//services
	ServiceMainTxt, _   = template.New("").Parse(serviceMainTxt)
	ServiceServerTxt, _ = template.New("").Parse(serviceServerTxt)
	ServiceImlTxt, _    = template.New("").Parse(serviceImlTxt)

	//config
	ConfigYmlTxt, _ = template.New("").Parse(configYmlTxt)

	//gomod
	GomodTxt, _ = template.New("").Parse(gomodTxt)

	//main.go
	SystemMainTxt, _ = template.New("").Parse(systemMainTxt)

	//make_file
	MakeFileTxt, _ = template.New("").Parse(makeFileTxt)
)
