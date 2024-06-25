data "replicate_hardware" "available" {}

output "hardware_options" {
  value = data.replicate_hardware.available.hardware
}
