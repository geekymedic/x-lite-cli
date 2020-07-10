package util

import (
	"github.com/geekymedic/x-lite/xerrors"
	"github.com/kyokomi/emoji"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
)

type MarkdownProperty struct {
	Login         string
	Page          []string
	Zh            string
	URI           string
	RequestTable  []MarkdownTable
	RequestJson   interface{} // 请求参数示例
	ResponseTable []MarkdownTable
	ResponseJson  interface{} // 应答参数示例
}

type MarkdownTable struct {
	Title   string
	Columns []MarkdownTableColumn
}

type MarkdownTableColumn struct {
	FieldName   string      // 参数名称
	FieldType   string      // 类型
	FieldDesc   string      // 参数含义
	FieldIgnore string      // 必填
	DefValue    interface{} // 默认值
	FieldRemark interface{} // 备注
}

func StdoutExit(exitCode int, format string, args ...interface{}) {
	if exitCode == 0 {
		emoji.Fprintf(os.Stdout, ":heavy_check_mark:"+format, args...)
	} else {
		emoji.Fprintf(os.Stdout, ":heavy_multiplication_x:"+format, args...)
	}
	os.Exit(exitCode)
}

func SystemBaseDir() (sysDir string, systemName string, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	sysDir = regexp.MustCompile(".*-system").FindString(dir)
	if sysDir == "" {
		err = xerrors.Format("invalid system directory:%v", dir)
	}

	_, systemName = filepath.Split(sysDir)
	return
}

func defValue(typ string) interface{} {
	switch typ {
	case reflect.Int.String(), reflect.Int8.String(), reflect.Int16.String(), reflect.Int32.String(), reflect.Int64.String(),
		reflect.Uint.String(), reflect.Uint8.String(), reflect.Uint16.String(), reflect.Uint32.String(), reflect.Uint64.String():
		return rand.Intn(1<<8 - 1)
	case reflect.Float32.String(), reflect.Float64.String():
		return rand.Float32()
	case reflect.Bool.String():
		return (rand.Intn(1<<8-1) % 2) == 0
	case reflect.String.String():
		return uuid.Must(uuid.NewV4(), nil).String()
	default:
		return "Not Supper Type"
	}
}
