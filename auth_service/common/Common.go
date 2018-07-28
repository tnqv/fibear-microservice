package common

import (
		"bytes"
		"html/template"

)
type AppError struct {
	Err  error
	Code string
	Msg  string
}

type NewAppError struct {
	*AppError
	Show string
}

type DbConfig struct {
	DBHost     string
	DBUser     string
	DBName     string
	DBPassword string
	DBPort 		 int
}

var (
	MAIN_DB_CONSTRING string
)
const (
	AUTHORIZATION_FAILED        = "AUTHORIZATION_FAILED"
	ERROR_DATABASE_CONNECTION   = "ERROR_DATABASE_CONNECTION"
	ERROR_DATABASE_QUERY        = "ERROR_DATABASE_QUERY"
	ERROR_FILE_CREATE           = "ERROR_FILE_CREATE"
	ERROR_FILE_COPY             = "ERROR_FILE_COPY"
	ERROR_FILE_REMOVE           = "ERROR_FILE_REMOVE"
	ERROR_EMAIL_SEND            = "ERROR_EMAIL_SEND"
	INVALID_PARAMS              = "INVALID_PARAMS"
	INCORRECT_PHONE_OR_PASSWORD = "INCORRECT_PHONE_OR_PASSWORD"
	INCORRECT_PASSWORD          = "INCORRECT_PASSWORD"
	ERROR_DEFAULT               = "Có lỗi xảy ra, vui lòng thử lại!"
)

const (
	LETTERS          = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	NUMBER_CHAR_CODE = 6
	MAIN_DB_DRIVER = "mysql"
	TEMPLATE_DB_CONSTRING = `{{ .DBUser }}:{{ .DBPassword }}@tcp({{ .DBHost }}:{{ .DBPort }})/{{ .DBName }}`

//	"user:password@/dbname?charset=utf8&parseTime=True&loc=Local"

)

func LoadTemplate(TemplateStr string, data interface{}) string {
	var msg bytes.Buffer
	content := template.New("TEMPLATE")
	content, _ = content.Parse(TemplateStr)
	content.Execute(&msg, data)
	return string(msg.Bytes())
}
