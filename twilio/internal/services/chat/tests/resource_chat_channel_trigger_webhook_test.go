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

var channelTriggerWebhookResourceName = "twilio_chat_channel_trigger_webhook"

func TestAccTwilioChatChannelTriggerWebhook_basic(t *testing.T) {
	stateResourceName := fmt.Sprintf("%s.trigger_webhook", channelTriggerWebhookResourceName)
	friendlyName := acctest.RandString(10)
	webhookURL := "https://localhost.com/webhook"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTwilioChatChannelTriggerWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTwilioChatChannelTriggerWebhook_basic(friendlyName, webhookURL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioChatChannelTriggerWebhookExists(stateResourceName),
					resource.TestCheckResourceAttr(stateResourceName, "type", "trigger"),
					resource.TestCheckResourceAttr(stateResourceName, "webhook_url", webhookURL),
					resource.TestCheckResourceAttr(stateResourceName, "triggers.#", "1"),
					resource.TestCheckResourceAttr(stateResourceName, "triggers.0", "keyword"),
					resource.TestCheckResourceAttrSet(stateResourceName, "id"),
					resource.TestCheckResourceAttrSet(stateResourceName, "sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "service_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "channel_sid"),
					resource.TestCheckResourceAttr(stateResourceName, "method", "POST"),
					resource.TestCheckResourceAttr(stateResourceName, "retry_count", "0"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateResourceName, "url"),
				),
			},
			{
				ResourceName:      stateResourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTwilioChatChannelStudioWebhookImportStateIdFunc(stateResourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTwilioChatChannelTriggerWebhook_update(t *testing.T) {
	stateResourceName := fmt.Sprintf("%s.trigger_webhook", channelTriggerWebhookResourceName)
	friendlyName := acctest.RandString(10)
	webhookURL := "https://localhost.com/webhook"
	newWebhookURL := "https://localhost.com/new"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTwilioChatChannelTriggerWebhookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTwilioChatChannelTriggerWebhook_basic(friendlyName, webhookURL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioChatChannelTriggerWebhookExists(stateResourceName),
					resource.TestCheckResourceAttr(stateResourceName, "type", "trigger"),
					resource.TestCheckResourceAttr(stateResourceName, "webhook_url", webhookURL),
					resource.TestCheckResourceAttr(stateResourceName, "triggers.#", "1"),
					resource.TestCheckResourceAttr(stateResourceName, "triggers.0", "keyword"),
					resource.TestCheckResourceAttrSet(stateResourceName, "id"),
					resource.TestCheckResourceAttrSet(stateResourceName, "sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "service_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "channel_sid"),
					resource.TestCheckResourceAttr(stateResourceName, "method", "POST"),
					resource.TestCheckResourceAttr(stateResourceName, "retry_count", "0"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateResourceName, "url"),
				),
			},
			{
				Config: testAccTwilioChatChannelTriggerWebhook_basic(friendlyName, newWebhookURL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioChatChannelTriggerWebhookExists(stateResourceName),
					resource.TestCheckResourceAttr(stateResourceName, "type", "trigger"),
					resource.TestCheckResourceAttr(stateResourceName, "webhook_url", newWebhookURL),
					resource.TestCheckResourceAttr(stateResourceName, "triggers.#", "1"),
					resource.TestCheckResourceAttr(stateResourceName, "triggers.0", "keyword"),
					resource.TestCheckResourceAttrSet(stateResourceName, "id"),
					resource.TestCheckResourceAttrSet(stateResourceName, "sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "service_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "channel_sid"),
					resource.TestCheckResourceAttr(stateResourceName, "method", "POST"),
					resource.TestCheckResourceAttr(stateResourceName, "retry_count", "0"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateResourceName, "url"),
				),
			},
		},
	})
}

func TestAccTwilioChatChannelTriggerWebhook_invalidWebhookURL(t *testing.T) {
	friendlyName := acctest.RandString(10)
	webhookURL := "webhook"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccTwilioChatChannelTriggerWebhook_basic(friendlyName, webhookURL),
				ExpectError: regexp.MustCompile(`(?s)expected "webhook_url" to have a host, got webhook`),
			},
		},
	})
}

func TestAccTwilioChatChannelTriggerWebhook_invalidServiceSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccTwilioChatChannelTriggerWebhook_invalidServiceSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of service_sid to match regular expression "\^IS\[0-9a-fA-F\]\{32\}\$", got service_sid`),
			},
		},
	})
}

func TestAccTwilioChatChannelTriggerWebhook_invalidChannelSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccTwilioChatChannelTriggerWebhook_invalidChannelSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of channel_sid to match regular expression "\^CH\[0-9a-fA-F\]\{32\}\$", got channel_sid`),
			},
		},
	})
}

func testAccCheckTwilioChatChannelTriggerWebhookDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*common.TwilioClient).Chat

	for _, rs := range s.RootModule().Resources {
		if rs.Type != channelTriggerWebhookResourceName {
			continue
		}

		if _, err := client.Service(rs.Primary.Attributes["service_sid"]).Channel(rs.Primary.Attributes["channel_sid"]).Webhook(rs.Primary.ID).Fetch(); err != nil {
			if twilioError, ok := err.(*sdkUtils.TwilioError); ok {
				// currently programmable chat returns a 403 if the service instance does not exist
				if twilioError.Status == 403 && twilioError.Message == "Service instance not found" {
					return nil
				}
			}
			return fmt.Errorf("Error occurred when retrieving channel webhook information %s", err.Error())
		}
	}

	return nil
}

func testAccCheckTwilioChatChannelTriggerWebhookExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*common.TwilioClient).Chat

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if _, err := client.Service(rs.Primary.Attributes["service_sid"]).Channel(rs.Primary.Attributes["channel_sid"]).Webhook(rs.Primary.ID).Fetch(); err != nil {
			return fmt.Errorf("Error occurred when retrieving channel webhook information %s", err.Error())
		}

		return nil
	}
}

func testAccTwilioChatChannelTriggerWebhookImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Not found: %s", name)
		}

		return fmt.Sprintf("/Services/%s/Channels/%s/Webhooks/%s", rs.Primary.Attributes["service_sid"], rs.Primary.Attributes["channel_sid"], rs.Primary.Attributes["sid"]), nil
	}
}

func testAccTwilioChatChannelTriggerWebhook_basic(friendlyName string, webhookUrl string) string {
	return fmt.Sprintf(`
resource "twilio_chat_service" "service" {
  friendly_name = "%[1]s"
}

resource "twilio_chat_channel" "channel" {
  service_sid   = twilio_chat_service.service.sid
  friendly_name = "%[1]s"
  type          = "private"
}

resource "twilio_chat_channel_trigger_webhook" "trigger_webhook" {
  service_sid = twilio_chat_service.service.sid
  channel_sid = twilio_chat_channel.channel.sid
  webhook_url = "%[2]s"
  triggers    = ["keyword"]
}
`, friendlyName, webhookUrl)
}

func testAccTwilioChatChannelTriggerWebhook_invalidServiceSid() string {
	return `
resource "twilio_chat_channel_trigger_webhook" "trigger_webhook" {
  service_sid = "service_sid"
  channel_sid = "CHaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  webhook_url = "https://localhost.com/webhook"
  triggers    = ["keyword"]
}
`
}

func testAccTwilioChatChannelTriggerWebhook_invalidChannelSid() string {
	return `
resource "twilio_chat_channel_trigger_webhook" "trigger_webhook" {
  service_sid = "ISaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  channel_sid = "channel_sid"
  webhook_url = "https://localhost.com/webhook"
  triggers    = ["keyword"]
}
`
}
