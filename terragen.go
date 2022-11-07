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

// DeploymentsFile represents a file containing a set of Deployments.
type DeploymentsFile struct {
	TerragruntDeployments []TerragruntDeployment     `hcl:"terragrunt_deployment,block"`
	TerragruntIncludeAlls []terragrunt.IncludeConfig `hcl:"terragrunt_include_all,block"`
}

func DecodeFromFile(filename string, out interface{}) error {
	generator, err := NewGenerator()
	if err != nil {
		return err
	}
	return generator.DecodeFromFile(filename, out)
}

func GenerateFromFile(deploymentsFilePath string) error {
	generator, err := NewGenerator()
	if err != nil {
		return err
	}
	return generator.GenerateFromFile(deploymentsFilePath)
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

func (g *Generator) GenerateFromFile(deploymentsFilePath string) error {
	deploymentsFile := &DeploymentsFile{}
	err := DecodeFromFile(deploymentsFilePath, deploymentsFile)
	if err != nil {
		return err
	}
	for _, tgDeployment := range deploymentsFile.TerragruntDeployments {
		generatedConfig := hclwrite.NewEmptyFile()
		allIncludes := append(deploymentsFile.TerragruntIncludeAlls, tgDeployment.Includes...)
		tgConfig := terragrunt.TerragruntConfig{
			Terraform: terragrunt.TerraformConfig{
				Source: tgDeployment.Source,
			},
			IncludeConfigs: allIncludes,
		}
		gohcl.EncodeIntoBody(&tgConfig, generatedConfig.Body())

		destFile, err := g.fs.Create(tgDeployment.DestinationPath)
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
