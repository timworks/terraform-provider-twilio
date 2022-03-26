package sip_trunking

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/timworks/terraform-provider-twilio/twilio/common"
	"github.com/timworks/terraform-provider-twilio/twilio/utils"
	"github.com/timworks/twilio-sdk-go/service/trunking/v1/trunk/ip_access_control_lists"
)

func resourceSIPTrunkingIPAccessControlList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSIPTrunkingIPAccessControlListCreate,
		ReadContext:   resourceSIPTrunkingIPAccessControlListRead,
		DeleteContext: resourceSIPTrunkingIPAccessControlListDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				format := "/Trunks/(.*)/IpAccessControlLists/(.*)"
				regex := regexp.MustCompile(format)
				match := regex.FindStringSubmatch(d.Id())

				if len(match) != 3 {
					return nil, fmt.Errorf("The imported ID (%s) does not match the format (%s)", d.Id(), format)
				}

				d.Set("trunk_sid", match[1])
				d.Set("sid", match[2])
				d.SetId(match[2])
				return []*schema.ResourceData{d}, nil
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trunk_sid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: utils.SIPTrunkSidValidation(),
			},
			"ip_access_control_list_sid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: utils.SIPIPAccessControlListSidValidation(),
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

func resourceSIPTrunkingIPAccessControlListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).SIPTrunking

	createInput := &ip_access_control_lists.CreateIpAccessControlListInput{
		IpAccessControlListSid: d.Get("ip_access_control_list_sid").(string),
	}

	createResult, err := client.Trunk(d.Get("trunk_sid").(string)).IpAccessControlLists.CreateWithContext(ctx, createInput)
	if err != nil {
		return diag.Errorf("Failed to create IP access control list: %s", err.Error())
	}

	d.SetId(createResult.Sid)
	return resourceSIPTrunkingIPAccessControlListRead(ctx, d, meta)
}

func resourceSIPTrunkingIPAccessControlListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).SIPTrunking

	getResponse, err := client.Trunk(d.Get("trunk_sid").(string)).IpAccessControlList(d.Id()).FetchWithContext(ctx)
	if err != nil {
		if utils.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Failed to read IP access control list: %s", err.Error())
	}

	d.Set("sid", getResponse.Sid)
	d.Set("trunk_sid", getResponse.TrunkSid)
	d.Set("ip_access_control_list_sid", getResponse.Sid) // The IpAccessControlListSid is stored as the resource sid
	d.Set("account_sid", getResponse.AccountSid)
	d.Set("friendly_name", getResponse.FriendlyName)
	d.Set("date_created", getResponse.DateCreated.Format(time.RFC3339))

	if getResponse.DateUpdated != nil {
		d.Set("date_updated", getResponse.DateUpdated.Format(time.RFC3339))
	}
	d.Set("url", getResponse.URL)

	return nil
}

func resourceSIPTrunkingIPAccessControlListDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).SIPTrunking

	if err := client.Trunk(d.Get("trunk_sid").(string)).IpAccessControlList(d.Id()).DeleteWithContext(ctx); err != nil {
		return diag.Errorf("Failed to delete IP access control list: %s", err.Error())
	}
	d.SetId("")
	return nil
}
