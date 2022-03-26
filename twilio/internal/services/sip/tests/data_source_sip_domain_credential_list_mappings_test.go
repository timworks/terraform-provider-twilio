package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/acceptance"
)

var domainCredentialListMappingsDataSourceName = "twilio_sip_domain_credential_list_mappings"

func TestAccDataSourceTwilioSIPDomainCredentialListMappings_basic(t *testing.T) {
	stateDataSourceName := fmt.Sprintf("data.%s.credential_list_mappings", domainCredentialListMappingsDataSourceName)

	testData := acceptance.TestAccData
	friendlyName := acctest.RandString(10)
	domainName := acctest.RandString(10) + ".sip.twilio.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTwilioSIPDomainCredentialListMappings_basic(testData, friendlyName, domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(stateDataSourceName, "id"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "domain_sid"),
					resource.TestCheckResourceAttr(stateDataSourceName, "credential_list_mappings.#", "1"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "credential_list_mappings.0.sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "credential_list_mappings.0.friendly_name"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "credential_list_mappings.0.date_created"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "credential_list_mappings.0.date_updated"),
				),
			},
		},
	})
}

func TestAccDataSourceTwilioSIPDomainCredentialListMappings_invalidAccountSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceTwilioSIPDomainCredentialListMappings_invalidAccountSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of account_sid to match regular expression "\^AC\[0-9a-fA-F\]\{32\}\$", got account_sid`),
			},
		},
	})
}

func TestAccDataSourceTwilioSIPDomainCredentialListMappings_invalidDomainSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceTwilioSIPDomainCredentialListMappings_invalidDomainSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of domain_sid to match regular expression "\^SD\[0-9a-fA-F\]\{32\}\$", got domain_sid`),
			},
		},
	})
}

func testAccDataSourceTwilioSIPDomainCredentialListMappings_basic(testData *acceptance.TestData, friendlyName string, domainName string) string {
	return fmt.Sprintf(`
resource "twilio_sip_credential_list" "credential_list" {
  account_sid   = "%[1]s"
  friendly_name = "%[2]s"
}

resource "twilio_sip_domain" "domain" {
  account_sid = "%[1]s"
  domain_name = "%[3]s"
}

resource "twilio_sip_domain_credential_list_mapping" "credential_list_mapping" {
  account_sid         = twilio_sip_domain.domain.account_sid
  domain_sid          = twilio_sip_domain.domain.sid
  credential_list_sid = twilio_sip_credential_list.credential_list.sid
}

data "twilio_sip_domain_credential_list_mappings" "credential_list_mappings" {
  account_sid = twilio_sip_domain_credential_list_mapping.credential_list_mapping.account_sid
  domain_sid  = twilio_sip_domain_credential_list_mapping.credential_list_mapping.domain_sid
}
`, testData.AccountSid, friendlyName, domainName)
}

func testAccDataSourceTwilioSIPDomainCredentialListMappings_invalidAccountSid() string {
	return `
data "twilio_sip_domain_credential_list_mappings" "credential_list_mappings" {
  account_sid = "account_sid"
  domain_sid  = "SDaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}
`
}

func testAccDataSourceTwilioSIPDomainCredentialListMappings_invalidDomainSid() string {
	return `
data "twilio_sip_domain_credential_list_mappings" "credential_list_mappings" {
  account_sid = "ACaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  domain_sid  = "domain_sid"
}
`
}
