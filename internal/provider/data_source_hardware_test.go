package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccHardwareDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccHardwareDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.replicate_hardware.test", "id", "replicate_hardware"),
					resource.TestCheckResourceAttrSet("data.replicate_hardware.test", "hardware.#"),
					resource.TestCheckResourceAttrSet("data.replicate_hardware.test", "hardware.0.name"),
					resource.TestCheckResourceAttrSet("data.replicate_hardware.test", "hardware.0.sku"),
				),
			},
		},
	})
}

func testAccHardwareDataSourceConfig() string {
	return testAccProviderConfig() + `
data "replicate_hardware" "test" {}
`
}
