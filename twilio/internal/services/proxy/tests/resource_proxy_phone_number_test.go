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
	sdkUtils "github.com/timworks/twilio-sdk-go/utils"
)

var proxyPhoneNumberResourceName = "twilio_proxy_phone_number"

// Tests have to run sequentially as a phone number cannot be associated with more than 1 proxy service at a given time

func TestAccTwilioProxyPhoneNumber_basic(t *testing.T) {
	stateResourceName := fmt.Sprintf("%s.phone_number", proxyPhoneNumberResourceName)

	testData := acceptance.TestAccData
	uniqueName := acctest.RandString(10)
	isReserved := true

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTwilioProxyPhoneNumberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTwilioProxyPhoneNumber_basic(testData, uniqueName, isReserved),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioProxyPhoneNumberExists(stateResourceName),
					resource.TestCheckResourceAttr(stateResourceName, "sid", testData.PhoneNumberSid),
					resource.TestCheckResourceAttrSet(stateResourceName, "id"),
					resource.TestCheckResourceAttrSet(stateResourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "service_sid"),
					resource.TestCheckResourceAttr(stateResourceName, "is_reserved", "true"),
					resource.TestCheckResourceAttrSet(stateResourceName, "phone_number"),
					resource.TestCheckResourceAttrSet(stateResourceName, "friendly_name"),
					resource.TestCheckResourceAttrSet(stateResourceName, "iso_country"),
					resource.TestCheckResourceAttrSet(stateResourceName, "in_use"),
					resource.TestCheckResourceAttr(stateResourceName, "capabilities.#", "1"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.fax_inbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.fax_outbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.mms_inbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.mms_outbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.restriction_fax_domestic"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.restriction_mms_domestic"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.restriction_sms_domestic"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.restriction_voice_domestic"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.sip_trunking"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.sms_inbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.sms_outbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.voice_inbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.voice_outbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateResourceName, "url"),
				),
			},
			{
				ResourceName:      stateResourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTwilioProxyPhoneNumberImportStateIdFunc(stateResourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTwilioProxyPhoneNumber_update(t *testing.T) {
	stateResourceName := fmt.Sprintf("%s.phone_number", proxyPhoneNumberResourceName)

	testData := acceptance.TestAccData
	uniqueName := acctest.RandString(10)
	isReserved := true
	newIsReserved := false

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTwilioProxyPhoneNumberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTwilioProxyPhoneNumber_basic(testData, uniqueName, isReserved),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioProxyPhoneNumberExists(stateResourceName),
					resource.TestCheckResourceAttr(stateResourceName, "sid", testData.PhoneNumberSid),
					resource.TestCheckResourceAttrSet(stateResourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "service_sid"),
					resource.TestCheckResourceAttr(stateResourceName, "is_reserved", "true"),
					resource.TestCheckResourceAttrSet(stateResourceName, "phone_number"),
					resource.TestCheckResourceAttrSet(stateResourceName, "friendly_name"),
					resource.TestCheckResourceAttrSet(stateResourceName, "iso_country"),
					resource.TestCheckResourceAttrSet(stateResourceName, "in_use"),
					resource.TestCheckResourceAttr(stateResourceName, "capabilities.#", "1"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.fax_inbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.fax_outbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.mms_inbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.mms_outbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.restriction_fax_domestic"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.restriction_mms_domestic"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.restriction_sms_domestic"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.restriction_voice_domestic"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.sip_trunking"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.sms_inbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.sms_outbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.voice_inbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.voice_outbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateResourceName, "url"),
				),
			},
			{
				Config: testAccTwilioProxyPhoneNumber_basic(acceptance.TestAccData, uniqueName, newIsReserved),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioProxyPhoneNumberExists(stateResourceName),
					resource.TestCheckResourceAttr(stateResourceName, "sid", testData.PhoneNumberSid),
					resource.TestCheckResourceAttrSet(stateResourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "service_sid"),
					resource.TestCheckResourceAttr(stateResourceName, "is_reserved", "false"),
					resource.TestCheckResourceAttrSet(stateResourceName, "phone_number"),
					resource.TestCheckResourceAttrSet(stateResourceName, "friendly_name"),
					resource.TestCheckResourceAttrSet(stateResourceName, "iso_country"),
					resource.TestCheckResourceAttrSet(stateResourceName, "in_use"),
					resource.TestCheckResourceAttr(stateResourceName, "capabilities.#", "1"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.fax_inbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.fax_outbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.mms_inbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.mms_outbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.restriction_fax_domestic"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.restriction_mms_domestic"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.restriction_sms_domestic"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.restriction_voice_domestic"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.sip_trunking"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.sms_inbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.sms_outbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.voice_inbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "capabilities.0.voice_outbound"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateResourceName, "url"),
				),
			},
		},
	})
}

func TestAccTwilioProxyPhoneNumber_invalidServiceSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccTwilioProxyPhoneNumber_invalidServiceSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of service_sid to match regular expression "\^KS\[0-9a-fA-F\]\{32\}\$", got service_sid`),
			},
		},
	})
}

func TestAccTwilioProxyPhoneNumber_invalidPhoneNumberSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccTwilioProxyPhoneNumber_invalidPhoneNumberSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of sid to match regular expression "\^PN\[0-9a-fA-F\]\{32\}\$", got phone_number_sid`),
			},
		},
	})
}

func TestAccTwilioProxyPhoneNumber_invalidPhoneNumber(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccTwilioProxyPhoneNumber_invalidPhoneNumber(),
				ExpectError: regexp.MustCompile(`(?s)expected value of phone_number to match regular expression "\^\\\\\+\[1-9\]\\\\d\{1,14\}\$", got phone_number`),
			},
		},
	})
}

func testAccCheckTwilioProxyPhoneNumberDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*common.TwilioClient).Proxy

	for _, rs := range s.RootModule().Resources {
		if rs.Type != proxyPhoneNumberResourceName {
			continue
		}

		if _, err := client.Service(rs.Primary.Attributes["service_sid"]).PhoneNumber(rs.Primary.ID).Fetch(); err != nil {
			if twilioError, ok := err.(*sdkUtils.TwilioError); ok {
				// currently proxy returns a 400 if the proxy phone number instance does not exist
				if twilioError.Status == 400 && twilioError.Message == "Invalid Phone Number Sid" {
					return nil
				}
			}
			return fmt.Errorf("Error occurred when retrieving proxy phone number information %s", err.Error())
		}
	}

	return nil
}

func testAccCheckTwilioProxyPhoneNumberExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*common.TwilioClient).Proxy

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if _, err := client.Service(rs.Primary.Attributes["service_sid"]).PhoneNumber(rs.Primary.ID).Fetch(); err != nil {
			return fmt.Errorf("Error occurred when retrieving proxy phone number information %s", err.Error())
		}

		return nil
	}
}

func testAccTwilioProxyPhoneNumberImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Not found: %s", name)
		}

		return fmt.Sprintf("/Services/%s/PhoneNumbers/%s", rs.Primary.Attributes["service_sid"], rs.Primary.Attributes["sid"]), nil
	}
}

func testAccTwilioProxyPhoneNumber_basic(testData *acceptance.TestData, uniqueName string, isReserved bool) string {
	return fmt.Sprintf(`
resource "twilio_proxy_service" "service" {
  unique_name = "%s"
}

resource "twilio_proxy_phone_number" "phone_number" {
  service_sid = twilio_proxy_service.service.sid
  sid         = "%s"
  is_reserved = %t
}
`, uniqueName, testData.PhoneNumberSid, isReserved)
}

func testAccTwilioProxyPhoneNumber_invalidServiceSid() string {
	return `
resource "twilio_proxy_phone_number" "phone_number" {
  service_sid = "service_sid"
  sid         = "PNaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  is_reserved = false
}
`
}

func testAccTwilioProxyPhoneNumber_invalidPhoneNumberSid() string {
	return `
resource "twilio_proxy_phone_number" "phone_number" {
  service_sid = "KSaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  sid         = "phone_number_sid"
  is_reserved = false
}
`
}

func testAccTwilioProxyPhoneNumber_invalidPhoneNumber() string {
	return `
resource "twilio_proxy_phone_number" "phone_number" {
  service_sid = "KSaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  phone_number         = "phone_number"
  is_reserved = false
}
`
}
