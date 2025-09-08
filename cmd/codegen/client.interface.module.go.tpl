package interfaces

import (
    "context"

{{if eq .ModuleName "emissions" "mint"}}    {{ .PackageName }} "{{ .ImportPath }}"{{else}}    {{ .PackageName }} "{{ .ImportPath }}"{{end}}
)

type {{ .ModuleName | title }}Client interface {
{{ range .Methods -}}
{{ "\t" }}{{ .Name }}(ctx context.Context, req *{{ $.PackageName }}.{{ .RequestType }}) (*{{ $.PackageName }}.{{ .ResponseType }}, error)
{{ end -}}
}

