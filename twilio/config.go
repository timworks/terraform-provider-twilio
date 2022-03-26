package twilio

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/timworks/terraform-provider-twilio/twilio/common"
	"github.com/timworks/twilio-sdk-go/client"
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
	"github.com/timworks/twilio-sdk-go/session"
	"github.com/timworks/twilio-sdk-go/session/credentials"
	"github.com/timworks/twilio-sdk-go/utils"
)

type Config struct {
	AccountSid       string
	AuthToken        string
	APIKey           string
	APISecret        string
	RetryAttempts    int
	BackoffInterval  int
	terraformVersion string
}

func (config *Config) Client() (interface{}, diag.Diagnostics) {

	creds, err := credentials.New(getCredentials(config))
	if err != nil {
		return nil, diag.FromErr(err)
	}

	sess := session.New(creds)
	sdkConfig := &client.Config{
		RetryAttempts:   utils.Int(config.RetryAttempts),
		BackoffInterval: utils.Int(config.BackoffInterval),
	}

	client := &common.TwilioClient{
		AccountSid:       config.AccountSid,
		TerraformVersion: config.terraformVersion,

		Accounts:      accounts.New(sess, sdkConfig),
		API:           api.New(sess, sdkConfig),
		Autopilot:     autopilot.New(sess, sdkConfig),
		Chat:          chat.New(sess, sdkConfig),
		Conversations: conversations.New(sess, sdkConfig),
		Flex:          flex.New(sess, sdkConfig),
		Messaging:     messaging.New(sess, sdkConfig),
		Proxy:         proxy.New(sess, sdkConfig),
		Serverless:    serverless.New(sess, sdkConfig),
		SIPTrunking:   trunking.New(sess, sdkConfig),
		Studio:        studio.New(sess, sdkConfig),
		TaskRouter:    taskrouter.New(sess, sdkConfig),
		Video:         video.New(sess, sdkConfig),
	}
	return client, nil
}

func getCredentials(config *Config) credentials.TwilioCredentials {
	if config.APIKey != "" && config.APISecret != "" {
		return credentials.APIKey{
			Account: config.AccountSid,
			Sid:     config.APIKey,
			Value:   config.APISecret,
		}
	}
	return credentials.Account{
		Sid:       config.AccountSid,
		AuthToken: config.AuthToken,
	}
}
