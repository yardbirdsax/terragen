# terragen

Terragen is library and (eventually) tools for generating code in the [Terraform](https://terraform.io) / [Terragrunt](https://terragrunt.gruntwork.io) ecosystem.

Need to have a bunch of repetitive files for Terragrunt configurations? Don't want to maintain all
of them? Terragen lets you define them in one place, using the same familiar HCL syntax, then
generates the actual files for you.

**Example Terragen configuration file:**
```hcl
terragrunt_configuration "test" {
  source = "mymodule"
  destination_path = "path/to/test/terragrunt.hcl"
  include "something" {
    path = "hello"
  }
}

terragrunt_include_all "all" {
  path = "world"
}
```

Using the `terragen.GenerateFromFile` function, you would get a file at the path
`path/to/test/terragrunt.hcl` that looked like this:

```hcl
terraform {
  source = "mymodule"
}

include "something" {
  path = "hello"
}
include "all" {
  path = "world"
}
```

Right now, Terragen supports a very limited set of syntax for generating minimal Terragrunt
configuration files. The goal is to expand on this feature set, to eventually include things like:

- Support more complete Terragrunt configuration syntax, such as locals, dependencies, etc.
- Support generating Terraform files.
- Fetch file contents from remote sources, using something like [go-getter](https://github.com/hashicorp/go-getter).
- Include a CLI that will perform the generation for you.
