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

var domainIPAccessControlListMappingResourceName = "twilio_sip_domain_ip_access_control_list_mapping"

func TestAccTwilioSIPDomainIPAccessControlListMapping_basic(t *testing.T) {
	stateResourceName := fmt.Sprintf("%s.domain_ip_access_control_list_mapping", domainIPAccessControlListMappingResourceName)

	testData := acceptance.TestAccData
	friendlyName := acctest.RandString(10)
	domainName := acctest.RandString(10) + ".sip.twilio.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTwilioSIPDomainIPAccessControlListMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTwilioSIPDomainIPAccessControlListMapping_basic(testData, friendlyName, domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioSIPDomainIPAccessControlListMappingExists(stateResourceName),
					resource.TestCheckResourceAttrSet(stateResourceName, "id"),
					resource.TestCheckResourceAttrSet(stateResourceName, "sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "domain_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "ip_access_control_list_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "friendly_name"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_updated"),
				),
			},
			{
				ResourceName:      stateResourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTwilioSIPDomainIPAccessControlListMappingImportStateIdFunc(stateResourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTwilioSIPDomainIPAccessControlListMapping_invalidAccountSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccTwilioSIPDomainIPAccessControlListMapping_invalidAccountSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of account_sid to match regular expression "\^AC\[0-9a-fA-F\]\{32\}\$", got account_sid`),
			},
		},
	})
}

func TestAccTwilioSIPDomainIPAccessControlListMapping_invalidDomainSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccTwilioSIPDomainIPAccessControlListMapping_invalidDomainSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of domain_sid to match regular expression "\^SD\[0-9a-fA-F\]\{32\}\$", got domain_sid`),
			},
		},
	})
}

func TestAccTwilioSIPDomainIPAccessControlListMapping_invalidIPAccessControlListSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccTwilioSIPDomainIPAccessControlListMapping_invalidIPAccessControlList(),
				ExpectError: regexp.MustCompile(`(?s)expected value of ip_access_control_list_sid to match regular expression "\^AL\[0-9a-fA-F\]\{32\}\$", got ip_access_control_list_sid`),
			},
		},
	})
}

func testAccCheckTwilioSIPDomainIPAccessControlListMappingDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*common.TwilioClient).API

	for _, rs := range s.RootModule().Resources {
		if rs.Type != domainIPAccessControlListMappingResourceName {
			continue
		}

		if _, err := client.Account(rs.Primary.Attributes["account_sid"]).Sip.Domain(rs.Primary.Attributes["domain_sid"]).Auth.Calls.IpAccessControlListMapping(rs.Primary.ID).Fetch(); err != nil {
			if utils.IsNotFoundError(err) {
				return nil
			}
			return fmt.Errorf("Error occurred when retrieving SIP domain IP access control list mapping information %s", err.Error())
		}
	}

	return nil
}

func testAccCheckTwilioSIPDomainIPAccessControlListMappingExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*common.TwilioClient).API

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if _, err := client.Account(rs.Primary.Attributes["account_sid"]).Sip.Domain(rs.Primary.Attributes["domain_sid"]).Auth.Calls.IpAccessControlListMapping(rs.Primary.ID).Fetch(); err != nil {
			return fmt.Errorf("Error occurred when retrieving SIP domain IP access control list mapping information %s", err.Error())
		}

		return nil
	}
}

func testAccTwilioSIPDomainIPAccessControlListMappingImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Not found: %s", name)
		}

		return fmt.Sprintf("/Accounts/%s/SIP/Domains/%s/Auth/Calls/IpAccessControlListMappings/%s", rs.Primary.Attributes["account_sid"], rs.Primary.Attributes["domain_sid"], rs.Primary.Attributes["sid"]), nil
	}
}

func testAccTwilioSIPDomainIPAccessControlListMapping_basic(testData *acceptance.TestData, friendlyName string, domainName string) string {
	return fmt.Sprintf(`
resource "twilio_sip_ip_access_control_list" "ip_access_control_list" {
  account_sid   = "%[1]s"
  friendly_name = "%[2]s"
}

resource "twilio_sip_domain" "domain" {
  account_sid = "%[1]s"
  domain_name = "%[3]s"
}

resource "twilio_sip_domain_ip_access_control_list_mapping" "domain_ip_access_control_list_mapping" {
  account_sid                = twilio_sip_domain.domain.account_sid
  domain_sid                 = twilio_sip_domain.domain.sid
  ip_access_control_list_sid = twilio_sip_ip_access_control_list.ip_access_control_list.sid
}
`, testData.AccountSid, friendlyName, domainName)
}

func testAccTwilioSIPDomainIPAccessControlListMapping_invalidAccountSid() string {
	return `
resource "twilio_sip_domain_ip_access_control_list_mapping" "domain_ip_access_control_list_mapping" {
  account_sid                = "account_sid"
  domain_sid                 = "SDaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  ip_access_control_list_sid = "ALaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}
`
}

func testAccTwilioSIPDomainIPAccessControlListMapping_invalidDomainSid() string {
	return `
resource "twilio_sip_domain_ip_access_control_list_mapping" "domain_ip_access_control_list_mapping" {
  account_sid                = "ACaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  domain_sid                 = "domain_sid"
  ip_access_control_list_sid = "ALaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}
`
}

func testAccTwilioSIPDomainIPAccessControlListMapping_invalidIPAccessControlList() string {
	return `
resource "twilio_sip_domain_ip_access_control_list_mapping" "domain_ip_access_control_list_mapping" {
  account_sid                = "ACaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  domain_sid                 = "SDaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  ip_access_control_list_sid = "ip_access_control_list_sid"
}
`
}
