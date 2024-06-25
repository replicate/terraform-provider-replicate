provider "replicate" {
  # NOTE: This is populated from the `TF_VAR_REPLICATE_API_TOKEN` environment variable.
  api_token = var.replicate_api_token
}

# # Data source to get the latest AMI ID
# data "replicate_model" "stability-ai/sdxl" {
#   most_recent = true
#   owners      = ["self"]
#   filter {
#     name   = "name"
#     values = ["my-ami-*"]
#   }
# }

# # Use the data source to create an EC2 instance
# resource "aws_instance" "example" {
#   ami           = data.aws_ami.latest.id
#   instance_type = "t2.micro"
# }
