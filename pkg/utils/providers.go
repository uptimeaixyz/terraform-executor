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
	BackendPath string
	Providers   []ProviderConfig
	Workspace   string
}

const terraformTemplate = `
terraform {
    backend "local" {
        path = "{{ .BackendPath }}"
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

// func main() {
// 	// Example data for generating the template
// 	data := TerraformTemplateData{
// 		BackendPath: "terraform.tfstate",
// 		Providers: []ProviderConfig{
// 			{Name: "digitalocean", Source: "digitalocean/digitalocean", Version: "2.5.0"},
// 			{Name: "aws", Source: "hashicorp/aws", Version: "4.0.0"},
// 		},
// 		Workspace: "prod",
// 	}

// 	// Generate the Terraform configuration
// 	config, err := GenerateTerraformConfig(data)
// 	if err != nil {
// 		fmt.Println("Error generating Terraform config:", err)
// 		return
// 	}

// 	fmt.Println(config)
// }
