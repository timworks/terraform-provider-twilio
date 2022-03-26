package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/timworks/terraform-provider-twilio/twilio/common"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/acceptance"
	"github.com/timworks/terraform-provider-twilio/twilio/utils"
)

var userResourceName = "twilio_conversations_user"

func TestAccTwilioConversationsUser_basic(t *testing.T) {
	stateResourceName := fmt.Sprintf("%s.user", userResourceName)
	friendlyName := acctest.RandString(10)
	identity := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTwilioConversationsUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTwilioConversationsUser_basic(friendlyName, identity),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioConversationsUserExists(stateResourceName),
					resource.TestCheckResourceAttr(stateResourceName, "identity", identity),
					resource.TestCheckResourceAttrSet(stateResourceName, "id"),
					resource.TestCheckResourceAttrSet(stateResourceName, "sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "service_sid"),
					resource.TestCheckResourceAttr(stateResourceName, "friendly_name", ""),
					resource.TestCheckResourceAttr(stateResourceName, "attributes", "{}"),
					resource.TestCheckResourceAttrSet(stateResourceName, "is_notifiable"),
					resource.TestCheckResourceAttrSet(stateResourceName, "is_online"),
					resource.TestCheckResourceAttrSet(stateResourceName, "role_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateResourceName, "url"),
				),
			},
			{
				ResourceName:      stateResourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTwilioConversationsUserImportStateIdFunc(stateResourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTwilioConversationsUser_friendlyName(t *testing.T) {
	stateResourceName := fmt.Sprintf("%s.user", userResourceName)

	friendlyName := acctest.RandString(10)
	userFriendlyName := ""
	newUserFriendlyName := acctest.RandString(256)
	identity := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTwilioConversationsUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTwilioConversationsUser_friendlyName(friendlyName, userFriendlyName, identity),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioConversationsUserExists(stateResourceName),
					resource.TestCheckResourceAttrSet(stateResourceName, "service_sid"),
					resource.TestCheckResourceAttr(stateResourceName, "friendly_name", userFriendlyName),
				),
			},
			{
				Config: testAccTwilioConversationsUser_friendlyName(friendlyName, newUserFriendlyName, identity),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioConversationsUserExists(stateResourceName),
					resource.TestCheckResourceAttr(stateResourceName, "friendly_name", newUserFriendlyName),
				),
			},
			{
				Config: testAccTwilioConversationsUser_basic(friendlyName, identity),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioConversationsUserExists(stateResourceName),
					resource.TestCheckResourceAttr(stateResourceName, "friendly_name", ""),
				),
			},
		},
	})
}

func TestAccTwilioConversationsUser_invalidFriendlyNameWith257Characters(t *testing.T) {
	friendlyName := acctest.RandString(10)
	userFriendlyName := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	identity := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccTwilioConversationsUser_friendlyName(friendlyName, userFriendlyName, identity),
				ExpectError: regexp.MustCompile(`(?s)expected length of friendly_name to be in the range \(0 - 256\), got aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`),
			},
		},
	})
}

func TestAccTwilioConversationsUser_attibutes(t *testing.T) {
	stateResourceName := fmt.Sprintf("%s.user", userResourceName)
	friendlyName := acctest.RandString(10)
	identity := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTwilioConversationsUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTwilioConversationsUser_basic(friendlyName, identity),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioConversationsUserExists(stateResourceName),
					resource.TestCheckResourceAttr(stateResourceName, "attributes", "{}"),
				),
			},
			{
				Config: testAccTwilioConversationsUser_withAttributes(friendlyName, identity),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioConversationsUserExists(stateResourceName),
					resource.TestCheckResourceAttr(stateResourceName, "attributes", `{"test":true}`),
				),
			},
			{
				Config: testAccTwilioConversationsUser_basic(friendlyName, identity),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioConversationsUserExists(stateResourceName),
					resource.TestCheckResourceAttr(stateResourceName, "attributes", "{}"),
				),
			},
		},
	})
}

func TestAccTwilioConversationsUser_invalidAttributesString(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccTwilioConversationsUser_invalidAttributesString(),
				ExpectError: regexp.MustCompile(`(?s)"attributes" contains an invalid JSON`),
			},
		},
	})
}

func TestAccTwilioConversationsUser_invalidServiceSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccTwilioConversationsUser_invalidServiceSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of service_sid to match regular expression "\^IS\[0-9a-fA-F\]\{32\}\$", got service_sid`),
			},
		},
	})
}

func testAccCheckTwilioConversationsUserDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*common.TwilioClient).Conversations

	for _, rs := range s.RootModule().Resources {
		if rs.Type != userResourceName {
			continue
		}

		if _, err := client.Service(rs.Primary.Attributes["service_sid"]).User(rs.Primary.ID).Fetch(); err != nil {
			if utils.IsNotFoundError(err) {
				return nil
			}
			return fmt.Errorf("Error occurred when retrieving user information %s", err.Error())
		}
	}

	return nil
}

func testAccCheckTwilioConversationsUserExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*common.TwilioClient).Conversations

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if _, err := client.Service(rs.Primary.Attributes["service_sid"]).User(rs.Primary.ID).Fetch(); err != nil {
			return fmt.Errorf("Error occurred when retrieving user information %s", err.Error())
		}

		return nil
	}
}

func testAccTwilioConversationsUserImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Not found: %s", name)
		}

		return fmt.Sprintf("/Services/%s/Users/%s", rs.Primary.Attributes["service_sid"], rs.Primary.Attributes["sid"]), nil
	}
}

func testAccTwilioConversationsUser_basic(friendlyName string, identity string) string {
	return fmt.Sprintf(`
resource "twilio_conversations_service" "service" {
  friendly_name = "%s"
}

resource "twilio_conversations_user" "user" {
  service_sid = twilio_conversations_service.service.sid
  identity    = "%s"
}
`, friendlyName, identity)
}

func testAccTwilioConversationsUser_friendlyName(friendlyName string, userFriendlyName string, identity string) string {
	return fmt.Sprintf(`
resource "twilio_conversations_service" "service" {
  friendly_name = "%s"
}

resource "twilio_conversations_user" "user" {
  service_sid   = twilio_conversations_service.service.sid
  friendly_name = "%s"
  identity      = "%s"
}
`, friendlyName, userFriendlyName, identity)
}

func testAccTwilioConversationsUser_withAttributes(friendlyName string, identity string) string {
	return fmt.Sprintf(`
resource "twilio_conversations_service" "service" {
  friendly_name = "%s"
}

resource "twilio_conversations_user" "user" {
  service_sid = twilio_conversations_service.service.sid
  identity    = "%s"
  attributes  = "{\"test\": true}"
}
`, friendlyName, identity)
}

func testAccTwilioConversationsUser_invalidAttributesString() string {
	return `
resource "twilio_conversations_user" "user" {
  service_sid = "ISaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  identity    = "invalid_attributes"
  attributes  = "attributes"
}
`
}

func testAccTwilioConversationsUser_invalidServiceSid() string {
	return `
resource "twilio_conversations_user" "user" {
  service_sid = "service_sid"
  identity    = "invalid_service_sid"
}
`
}
