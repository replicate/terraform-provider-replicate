package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccModelVersionDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccModelVersionDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.replicate_model_version.sdxl", "model", "stability-ai/sdxl"),
					resource.TestCheckResourceAttrSet("data.replicate_model_version.sdxl", "id"),
					resource.TestCheckResourceAttrSet("data.replicate_model_version.sdxl", "versions.#"),
					resource.TestCheckResourceAttrSet("data.replicate_model_version.sdxl", "versions.0.id"),
					resource.TestCheckResourceAttrSet("data.replicate_model_version.sdxl", "versions.0.created_at"),
					resource.TestCheckResourceAttrSet("data.replicate_model_version.sdxl", "versions.0.cog_version"),
				),
			},
		},
	})
}

func testAccModelVersionDataSourceConfig() string {
	return testAccProviderConfig() + `
data "replicate_model_version" "sdxl" {
  model = "stability-ai/sdxl"
}
`
}
