package terragen

import (
	"github.com/yardbirdsax/terragen/internal/terragrunt"
)

// TerragruntDeployment is a single configuration for the generation of a Terragrunt configuration.
type TerragruntDeployment struct {
	Name string `hcl:"name,label"`
	// Source is the value passed for the `terraform` block's `source` attribute.
	Source string `hcl:"source,attr"`
	// DestinationPath is the path at which the file should be generated.
	DestinationPath string `hcl:"destination_path,attr"`
	// Includes is used to generate an `includes` block in the generated file.
	Includes []terragrunt.IncludeConfig `hcl:"include,block"`
}
