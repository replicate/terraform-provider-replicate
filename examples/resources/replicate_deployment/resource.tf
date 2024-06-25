resource "replicate_deployment" "terraform-example" {
  name          = "terraform-example"
  model         = "replicate/hello-world"
  version       = "5c7d5dc6dd8bf75c1acaa8565735e7986bc5b66206b55cca93cb72c9bf15ccaa"
  hardware      = "cpu"
  min_instances = 1
  max_instances = 2
}
