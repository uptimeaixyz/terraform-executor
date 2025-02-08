package utils

import (
	"bytes"
	"text/template"
)

type ProviderConfig struct {
	Name    string
	Source  string
	Version string
}

type TerraformTemplateData struct {
	Bucket    string
	UserID    string
	Context   string
	Workspace string
	Providers []ProviderConfig
}

const terraformTemplate = `
terraform {
    backend "s3" {
        bucket  = "{{ .Bucket }}"
		key     = "{{ .UserID }}/{{ .Context }}/{{ .Workspace }}/terraform.tfstate"
		region  = "eu-west-3"
		profile = "tfstate"
    }
    required_providers {
{{- range .Providers }}
        {{ .Name }} = {
            source = "{{ .Source }}"
            version = "{{ .Version }}"
        }
{{- end }}
    }
}
`

func GenerateTerraformConfig(data TerraformTemplateData) (string, error) {
	tmpl, err := template.New("terraform").Parse(terraformTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
