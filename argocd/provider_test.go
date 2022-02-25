package argocd

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]func() (*schema.Provider, error)

func init() {
	testAccProviders = map[string]func() (*schema.Provider, error){
		"argocd": func() (*schema.Provider, error) {
			return Provider(), nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ARGOCD_AUTH_USERNAME"); v == "" {
		t.Fatal("ARGOCD_AUTH_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("ARGOCD_AUTH_PASSWORD"); v == "" {
		t.Fatal("ARGOCD_AUTH_PASSWORD must be set for acceptance tests")
	}
	if v := os.Getenv("ARGOCD_SERVER"); v == "" {
		t.Fatal("ARGOCD_SERVER must be set for acceptance tests")
	}
	if v := os.Getenv("ARGOCD_INSECURE"); v == "" {
		t.Fatal("ARGOCD_INSECURE should be set for acceptance tests")
	}
}

func testAccPreCheckFeatureSupported(t *testing.T, feature int) {
	testAccProvider, err := testAccProviders["argocd"]()
	if err != nil {
		t.Fatal(err)
	}

	d := schema.TestResourceDataRaw(t, testAccProvider.Schema, make(map[string]interface{}))
	server := ServerInterface{ProviderData: d}
	err = server.initClients()
	if err != nil {
		t.Fatal(err)
	}

	featureSupported, err := server.isFeatureSupported(feature)
	if err != nil {
		t.Fatal(err)
	}

	if !featureSupported {
		t.Skip("feature not supported in tested argocd version")
	}
}
