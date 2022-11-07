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
