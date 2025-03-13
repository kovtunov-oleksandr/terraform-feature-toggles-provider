package feature_toggles

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Unit Tests

func TestResourceFeatureToggle_Schema(t *testing.T) {
	r := resourceFeatureToggle()

	// Test that required fields are marked as required
	if !r.Schema["name"].Required {
		t.Fatalf("name should be required")
	}
	if !r.Schema["enabled"].Required {
		t.Fatalf("enabled should be required")
	}
	if !r.Schema["environments"].Required {
		t.Fatalf("environments should be required")
	}

	// Test that optional fields are marked as optional
	if r.Schema["description"].Required {
		t.Fatalf("description should be optional")
	}
	if r.Schema["tags"].Required {
		t.Fatalf("tags should be optional")
	}
}

func TestResourceFeatureToggle_Create(t *testing.T) {
	r := resourceFeatureToggle()
	d := r.Data(nil)

	d.Set("name", "test-feature")
	d.Set("description", "Test feature toggle")
	d.Set("enabled", true)
	d.Set("tags", []string{"test", "example"})
	d.Set("environments", []string{"dev", "staging"})

	diags := resourceFeatureToggleCreate(context.Background(), d, nil)
	if diags.HasError() {
		t.Fatalf("got unexpected error: %#v", diags)
	}

	// Test that ID is set to the name
	if d.Id() != "test-feature" {
		t.Fatalf("expected ID to be 'test-feature', got %s", d.Id())
	}
}

func TestResourceFeatureToggle_Delete(t *testing.T) {
	r := resourceFeatureToggle()
	d := r.Data(nil)

	d.SetId("test-feature")

	diags := resourceFeatureToggleDelete(context.Background(), d, nil)
	if diags.HasError() {
		t.Fatalf("got unexpected error: %#v", diags)
	}

	// Test that ID is cleared
	if d.Id() != "" {
		t.Fatalf("expected ID to be empty, got %s", d.Id())
	}
}

// Acceptance Tests

func TestAccFeatureToggle_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFeatureToggleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureToggleConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeatureToggleExists("feature_toggles_feature.test"),
					resource.TestCheckResourceAttr("feature_toggles_feature.test", "name", "test_feature"),
					resource.TestCheckResourceAttr("feature_toggles_feature.test", "enabled", "true"),
					resource.TestCheckResourceAttr("feature_toggles_feature.test", "description", "Test feature toggle"),
					resource.TestCheckResourceAttr("feature_toggles_feature.test", "environments.#", "2"),
					resource.TestCheckResourceAttr("feature_toggles_feature.test", "environments.0", "dev"),
					resource.TestCheckResourceAttr("feature_toggles_feature.test", "environments.1", "staging"),
				),
			},
		},
	})
}

func TestAccFeatureToggle_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFeatureToggleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureToggleConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeatureToggleExists("feature_toggles_feature.test"),
					resource.TestCheckResourceAttr("feature_toggles_feature.test", "enabled", "true"),
				),
			},
			{
				Config: testAccFeatureToggleConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFeatureToggleExists("feature_toggles_feature.test"),
					resource.TestCheckResourceAttr("feature_toggles_feature.test", "enabled", "false"),
					resource.TestCheckResourceAttr("feature_toggles_feature.test", "description", "Updated test feature toggle"),
					resource.TestCheckResourceAttr("feature_toggles_feature.test", "tags.#", "3"),
					resource.TestCheckResourceAttr("feature_toggles_feature.test", "tags.0", "test"),
					resource.TestCheckResourceAttr("feature_toggles_feature.test", "tags.1", "example"),
					resource.TestCheckResourceAttr("feature_toggles_feature.test", "tags.2", "updated"),
				),
			},
		},
	})
}

// Helper functions and test fixtures

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"feature_toggles": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	// Add any necessary precheck logic here
}

func testAccCheckFeatureToggleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return nil
		}
		if rs.Primary.ID == "" {
			return nil
		}

		// Implement actual check logic here
		// For example, check if the feature toggle exists in the backend

		return nil
	}
}

func testAccCheckFeatureToggleDestroy(s *terraform.State) error {
	// Implement check logic to ensure the resource is destroyed,
	// This is typically done by checking that the resource no longer exists in the backend
	return nil
}

const testAccFeatureToggleConfig_basic = `
resource "feature_toggles_feature" "test" {
  name         = "test_feature"
  description  = "Test feature toggle"
  enabled      = true
  environments = ["dev", "staging"]
}
`

const testAccFeatureToggleConfig_update = `
resource "feature_toggles_feature" "test" {
  name         = "test_feature"
  description  = "Updated test feature toggle"
  enabled      = false
  tags         = ["test", "example", "updated"]
  environments = ["dev", "staging"]
}
`
