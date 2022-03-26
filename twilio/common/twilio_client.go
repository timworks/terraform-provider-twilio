package common

import (
	accounts "github.com/timworks/twilio-sdk-go/service/accounts/v1"
	api "github.com/timworks/twilio-sdk-go/service/api/v2010"
	autopilot "github.com/timworks/twilio-sdk-go/service/autopilot/v1"
	chat "github.com/timworks/twilio-sdk-go/service/chat/v2"
	conversations "github.com/timworks/twilio-sdk-go/service/conversations/v1"
	flex "github.com/timworks/twilio-sdk-go/service/flex/v1"
	messaging "github.com/timworks/twilio-sdk-go/service/messaging/v1"
	proxy "github.com/timworks/twilio-sdk-go/service/proxy/v1"
	serverless "github.com/timworks/twilio-sdk-go/service/serverless/v1"
	studio "github.com/timworks/twilio-sdk-go/service/studio/v2"
	taskrouter "github.com/timworks/twilio-sdk-go/service/taskrouter/v1"
	trunking "github.com/timworks/twilio-sdk-go/service/trunking/v1"
	video "github.com/timworks/twilio-sdk-go/service/video/v1"
)

type TwilioClient struct {
	AccountSid       string
	TerraformVersion string

	Accounts      *accounts.Accounts
	API           *api.V2010
	Autopilot     *autopilot.Autopilot
	Chat          *chat.Chat
	Conversations *conversations.Conversations
	Flex          *flex.Flex
	Messaging     *messaging.Messaging
	Proxy         *proxy.Proxy
	Serverless    *serverless.Serverless
	SIPTrunking   *trunking.Trunking
	Studio        *studio.Studio
	TaskRouter    *taskrouter.TaskRouter
	Video         *video.Video
}
