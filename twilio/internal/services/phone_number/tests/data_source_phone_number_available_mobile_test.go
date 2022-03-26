package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/acceptance"
)

var phoneNumberAvailableMobileDataSourceName = "twilio_phone_number_available_mobile_numbers"

func TestAccDataSourceTwilioPhoneNumberAvailableMobile_complete(t *testing.T) {
	stateDataSourceName := fmt.Sprintf("data.%s.available_mobile_numbers", phoneNumberAvailableMobileDataSourceName)
	testData := acceptance.TestAccData

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTwilioPhoneNumberAvailableMobile_complete(testData),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(stateDataSourceName, "id"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "account_sid"),
					resource.TestCheckResourceAttr(stateDataSourceName, "iso_country", "GB"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "available_phone_numbers.#"),
				),
			},
		},
	})
}

func testAccTwilioPhoneNumberAvailableMobile_complete(testData *acceptance.TestData) string {
	return fmt.Sprintf(`
data "twilio_phone_number_available_mobile_numbers" "available_mobile_numbers" {
  account_sid = "%s"
  iso_country = "GB"
}
`, testData.AccountSid)
}
