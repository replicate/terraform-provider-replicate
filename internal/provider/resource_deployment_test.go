package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeploymentResource(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDeploymentResourceConfig("replicate-testing", rName, "replicate/hello-world", "5c7d5dc6dd8bf75c1acaa8565735e7986bc5b66206b55cca93cb72c9bf15ccaa", "cpu", 0, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("replicate_deployment.test", "owner", "replicate-testing"),
					resource.TestCheckResourceAttr("replicate_deployment.test", "name", rName),
					resource.TestCheckResourceAttr("replicate_deployment.test", "model", "replicate/hello-world"),
					resource.TestCheckResourceAttr("replicate_deployment.test", "version", "5c7d5dc6dd8bf75c1acaa8565735e7986bc5b66206b55cca93cb72c9bf15ccaa"),
					resource.TestCheckResourceAttr("replicate_deployment.test", "hardware", "cpu"),
					resource.TestCheckResourceAttr("replicate_deployment.test", "min_instances", "0"),
					resource.TestCheckResourceAttr("replicate_deployment.test", "max_instances", "1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "replicate_deployment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDeploymentResourceConfig("replicate-testing", rName, "replicate/hello-world", "5c7d5dc6dd8bf75c1acaa8565735e7986bc5b66206b55cca93cb72c9bf15ccaa", "gpu-t4", 2, 4),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("replicate_deployment.test", "hardware", "gpu-t4"),
					resource.TestCheckResourceAttr("replicate_deployment.test", "min_instances", "2"),
					resource.TestCheckResourceAttr("replicate_deployment.test", "max_instances", "4"),
				),
			},
			{
				Config: testAccDeploymentResourceConfig("replicate-testing", rName, "replicate/hello-world", "5c7d5dc6dd8bf75c1acaa8565735e7986bc5b66206b55cca93cb72c9bf15ccaa", "gpu-t4", 2, 4),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("replicate_deployment.test", "hardware", "cpu"),
					resource.TestCheckResourceAttr("replicate_deployment.test", "min_instances", "0"),
					resource.TestCheckResourceAttr("replicate_deployment.test", "max_instances", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDeploymentResourceConfig(owner, name, model, version, hardware string, minInstances, maxInstances int) string {
	return fmt.Sprintf(testAccProviderConfig()+`
resource "replicate_deployment" "test" {
  owner         = %[1]q
  name          = %[2]q
  model         = %[3]q
  version       = %[4]q
  hardware      = %[5]q
  min_instances = %[6]d
  max_instances = %[7]d
}
`, owner, name, model, version, hardware, minInstances, maxInstances)
}
