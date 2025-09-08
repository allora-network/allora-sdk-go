package rest

import (
    "context"

    "github.com/brynbellomy/go-utils/errors"
    "github.com/rs/zerolog"

{{if eq .ModuleName "emissions" "mint"}}    {{ .PackageName }} "{{ .ImportPath }}"{{else}}    {{ .PackageName }} "{{ .ImportPath }}"{{end}}
)

// {{ .ModuleName | title }}RESTClient provides REST access to the {{ .ModuleName }} module
type {{ .ModuleName | title }}RESTClient struct {
    *RESTClientCore
    logger zerolog.Logger
}

// New{{ .ModuleName | title }}RESTClient creates a new {{ .ModuleName }} REST client
func New{{ .ModuleName | title }}RESTClient(core *RESTClientCore, logger zerolog.Logger) *{{ .ModuleName | title }}RESTClient {
    return &{{ .ModuleName | title }}RESTClient{
        RESTClientCore: core,
        logger: logger.With().Str("module", "{{ .ModuleName }}").Str("protocol", "rest").Logger(),
    }
}

{{range .Methods}}{{if and .Comment}}// {{ .Comment }}{{end}}
func (c *{{ $.ModuleName | title }}RESTClient) {{ .Name }}(ctx context.Context, req *{{ $.PackageName }}.{{ .RequestType }}) (*{{ $.PackageName }}.{{ .ResponseType }}, error) {
    resp := &{{ $.PackageName }}.{{ .ResponseType }}{}
    err := c.RESTClientCore.executeRequest(ctx,
        "{{ .HTTPMethod }}", "{{ .HTTPPath }}",
        {{if .PathParams}}[]string{ {{range $i, $param := .PathParams}}{{if $i}}, {{end}}"{{ $param }}"{{end}} }{{else}}nil{{end}}, {{if .QueryParams}}[]string{ {{range $i, $param := .QueryParams}}{{if $i}}, {{end}}"{{ $param }}"{{end}} }{{else}}nil{{end}},
        req, resp,
    )
    return resp, errors.Wrap(err, "while calling {{ $.ModuleName | title }}RESTClient.{{ .Name }}")
}

{{end}}