/*
Package Terragrunt defines the structures of generated Terragrunt configuration files. They are
defined here rather than importing the types from the Terragrunt project because of the huge
amount of transitive dependencies that would introduce to the project.
*/
package terragrunt

import "github.com/zclconf/go-cty/cty"

// TerragruntConfig is a complete Terragrunt configuration file, and seeks to mirror
// https://pkg.go.dev/github.com/gruntwork-io/terragrunt@v0.40.2/config#TerragruntConfig.
type TerragruntConfig struct {
	Terraform      TerraformConfig `hcl:"terraform,block"`
	IncludeConfigs []IncludeConfig `hcl:"include,block"`
}

// TerraformConfig represents the `terraform` block in a Terragrunt configuration.
// Mirrors: https://pkg.go.dev/github.com/gruntwork-io/terragrunt@v0.40.2/config#TerraformConfig
type TerraformConfig struct {
	Source string `hcl:"source,attr"`
}

type IncludeConfig struct {
	Name          string  `hcl:"name,label"`
	Path          string  `hcl:"path,attr"`
	Expose        *bool   `hcl:"expose,attr"`
	MergeStrategy *string `hcl:"merge_strategy,attr"`
}

// Dependency represents the `dependency` block in a Terragrunt configuration.
// Mirrors: https://pkg.go.dev/github.com/gruntwork-io/terragrunt@v0.40.2/config#Dependency
type Dependency struct {
	Name                                string     `hcl:",label"`
	ConfigPath                          string     `hcl:"config_path,attr"`
	SkipOutputs                         *bool      `hcl:"skip_outputs,attr"`
	MockOutputs                         *cty.Value `hcl:"mock_outputs,attr"`
	MockOutputsAllowedTerraformCommands *[]string  `hcl:"mock_outputs_allowed_terraform_commands,attr"`
	MockOutputsMergeStrategyWithState		*MergeStrategyType `hcl:"mock_outputs_merge_strategy_with_state"`
}

type MergeStrategyType string

const (
	NoMerge          MergeStrategyType = "no_merge"
	ShallowMerge     MergeStrategyType = "shallow"
	DeepMerge        MergeStrategyType = "deep"
	DeepMergeMapOnly MergeStrategyType = "deep_map_only"
)
