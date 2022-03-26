package conversations

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/timworks/terraform-provider-twilio/twilio/common"
	"github.com/timworks/terraform-provider-twilio/twilio/utils"
	"github.com/timworks/twilio-sdk-go/service/conversations/v1/credential"
	"github.com/timworks/twilio-sdk-go/service/conversations/v1/credentials"
)

func resourceConversationsPushCredentialFCM() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConversationsPushCredentialFCMCreate,
		ReadContext:   resourceConversationsPushCredentialFCMRead,
		UpdateContext: resourceConversationsPushCredentialFCMUpdate,
		DeleteContext: resourceConversationsPushCredentialFCMDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				format := "/Credentials/(.*)"
				regex := regexp.MustCompile(format)
				match := regex.FindStringSubmatch(d.Id())

				if len(match) != 2 {
					return nil, fmt.Errorf("The imported ID (%s) does not match the format (%s)", d.Id(), format)
				}

				d.Set("sid", match[1])
				d.SetId(match[1])
				return []*schema.ResourceData{d}, nil
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
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
			"friendly_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"secret": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"type": {
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

func resourceConversationsPushCredentialFCMCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Conversations

	createInput := &credentials.CreateCredentialInput{
		Type:         "fcm",
		Secret:       utils.OptionalString(d, "secret"),
		FriendlyName: utils.OptionalString(d, "friendly_name"),
	}

	createResult, err := client.Credentials.CreateWithContext(ctx, createInput)
	if err != nil {
		return diag.Errorf("Failed to create conversations credential: %s", err.Error())
	}

	d.SetId(createResult.Sid)
	return resourceConversationsPushCredentialFCMRead(ctx, d, meta)
}

func resourceConversationsPushCredentialFCMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Conversations

	getResponse, err := client.Credential(d.Id()).FetchWithContext(ctx)
	if err != nil {
		if utils.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Failed to read conversations credential: %s", err.Error())
	}

	d.Set("sid", getResponse.Sid)
	d.Set("account_sid", getResponse.AccountSid)
	d.Set("friendly_name", getResponse.FriendlyName)
	d.Set("type", getResponse.Type)
	d.Set("date_created", getResponse.DateCreated.Format(time.RFC3339))

	if getResponse.DateUpdated != nil {
		d.Set("date_updated", getResponse.DateUpdated.Format(time.RFC3339))
	}

	d.Set("url", getResponse.URL)

	return nil
}

func resourceConversationsPushCredentialFCMUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Conversations

	updateInput := &credential.UpdateCredentialInput{
		Secret:       utils.OptionalString(d, "secret"),
		FriendlyName: utils.OptionalString(d, "friendly_name"),
	}

	updateResp, err := client.Credential(d.Id()).UpdateWithContext(ctx, updateInput)
	if err != nil {
		return diag.Errorf("Failed to update conversations credential: %s", err.Error())
	}

	d.SetId(updateResp.Sid)
	return resourceConversationsPushCredentialFCMRead(ctx, d, meta)
}

func resourceConversationsPushCredentialFCMDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Conversations

	if err := client.Credential(d.Id()).DeleteWithContext(ctx); err != nil {
		return diag.Errorf("Failed to delete conversations credential: %s", err.Error())
	}
	d.SetId("")
	return nil
}
