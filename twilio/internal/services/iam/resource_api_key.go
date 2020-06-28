package iam

import (
	"fmt"
	"time"

	"github.com/RJPearson94/terraform-provider-twilio/twilio/common"
	"github.com/RJPearson94/terraform-provider-twilio/twilio/utils"
	"github.com/RJPearson94/twilio-sdk-go/service/api/v2010/account/key"
	"github.com/RJPearson94/twilio-sdk-go/service/api/v2010/account/keys"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIamApiKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiKeyCreate,
		Read:   resourceApiKeyRead,
		Update: resourceApiKeyUpdate,
		Delete: resourceApiKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_sid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApiKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*common.TwilioClient).API

	createInput := &keys.CreateKeyInput{
		FriendlyName: utils.OptionalString(d, "friendly_name"),
	}

	createResult, err := client.Account(d.Get("account_sid").(string)).Keys.Create(createInput)
	if err != nil {
		return fmt.Errorf("Failed to create account api key: %s", err.Error())
	}

	d.SetId(createResult.Sid)
	d.Set("secret", createResult.Secret)
	return resourceApiKeyRead(d, meta)
}

func resourceApiKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*common.TwilioClient).API

	getResponse, err := client.Account(d.Get("account_sid").(string)).Key(d.Id()).Get()
	if err != nil {
		if utils.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Failed to read account api key: %s", err)
	}

	d.Set("sid", getResponse.Sid)
	d.Set("account_sid", d.Get("account_sid").(string))
	d.Set("friendly_name", getResponse.FriendlyName)
	d.Set("date_created", getResponse.DateCreated.Time.Format(time.RFC3339))

	if getResponse.DateUpdated != nil {
		d.Set("date_updated", getResponse.DateUpdated.Time.Format(time.RFC3339))
	}

	return nil
}

func resourceApiKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*common.TwilioClient).API

	updateInput := &key.UpdateKeyInput{
		FriendlyName: utils.OptionalString(d, "friendly_name"),
	}

	updateResp, err := client.Account(d.Get("account_sid").(string)).Key(d.Id()).Update(updateInput)
	if err != nil {
		return fmt.Errorf("Failed to update account api key: %s", err.Error())
	}

	d.SetId(updateResp.Sid)
	return resourceApiKeyRead(d, meta)
}

func resourceApiKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*common.TwilioClient).API

	if err := client.Account(d.Get("account_sid").(string)).Key(d.Id()).Delete(); err != nil {
		return fmt.Errorf("Failed to delete account api key: %s", err.Error())
	}

	d.SetId("")
	return nil
}
