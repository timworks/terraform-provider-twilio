package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/acceptance"
)

var modelBuildDataSourceName = "twilio_autopilot_model_build"

func TestAccDataSourceTwilioAutopilotModelBuild_sid(t *testing.T) {
	stateDataSourceName := fmt.Sprintf("data.%s.model_build", modelBuildDataSourceName)
	uniqueName := acctest.RandString(10)
	modelBuildUniqueNamePrefix := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTwilioAutopilotModelBuild_sid(uniqueName, modelBuildUniqueNamePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(stateDataSourceName, "unique_name", regexp.MustCompile(fmt.Sprintf("^%s", modelBuildUniqueNamePrefix))),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "id"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "assistant_sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "build_duration"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "status"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "error_code"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "url"),
				),
			},
		},
	})
}

func TestAccDataSourceTwilioAutopilotModelBuild_uniqueName(t *testing.T) {
	stateDataSourceName := fmt.Sprintf("data.%s.model_build", modelBuildDataSourceName)
	uniqueName := acctest.RandString(10)
	modelBuildUniqueNamePrefix := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTwilioAutopilotModelBuild_uniqueName(uniqueName, modelBuildUniqueNamePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(stateDataSourceName, "unique_name", regexp.MustCompile(fmt.Sprintf("^%s", modelBuildUniqueNamePrefix))),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "id"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "assistant_sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "build_duration"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "status"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "error_code"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "url"),
				),
			},
		},
	})
}

func TestAccDataSourceTwilioAutopilotModelBuild_invalidAssistantSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceTwilioAutopilotModelBuild_invalidAssistantSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of assistant_sid to match regular expression "\^UA\[0-9a-fA-F\]\{32\}\$", got assistant_sid`),
			},
		},
	})
}

func TestAccDataSourceTwilioAutopilotModelBuild_invalidSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceTwilioAutopilotModelBuild_invalidSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of sid to match regular expression "\^UG\[0-9a-fA-F\]\{32\}\$", got sid`),
			},
		},
	})
}

func testAccDataSourceTwilioAutopilotModelBuild_sid(uniqueName string, modelBuildUniqueNamePrefix string) string {
	return fmt.Sprintf(`
resource "twilio_autopilot_assistant" "assistant" {
  unique_name = "%[1]s"
}

resource "twilio_autopilot_task" "task" {
  assistant_sid = twilio_autopilot_assistant.assistant.sid
  unique_name   = "%[1]s"
}

resource "twilio_autopilot_task_sample" "task_sample" {
  assistant_sid = twilio_autopilot_assistant.assistant.sid
  task_sid      = twilio_autopilot_task.task.sid
  language      = "en-US"
  tagged_text   = "test"
}

resource "twilio_autopilot_model_build" "model_build" {
  assistant_sid      = twilio_autopilot_assistant.assistant.sid
  unique_name_prefix = "%[2]s"

  triggers = {
    redeployment = sha1(join(",", tolist([
      twilio_autopilot_task_sample.task_sample.sid,
      twilio_autopilot_task_sample.task_sample.language,
      twilio_autopilot_task_sample.task_sample.tagged_text,
    ])))
  }

  lifecycle {
    create_before_destroy = true
  }

  polling {
    enabled = true
  }
}

data "twilio_autopilot_model_build" "model_build" {
  assistant_sid = twilio_autopilot_model_build.model_build.assistant_sid
  sid           = twilio_autopilot_model_build.model_build.sid
}
`, uniqueName, modelBuildUniqueNamePrefix)
}

func testAccDataSourceTwilioAutopilotModelBuild_uniqueName(uniqueName string, modelBuildUniqueNamePrefix string) string {
	return fmt.Sprintf(`
resource "twilio_autopilot_assistant" "assistant" {
  unique_name = "%[1]s"
}

resource "twilio_autopilot_task" "task" {
  assistant_sid = twilio_autopilot_assistant.assistant.sid
  unique_name   = "%[1]s"
}

resource "twilio_autopilot_task_sample" "task_sample" {
  assistant_sid = twilio_autopilot_assistant.assistant.sid
  task_sid      = twilio_autopilot_task.task.sid
  language      = "en-US"
  tagged_text   = "test"
}

resource "twilio_autopilot_model_build" "model_build" {
  assistant_sid      = twilio_autopilot_assistant.assistant.sid
  unique_name_prefix = "%[2]s"

  triggers = {
    redeployment = sha1(join(",", tolist([
      twilio_autopilot_task_sample.task_sample.sid,
      twilio_autopilot_task_sample.task_sample.language,
      twilio_autopilot_task_sample.task_sample.tagged_text,
    ])))
  }

  lifecycle {
    create_before_destroy = true
  }

  polling {
    enabled = true
  }
}

data "twilio_autopilot_model_build" "model_build" {
  assistant_sid = twilio_autopilot_model_build.model_build.assistant_sid
  unique_name   = twilio_autopilot_model_build.model_build.unique_name
}
`, uniqueName, modelBuildUniqueNamePrefix)
}

func testAccDataSourceTwilioAutopilotModelBuild_invalidAssistantSid() string {
	return `
data "twilio_autopilot_model_build" "model_build" {
  assistant_sid = "assistant_sid"
  sid           = "UGaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}
`
}

func testAccDataSourceTwilioAutopilotModelBuild_invalidSid() string {
	return `
data "twilio_autopilot_model_build" "model_build" {
  assistant_sid = "UAaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  sid           = "sid"
}
`
}
