package tests

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/acceptance"
)

var webhookDataSourceName = "twilio_autopilot_webhook"

func TestAccDataSourceTwilioAutopilotWebhook_sid(t *testing.T) {
	stateDataSourceName := fmt.Sprintf("data.%s.webhook", webhookDataSourceName)
	uniqueName := acctest.RandString(10)
	url := "http://localhost/webhook"
	events := []string{"onDialogueStart", "onDialogueEnd"}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTwilioAutopilotWebhook_sid(uniqueName, url, events),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(stateDataSourceName, "unique_name", uniqueName),
					resource.TestCheckResourceAttr(stateDataSourceName, "events.#", "2"),
					resource.TestCheckResourceAttr(stateDataSourceName, "events.0", "onDialogueStart"),
					resource.TestCheckResourceAttr(stateDataSourceName, "events.1", "onDialogueEnd"),
					resource.TestCheckResourceAttr(stateDataSourceName, "webhook_url", url),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "id"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "assistant_sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "webhook_method"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "url"),
				),
			},
		},
	})
}

func TestAccDataSourceTwilioAutopilotWebhook_uniqueName(t *testing.T) {
	stateDataSourceName := fmt.Sprintf("data.%s.webhook", webhookDataSourceName)
	uniqueName := acctest.RandString(10)
	url := "http://localhost/webhook"
	events := []string{"onDialogueStart", "onDialogueEnd"}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTwilioAutopilotWebhook_uniqueName(uniqueName, url, events),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(stateDataSourceName, "unique_name", uniqueName),
					resource.TestCheckResourceAttr(stateDataSourceName, "events.#", "2"),
					resource.TestCheckResourceAttr(stateDataSourceName, "events.0", "onDialogueStart"),
					resource.TestCheckResourceAttr(stateDataSourceName, "events.1", "onDialogueEnd"),
					resource.TestCheckResourceAttr(stateDataSourceName, "webhook_url", url),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "id"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "assistant_sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "webhook_method"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "url"),
				),
			},
		},
	})
}

func TestAccDataSourceTwilioAutopilotWebhook_invalidAssistantSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceTwilioAutopilotWebhook_invalidAssistantSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of assistant_sid to match regular expression "\^UA\[0-9a-fA-F\]\{32\}\$", got assistant_sid`),
			},
		},
	})
}

func TestAccDataSourceTwilioAutopilotWebhook_invalidSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceTwilioAutopilotWebhook_invalidSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of sid to match regular expression "\^UM\[0-9a-fA-F\]\{32\}\$", got sid`),
			},
		},
	})
}

func testAccDataSourceTwilioAutopilotWebhook_sid(uniqueName string, url string, events []string) string {
	return fmt.Sprintf(`
resource "twilio_autopilot_assistant" "assistant" {
  unique_name = "%[1]s"
}

resource "twilio_autopilot_webhook" "webhook" {
  assistant_sid = twilio_autopilot_assistant.assistant.sid
  unique_name   = "%[1]s"
  webhook_url   = "%[2]s"
  events        = %[3]s
}

data "twilio_autopilot_webhook" "webhook" {
  assistant_sid = twilio_autopilot_webhook.webhook.assistant_sid
  sid           = twilio_autopilot_webhook.webhook.sid
}
`, uniqueName, url, `["`+strings.Join(events, `","`)+`"]`)
}

func testAccDataSourceTwilioAutopilotWebhook_uniqueName(uniqueName string, url string, events []string) string {
	return fmt.Sprintf(`
resource "twilio_autopilot_assistant" "assistant" {
  unique_name = "%[1]s"
}

resource "twilio_autopilot_webhook" "webhook" {
  assistant_sid = twilio_autopilot_assistant.assistant.sid
  unique_name   = "%[1]s"
  webhook_url   = "%[2]s"
  events        = %[3]s
}

data "twilio_autopilot_webhook" "webhook" {
  assistant_sid = twilio_autopilot_webhook.webhook.assistant_sid
  unique_name   = twilio_autopilot_webhook.webhook.unique_name
}
`, uniqueName, url, `["`+strings.Join(events, `","`)+`"]`)
}

func testAccDataSourceTwilioAutopilotWebhook_invalidAssistantSid() string {
	return `
data "twilio_autopilot_webhook" "webhook" {
  assistant_sid = "assistant_sid"
  sid           = "UMaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}
`
}

func testAccDataSourceTwilioAutopilotWebhook_invalidSid() string {
	return `
data "twilio_autopilot_webhook" "webhook" {
  assistant_sid = "UAaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  sid           = "sid"
}
`
}
