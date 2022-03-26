package video

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/timworks/terraform-provider-twilio/twilio/common"
	"github.com/timworks/terraform-provider-twilio/twilio/utils"
)

func dataSourceVideoRecordingSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVideoRecordingSettingsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"account_sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_credentials_sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_s3_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_storage_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"encryption_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"encryption_key_sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"friendly_name": {
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

func dataSourceVideoRecordingSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Video

	getResponse, err := client.RecordingSettings().FetchWithContext(ctx)
	if err != nil {
		if utils.IsNotFoundError(err) {
			return diag.Errorf("Video recording settings was not found")
		}
		return diag.Errorf("Failed to read recording settings: %s", err.Error())
	}

	d.SetId(getResponse.AccountSid)
	d.Set("account_sid", getResponse.AccountSid)
	d.Set("aws_credentials_sid", getResponse.AWSCredentialSid)
	d.Set("aws_s3_url", getResponse.AWSS3URL)
	d.Set("aws_storage_enabled", getResponse.AWSStorageEnabled)
	d.Set("encryption_enabled", getResponse.EncryptionEnabled)
	d.Set("encryption_key_sid", getResponse.EncryptionKeySid)
	d.Set("friendly_name", getResponse.FriendlyName)
	d.Set("url", getResponse.URL)

	return nil
}
