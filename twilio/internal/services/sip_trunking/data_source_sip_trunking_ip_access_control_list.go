package sip_trunking

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/timworks/terraform-provider-twilio/twilio/common"
	"github.com/timworks/terraform-provider-twilio/twilio/utils"
)

func dataSourceSIPTrunkingIPAccessControlList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSIPTrunkingIPAccessControlListRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"sid": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.SIPIPAccessControlListSidValidation(),
			},
			"trunk_sid": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.SIPTrunkSidValidation(),
			},
			"account_sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"friendly_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSIPTrunkingIPAccessControlListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).SIPTrunking

	trunkSid := d.Get("trunk_sid").(string)
	sid := d.Get("sid").(string)
	getResponse, err := client.Trunk(trunkSid).IpAccessControlList(sid).FetchWithContext(ctx)
	if err != nil {
		if utils.IsNotFoundError(err) {
			return diag.Errorf("SIP trunk IP access control list with sid (%s) was not found for SIP trunk with sid (%s)", sid, trunkSid)
		}
		return diag.Errorf("Failed to read SIP trunk IP access control list: %s", err.Error())
	}

	d.SetId(getResponse.Sid)
	d.Set("sid", getResponse.Sid)
	d.Set("account_sid", getResponse.AccountSid)
	d.Set("trunk_sid", getResponse.TrunkSid)
	d.Set("friendly_name", getResponse.FriendlyName)
	d.Set("date_created", getResponse.DateCreated.Format(time.RFC3339))

	if getResponse.DateUpdated != nil {
		d.Set("date_updated", getResponse.DateUpdated.Format(time.RFC3339))
	}
	d.Set("url", getResponse.URL)

	return nil
}
