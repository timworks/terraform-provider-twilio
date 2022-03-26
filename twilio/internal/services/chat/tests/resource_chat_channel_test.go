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

var channelResourceName = "twilio_chat_channel"

func TestAccTwilioChatChannel_basic(t *testing.T) {
	stateResourceName := fmt.Sprintf("%s.channel", channelResourceName)
	friendlyName := acctest.RandString(10)
	channelType := "private"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTwilioChatChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTwilioChatChannel_basic(friendlyName, channelType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioChatChannelExists(stateResourceName),
					resource.TestCheckResourceAttr(stateResourceName, "friendly_name", friendlyName),
					resource.TestCheckResourceAttrSet(stateResourceName, "id"),
					resource.TestCheckResourceAttrSet(stateResourceName, "sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "service_sid"),
					resource.TestCheckResourceAttr(stateResourceName, "unique_name", ""),
					resource.TestCheckResourceAttr(stateResourceName, "attributes", "{}"),
					resource.TestCheckResourceAttr(stateResourceName, "type", channelType),
					resource.TestCheckResourceAttrSet(stateResourceName, "created_by"),
					resource.TestCheckResourceAttrSet(stateResourceName, "members_count"),
					resource.TestCheckResourceAttrSet(stateResourceName, "messages_count"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateResourceName, "url"),
				),
			},
			{
				ResourceName:      stateResourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTwilioChatChannelImportStateIdFunc(stateResourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTwilioChatChannel_update(t *testing.T) {
	stateResourceName := fmt.Sprintf("%s.channel", channelResourceName)
	friendlyName := acctest.RandString(10)
	newFriendlyName := acctest.RandString(10)
	channelType := "private"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTwilioChatChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTwilioChatChannel_basic(friendlyName, channelType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioChatChannelExists(stateResourceName),
					resource.TestCheckResourceAttr(stateResourceName, "friendly_name", friendlyName),
					resource.TestCheckResourceAttrSet(stateResourceName, "id"),
					resource.TestCheckResourceAttrSet(stateResourceName, "sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "service_sid"),
					resource.TestCheckResourceAttr(stateResourceName, "unique_name", ""),
					resource.TestCheckResourceAttr(stateResourceName, "attributes", "{}"),
					resource.TestCheckResourceAttr(stateResourceName, "type", channelType),
					resource.TestCheckResourceAttrSet(stateResourceName, "created_by"),
					resource.TestCheckResourceAttrSet(stateResourceName, "members_count"),
					resource.TestCheckResourceAttrSet(stateResourceName, "messages_count"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateResourceName, "url"),
				),
			},
			{
				Config: testAccTwilioChatChannel_basic(newFriendlyName, channelType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioChatChannelExists(stateResourceName),
					resource.TestCheckResourceAttr(stateResourceName, "friendly_name", newFriendlyName),
					resource.TestCheckResourceAttrSet(stateResourceName, "id"),
					resource.TestCheckResourceAttrSet(stateResourceName, "sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "service_sid"),
					resource.TestCheckResourceAttr(stateResourceName, "unique_name", ""),
					resource.TestCheckResourceAttr(stateResourceName, "attributes", "{}"),
					resource.TestCheckResourceAttr(stateResourceName, "type", channelType),
					resource.TestCheckResourceAttrSet(stateResourceName, "created_by"),
					resource.TestCheckResourceAttrSet(stateResourceName, "members_count"),
					resource.TestCheckResourceAttrSet(stateResourceName, "messages_count"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateResourceName, "url"),
				),
			},
		},
	})
}

func TestAccTwilioChatChannel_invalidType(t *testing.T) {
	friendlyName := acctest.RandString(10)
	channelType := "test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTwilioChatChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccTwilioChatChannel_basic(friendlyName, channelType),
				ExpectError: regexp.MustCompile(`(?s)expected type to be one of \[public private\], got test`),
			},
		},
	})
}

func TestAccTwilioChatChannel_invalidServiceSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccTwilioChatChannel_invalidServiceSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of service_sid to match regular expression "\^IS\[0-9a-fA-F\]\{32\}\$", got service_sid`),
			},
		},
	})
}

func testAccCheckTwilioChatChannelDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*common.TwilioClient).Chat

	for _, rs := range s.RootModule().Resources {
		if rs.Type != channelResourceName {
			continue
		}

		if _, err := client.Service(rs.Primary.Attributes["service_sid"]).Channel(rs.Primary.ID).Fetch(); err != nil {
			if twilioError, ok := err.(*sdkUtils.TwilioError); ok {
				// currently programmable chat returns a 403 if the service instance does not exist
				if twilioError.Status == 403 && twilioError.Message == "Service instance not found" {
					return nil
				}
			}
			return fmt.Errorf("Error occurred when retrieving channel information %s", err.Error())
		}
	}

	return nil
}

func testAccCheckTwilioChatChannelExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*common.TwilioClient).Chat

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if _, err := client.Service(rs.Primary.Attributes["service_sid"]).Channel(rs.Primary.ID).Fetch(); err != nil {
			return fmt.Errorf("Error occurred when retrieving channel information %s", err.Error())
		}

		return nil
	}
}

func testAccTwilioChatChannelImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Not found: %s", name)
		}

		return fmt.Sprintf("/Services/%s/Channels/%s", rs.Primary.Attributes["service_sid"], rs.Primary.Attributes["sid"]), nil
	}
}

func testAccTwilioChatChannel_basic(friendlyName string, channelType string) string {
	return fmt.Sprintf(`
resource "twilio_chat_service" "service" {
  friendly_name = "%[1]s"
}

resource "twilio_chat_channel" "channel" {
  service_sid   = twilio_chat_service.service.sid
  friendly_name = "%[1]s"
  type          = "%[2]s"
}
`, friendlyName, channelType)
}

func testAccTwilioChatChannel_invalidServiceSid() string {
	return `
resource "twilio_chat_channel" "channel" {
  service_sid   = "service_sid"
  friendly_name = "invalid_service_sid"
  type          = "private"
}
`
}
