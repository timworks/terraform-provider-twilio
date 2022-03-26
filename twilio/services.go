package twilio

import (
	"github.com/timworks/terraform-provider-twilio/twilio/common"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/account"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/autopilot"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/chat"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/conversations"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/credentials"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/flex"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/iam"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/messaging"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/phone_number"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/proxy"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/serverless"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/sip"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/sip_trunking"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/studio"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/taskrouter"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/twiml"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/video"
	"github.com/timworks/terraform-provider-twilio/twilio/internal/services/voice"
)

func SupportedServices() []common.ServiceRegistration {
	return []common.ServiceRegistration{
		account.Registration{},
		autopilot.Registration{},
		chat.Registration{},
		credentials.Registration{},
		conversations.Registration{},
		flex.Registration{},
		iam.Registration{},
		messaging.Registration{},
		phone_number.Registration{},
		proxy.Registration{},
		serverless.Registration{},
		studio.Registration{},
		sip.Registration{},
		sip_trunking.Registration{},
		taskrouter.Registration{},
		twiml.Registration{},
		video.Registration{},
		voice.Registration{},
	}
}
