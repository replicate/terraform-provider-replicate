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
  owner         = "replicate-testing"
  name          = "terraform-example"
  model         = "replicate/hello-world"
  version       = "5c7d5dc6dd8bf75c1acaa8565735e7986bc5b66206b55cca93cb72c9bf15ccaa"
  hardware      = "cpu"
  min_instances = 0
  max_instances = 1
}
