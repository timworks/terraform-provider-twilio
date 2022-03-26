package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/acceptance"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/studio/tests/helper"
)

func TestAccDataSourceTwilioStudioFlowDefinition_basic(t *testing.T) {
	stateDataSourceName := "data.twilio_studio_flow_definition.definition"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTwilioStudioFlowDefinition_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(stateDataSourceName, "json", `{"description":"Bot flow for creating a Flex webchat task","flags":{"allow_concurrent_calls":true},"initial_state":"Trigger","states":[{"name":"Trigger","properties":{"offset":{"x":200,"y":0}},"transitions":[{"event":"incomingCall"},{"event":"incomingMessage","next":"SendToAutopilot"},{"event":"incomingRequest"}],"type":"trigger"},{"name":"SendMessageToAgent","properties":{"attributes":"{\"channelSid\":\"{{trigger.message.ChannelSid}}\",\"channelType\":\"{{trigger.message.ChannelAttributes.channel_type}}\",\"name\":\"{{trigger.message.ChannelAttributes.from}}\"}","channel":"TCaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","offset":{"x":270,"y":540},"workflow":"WWaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"transitions":[{"event":"callComplete"},{"event":"callFailure"},{"event":"failedToEnqueue"}],"type":"send-to-flex"},{"name":"SendToAutopilot","properties":{"autopilot_assistant_sid":"UAaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","body":"{{trigger.message.Body}}","from":"{{flow.channel.address}}","offset":{"x":200,"y":240},"timeout":14400},"transitions":[{"event":"failure","next":"SendMessageToAgent"},{"event":"sessionEnded"},{"event":"timeout","next":"SendMessageToAgent"}],"type":"send-to-auto-pilot"}]}`),
					helper.ValidateFlowDefinition(stateDataSourceName),
				),
			},
		},
	})
}

func testAccDataSourceTwilioStudioFlowDefinition_basic() string {
	return `
data "twilio_studio_flow_widget_send_to_autopilot" "send_to_autopilot" {
	name  = "SendToAutopilot"
	
	offset {
		x = 200
		y = 240
	}
	
	transitions {
		failure = data.twilio_studio_flow_widget_send_to_flex.send_to_flex.name
		timeout = data.twilio_studio_flow_widget_send_to_flex.send_to_flex.name
	}
	
	autopilot_assistant_sid = "UAaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	from                    = "{{flow.channel.address}}"
	body                    = "{{trigger.message.Body}}"
	timeout                 = 14400
}
	
data "twilio_studio_flow_widget_send_to_flex" "send_to_flex" {
	name = "SendMessageToAgent"
	
	workflow_sid = "WWaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	channel_sid = "TCaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	attributes = jsonencode({
		"name" : "{{trigger.message.ChannelAttributes.from}}",
		"channelType" : "{{trigger.message.ChannelAttributes.channel_type}}",
		"channelSid" : "{{trigger.message.ChannelSid}}"
	})
	
	offset {
		x = 270
		y = 540
	}
}
	
data "twilio_studio_flow_widget_trigger" "trigger" {
	name = "Trigger"
	
	transitions {
		incoming_message = data.twilio_studio_flow_widget_send_to_autopilot.send_to_autopilot.name
	}
	
	offset {
		x = 200
		y = 0
	}
}
	
data "twilio_studio_flow_definition" "definition" {
	description   = "Bot flow for creating a Flex webchat task"
	initial_state = data.twilio_studio_flow_widget_trigger.trigger.name
	
	flags {
		allow_concurrent_calls = true
	}
	
	states {
		json = data.twilio_studio_flow_widget_trigger.trigger.json
	}

	states {
		json = data.twilio_studio_flow_widget_send_to_flex.send_to_flex.json
	}

	states {
		json = data.twilio_studio_flow_widget_send_to_autopilot.send_to_autopilot.json
	}
}
`
}
