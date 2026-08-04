package models

type AppManifest struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   any    `yaml:"metadata"`
	Spec       any    `yaml:"spec"`
}

type App struct {
	AppManifest `yaml:",inline"`
}
