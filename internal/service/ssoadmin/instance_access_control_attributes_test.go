package ssoadmin_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ssoadmin"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfssoadmin "github.com/hashicorp/terraform-provider-aws/internal/service/ssoadmin"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccSSOAdmin_serial(t *testing.T) {
	testCases := map[string]map[string]func(t *testing.T){
		"InstanceAccessControlAttributes": {
			"basic":      testAccInstanceAccessControlAttributes_basic,
			"disappears": testAccInstanceAccessControlAttributes_disappears,
			"multiple":   testAccInstanceAccessControlAttributes_multiple,
			"update":     testAccInstanceAccessControlAttributes_update,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccInstanceAccessControlAttributes_basic(t *testing.T) {
	resourceName := "aws_ssoadmin_instance_access_control_attributes.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckInstances(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ssoadmin.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckInstanceAccessControlAttributesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceAccessControlAttributesConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceAccessControlAttributesExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "attribute.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", "ENABLED"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccInstanceAccessControlAttributes_disappears(t *testing.T) {
	resourceName := "aws_ssoadmin_instance_access_control_attributes.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckInstances(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ssoadmin.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPermissionSetInlinePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceAccessControlAttributesConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceAccessControlAttributesExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfssoadmin.ResourceAccessControlAttributes(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccInstanceAccessControlAttributes_multiple(t *testing.T) {
	resourceName := "aws_ssoadmin_instance_access_control_attributes.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckInstances(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ssoadmin.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckInstanceAccessControlAttributesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceAccessControlAttributesConfig_multiple(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceAccessControlAttributesExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "attribute.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "status", "ENABLED"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccInstanceAccessControlAttributes_update(t *testing.T) {
	resourceName := "aws_ssoadmin_instance_access_control_attributes.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); testAccPreCheckInstances(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ssoadmin.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckInstanceAccessControlAttributesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceAccessControlAttributesConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceAccessControlAttributesExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccInstanceAccessControlAttributesConfig_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceAccessControlAttributesExists(resourceName),
				),
			},
		},
	})
}

func testAccCheckInstanceAccessControlAttributesDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).SSOAdminConn()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_ssoadmin_instance_access_control_attributes" {
			continue
		}

		_, err := tfssoadmin.FindInstanceAttributeControlAttributesByARN(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("SSO Instance Access Control Attributes %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckInstanceAccessControlAttributesExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SSO Instance Access Control Attributes ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SSOAdminConn()

		_, err := tfssoadmin.FindInstanceAttributeControlAttributesByARN(conn, rs.Primary.ID)

		return err
	}
}

func testAccInstanceAccessControlAttributesConfig_basic() string {
	return `
data "aws_ssoadmin_instances" "test" {}

resource "aws_ssoadmin_instance_access_control_attributes" "test" {
  instance_arn = tolist(data.aws_ssoadmin_instances.test.arns)[0]
  attribute {
    key = "name"
    value {
      source = ["$${path:name.givenName}"]
    }
  }
}
`
}
func testAccInstanceAccessControlAttributesConfig_multiple() string {
	return `
data "aws_ssoadmin_instances" "test" {}

resource "aws_ssoadmin_instance_access_control_attributes" "test" {
  instance_arn = tolist(data.aws_ssoadmin_instances.test.arns)[0]
  attribute {
    key = "name"
    value {
      source = ["$${path:name.givenName}"]
    }
  }
  attribute {
    key = "last"
    value {
      source = ["$${path:name.familyName}"]
    }
  }
}
`
}

func testAccInstanceAccessControlAttributesConfig_update() string {
	return `
data "aws_ssoadmin_instances" "test" {}

resource "aws_ssoadmin_instance_access_control_attributes" "test" {
  instance_arn = tolist(data.aws_ssoadmin_instances.test.arns)[0]
  attribute {
    key = "name"
    value {
      source = ["$${path:name.familyName}"]
    }
  }
}
`
}
