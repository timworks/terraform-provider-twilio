package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/acceptance"
)

var proxyServiceDataSourceName = "twilio_proxy_service"

func TestAccDataSourceTwilioProxyService_basic(t *testing.T) {
	stateDataSourceName := fmt.Sprintf("data.%s.service", proxyServiceDataSourceName)
	uniqueName := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTwilioProxyService_basic(uniqueName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(stateDataSourceName, "unique_name", uniqueName),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "id"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "account_sid"),
					resource.TestCheckResourceAttr(stateDataSourceName, "chat_instance_sid", ""),
					resource.TestCheckResourceAttr(stateDataSourceName, "default_ttl", "0"),
					resource.TestCheckResourceAttr(stateDataSourceName, "callback_url", ""),
					resource.TestCheckResourceAttr(stateDataSourceName, "geo_match_level", "country"),
					resource.TestCheckResourceAttr(stateDataSourceName, "number_selection_behavior", "prefer-sticky"),
					resource.TestCheckResourceAttr(stateDataSourceName, "intercept_callback_url", ""),
					resource.TestCheckResourceAttr(stateDataSourceName, "out_of_session_callback_url", ""),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "url"),
				),
			},
		},
	})
}

func TestAccDataSourceTwilioProxyService_invalidSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceTwilioProxyService_invalidSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of sid to match regular expression "\^KS\[0-9a-fA-F\]\{32\}\$", got sid`),
			},
		},
	})
}

func testAccDataSourceTwilioProxyService_basic(uniqueName string) string {
	return fmt.Sprintf(`
resource "twilio_proxy_service" "service" {
  unique_name = "%s"
}

data "twilio_proxy_service" "service" {
  sid = twilio_proxy_service.service.sid
}
`, uniqueName)
}

func testAccDataSourceTwilioProxyService_invalidSid() string {
	return `
data "twilio_proxy_service" "service" {
  sid = "sid"
}
`
}
