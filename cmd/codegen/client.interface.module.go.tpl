package interfaces

import (
    "context"

{{if eq .ModuleName "emissions" "mint"}}    {{ .PackageName }} "{{ .ImportPath }}"{{else}}    {{ .PackageName }} "{{ .ImportPath }}"{{end}}

    "github.com/allora-network/allora-sdk-go"
    "github.com/allora-network/allora-sdk-go/config"
)

type {{ .ModuleName | title }}Client interface {
{{ range .Methods -}}
{{ "\t" }}{{ .Name }}(ctx context.Context, req *{{ $.PackageName }}.{{ .RequestType }}, opts ...config.CallOpt) (*{{ $.PackageName }}.{{ .ResponseType }}, error)
{{ end -}}
}

