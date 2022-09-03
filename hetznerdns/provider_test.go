package hetznerdns

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]func() (*schema.Provider, error)

func init() {
	testAccProviders = map[string]func() (*schema.Provider, error){
		"hetznerdns": ProviderFactory,
	}
}

func ProviderFactory() (*schema.Provider, error) {
	return Provider(), nil
}

// See https://www.terraform.io/docs/plugins/provider.html
func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// The Provider requires the API Token in env and thus it is required
// to run the acceptance test as well. This function is used as a PreCheck
// in TestCases and if the the token is not in env, it prints a message.
func testAccAPITokenPresent(t *testing.T) {
	if v := os.Getenv("HETZNER_DNS_API_TOKEN"); v == "" {
		t.Fatal("HETZNER_DNS_API_TOKEN must be set for acceptance tests")
	}
}
