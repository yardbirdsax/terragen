package terragen

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/afero"
	"github.com/yardbirdsax/terragen/internal/terragrunt"
)

// ConfigurationsFile represents a file containing a set of configurations.
type ConfigurationsFile struct {
	// These are used for generating Terragrunt root files.
	TerragruntConfigurations []TerragruntConfiguration `hcl:"terragrunt_configuration,block"`
	// These `include` blocks will be present in all generated files.
	TerragruntIncludeAlls []terragrunt.IncludeConfig `hcl:"terragrunt_include_all,block"`
}

func DecodeFromFile(filename string, out interface{}) error {
	generator, err := NewGenerator()
	if err != nil {
		return err
	}
	return generator.DecodeFromFile(filename, out)
}

func GenerateFromFile(configurationsFilePath string) error {
	generator, err := NewGenerator()
	if err != nil {
		return err
	}
	return generator.GenerateFromFile(configurationsFilePath)
}

func NewGenerator(optFns ...GeneratorOptsFn) (*Generator, error) {
	generator := &Generator{
		fs: afero.NewOsFs(),
	}
	for _, f := range optFns {
		err := f(generator)
		if err != nil {
			return nil, err
		}
	}
	return generator, nil
}

type Generator struct {
	fs afero.Fs
}

func (g *Generator) DecodeFromFile(filename string, out interface{}) error {
	file, err := g.fs.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()
	fileContentBytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	err = hclsimple.Decode(filename, fileContentBytes, nil, out)
	return err
}

func (g *Generator) GenerateFromFile(configurationsFilePath string) error {
	configurationsFile := &ConfigurationsFile{}
	err := DecodeFromFile(configurationsFilePath, configurationsFile)
	if err != nil {
		return err
	}
	err = g.GenerateFromConfig(configurationsFile)
	return err
}

func (g *Generator) GenerateFromConfig(config *ConfigurationsFile) error {
	for _, configuration := range config.TerragruntConfigurations {
		generatedConfig := hclwrite.NewEmptyFile()
		allIncludes := append(config.TerragruntIncludeAlls, configuration.Includes...)
		tgConfig := terragrunt.TerragruntConfig{
			Terraform: terragrunt.TerraformConfig{
				Source: configuration.Source,
			},
			IncludeConfigs: allIncludes,
		}
		gohcl.EncodeIntoBody(&tgConfig, generatedConfig.Body())

		destFile, err := g.fs.Create(configuration.DestinationPath)
		if err != nil {
			return err
		}
		defer destFile.Close()
		_, err = destFile.Write(generatedConfig.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}

type GeneratorOptsFn func(*Generator) error

// WithFs sets the file system for the Generator.
func WithFs(fs afero.Fs) GeneratorOptsFn {
	return func(g *Generator) error {
		g.fs = fs
		return nil
	}
}
