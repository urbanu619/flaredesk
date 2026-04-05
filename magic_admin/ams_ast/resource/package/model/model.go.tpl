package {{.Module}}

// 引入关联包
{{- if or .HasDecimal .HasTime }}
import (
	{{- if .HasDecimal }}
	"github.com/shopspring/decimal"
	{{- end }}
	{{- if .HasTime }}
	"time"
	{{- end }}
	{{- if .HasJson }}
	"gorm.io/datatypes"
	{{- end }}
)
{{- end }}

type {{.StructName}} struct { {{range .Fields}}
	{{.Name}} {{.Type}} `{{.JsonTag}} {{.GormTag}}`{{end}}
}

func (*{{.StructName}}) TableName() string {
	return "{{.TableName}}"
}

func New{{.StructName}}() *{{.StructName}} {
	return &{{.StructName}}{}
}
