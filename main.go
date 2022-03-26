package main

import (
	"github.com/timworks/terraform-provider-twilio/twilio"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: twilio.Provider,
	})
}
