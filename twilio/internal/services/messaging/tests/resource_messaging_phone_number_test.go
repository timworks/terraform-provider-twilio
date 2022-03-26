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

var phoneNumberResourceName = "twilio_messaging_phone_number"

func TestAccTwilioMessagingPhoneNumber_basic(t *testing.T) {
	stateResourceName := fmt.Sprintf("%s.phone_number", phoneNumberResourceName)
	friendlyName := acctest.RandString(10)
	testData := acceptance.TestAccData

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTwilioMessagingPhoneNumberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTwilioMessagingPhoneNumber_basic(testData, friendlyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioMessagingPhoneNumberExists(stateResourceName),
					resource.TestCheckResourceAttr(stateResourceName, "sid", testData.PhoneNumberSid),
					resource.TestCheckResourceAttrSet(stateResourceName, "id"),
					resource.TestCheckResourceAttrSet(stateResourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.#"),
					resource.TestCheckResourceAttrSet(stateResourceName, "country_code"),
					resource.TestCheckResourceAttrSet(stateResourceName, "phone_number"),
					resource.TestCheckResourceAttrSet(stateResourceName, "service_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateResourceName, "url"),
				),
			},
			{
				ResourceName:      stateResourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTwilioMessagingPhoneNumberImportStateIdFunc(stateResourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTwilioMessagingPhoneNumber_invalidServiceSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccTwilioMessagingPhoneNumber_invalidServiceSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of service_sid to match regular expression "\^MG\[0-9a-fA-F\]\{32\}\$", got service_sid`),
			},
		},
	})
}

func TestAccTwilioMessagingPhoneNumber_invalidSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccTwilioMessagingPhoneNumber_invalidSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of sid to match regular expression "\^PN\[0-9a-fA-F\]\{32\}\$", got sid`),
			},
		},
	})
}

func testAccCheckTwilioMessagingPhoneNumberDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*common.TwilioClient).Messaging

	for _, rs := range s.RootModule().Resources {
		if rs.Type != phoneNumberResourceName {
			continue
		}

		if _, err := client.Service(rs.Primary.Attributes["service_sid"]).PhoneNumber(rs.Primary.ID).Fetch(); err != nil {
			if utils.IsNotFoundError(err) {
				return nil
			}
			return fmt.Errorf("Error occurred when retrieving phone number information %s", err.Error())
		}
	}

	return nil
}

func testAccCheckTwilioMessagingPhoneNumberExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*common.TwilioClient).Messaging

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if _, err := client.Service(rs.Primary.Attributes["service_sid"]).PhoneNumber(rs.Primary.ID).Fetch(); err != nil {
			return fmt.Errorf("Error occurred when retrieving phone number information %s", err.Error())
		}

		return nil
	}
}

func testAccTwilioMessagingPhoneNumberImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Not found: %s", name)
		}

		return fmt.Sprintf("/Services/%s/PhoneNumbers/%s", rs.Primary.Attributes["service_sid"], rs.Primary.Attributes["sid"]), nil
	}
}

func testAccTwilioMessagingPhoneNumber_basic(testData *acceptance.TestData, friendlyName string) string {
	return fmt.Sprintf(`
resource "twilio_messaging_service" "service" {
  friendly_name = "service-%s"
}

resource "twilio_messaging_phone_number" "phone_number" {
  service_sid = twilio_messaging_service.service.sid
  sid         = "%s"
}
`, friendlyName, testData.PhoneNumberSid)
}

func testAccTwilioMessagingPhoneNumber_invalidServiceSid() string {
	return `
resource "twilio_messaging_phone_number" "phone_number" {
  service_sid = "service_sid"
  sid         = "PNaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}
`
}

func testAccTwilioMessagingPhoneNumber_invalidSid() string {
	return `
resource "twilio_messaging_phone_number" "phone_number" {
  service_sid = "MGaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  sid         = "sid"
}
`
}
