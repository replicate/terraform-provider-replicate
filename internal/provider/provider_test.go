package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"replicate": providerserver.NewProtocol6WithError(New("acctest")()),
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv(EnvAccApiToken) == "" {
		t.Fatalf("%s must be set to run acceptance tests", EnvAccApiToken)
	}
}

func testAccProviderConfig() string {
	return fmt.Sprintf(`
provider "replicate" {
	api_token = "%s"
}

	`, os.Getenv(EnvAccApiToken))
}
