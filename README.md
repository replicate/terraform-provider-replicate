# Replicate Terraform Provider

[Terraform](https://developer.hashicorp.com/terraform) 
is a tool for [infrastructure as code](https://en.wikipedia.org/wiki/Infrastructure_as_code).
This repository contains a Terraform provider for Replicate
that lets you define your desired deployment configuration 
and run `terraform apply` to create or update it automatically.

## Requirements

- Terraform >= 1.0
- Go >= 1.21

## Getting Started

Clone the GitHub repository and install the provider:

```console
$ gh repo clone replicate/terraform-provider-replicate
$ cd terraform-provider-replicate
$ go install .
```

Configure Terraform to use your local development version of the provider
by creating a file named  `.terraformrc` in your user directory:

```hcl
# ~/.terraformrc
provider_installation {
  dev_overrides {
      "github.com/replicate/replicate" = "~/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```


Create a new file named `example.tf`:

```hcl
# example.tf
terraform {
  required_providers {
    replicate = {
      source = "github.com/replicate/replicate"
    }
  }
}

variable "replicate_api_token" {
  description = "API token for the Replicate provider"
  type        = string
  sensitive   = true
}

provider "replicate" {
  # NOTE: This is populated from the `TF_VAR_replicate_api_token` environment variable.
  api_token = var.replicate_api_token
}

resource "replicate_deployment" "terraform-example" {
  owner         = "mattt"
  name          = "terraform-example"
  model         = "replicate/hello-world"
  version       = "5c7d5dc6dd8bf75c1acaa8565735e7986bc5b66206b55cca93cb72c9bf15ccaa"
  hardware      = "cpu"
  min_instances = 0
  max_instances = 1
}
```

Set your [Replicate API token](https://replicate.com/account/api-tokens):

```console
$ export TF_VAR_replicate_api_token=r8_...
```

Preview the execution plan:

```console
$ terraform plan
# replicate_deployment.terraform-example will be created
  + resource "replicate_deployment" "example-deployment" {
      + owner         = "replicate-testing"
      + name          = "greeter"
      + model         = "replicate/hello-world"
      + version       = "5c7d5dc6dd8bf75c1acaa8565735e7986bc5b66206b55cca93cb72c9bf15ccaa"
      + hardware      = "cpu"
      + max_instances = 1
      + min_instances = 0
    }
```

Apply the changes:

```console
$ terraform apply
```
