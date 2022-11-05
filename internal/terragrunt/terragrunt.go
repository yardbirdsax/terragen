/*
Package Terragrunt defines the structures of generated Terragrunt configuration files. They are
defined here rather than importing the types from the Terragrunt project because of the huge
amount of transitive dependencies that would introduce to the project.
*/
package terragrunt

// TerragruntConfig is a complete Terragrunt configuration file, and seeks to mirror
// https://pkg.go.dev/github.com/gruntwork-io/terragrunt@v0.40.0/config#TerragruntConfig.
type TerragruntConfig struct {
	Terraform TerraformConfig `hcl:"terraform,block"`
	IncludeConfigs []IncludeConfig `hcl:"include,block"`
}

// TerraformConfig represents the `terraform` block in a Terragrunt configuration.
// Mirrors: https://pkg.go.dev/github.com/gruntwork-io/terragrunt@v0.40.0/config#TerraformConfig
type TerraformConfig struct {
	Source string `hcl:"source,attr"`
}

type IncludeConfig struct {
	Name string `hcl:"name,label"`
	Path string `hcl:"path,attr"`
	Expose *bool `hcl:"expose,attr"`
	MergeStrategy *string `hcl:"merge_strategy,attr"`
}
