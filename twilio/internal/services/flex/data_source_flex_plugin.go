package flex

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/timworks/terraform-provider-twilio/twilio/common"
	"github.com/timworks/terraform-provider-twilio/twilio/utils"
	"github.com/timworks/twilio-sdk-go/service/flex/v1/plugin/versions"
	sdkUtils "github.com/timworks/twilio-sdk-go/utils"
)

func dataSourceFlexPlugin() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFlexPluginRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"sid": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: utils.FlexPluginSidValidation(),
				ExactlyOneOf: []string{"sid", "unique_name"},
			},
			"unique_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"sid", "unique_name"},
			},
			"account_sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"friendly_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"archived": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"changelog": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plugin_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"version_archived": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"latest_version_sid": {
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

func dataSourceFlexPluginRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Flex

	var identifier string

	if v, ok := d.GetOk("sid"); ok {
		identifier = v.(string)
	} else {
		identifier = d.Get("unique_name").(string)
	}

	getResponse, err := client.Plugin(identifier).FetchWithContext(ctx)
	if err != nil {
		if utils.IsNotFoundError(err) {
			return diag.Errorf("Flex plugin with sid/ unique name (%s) was not found", identifier)
		}
		return diag.Errorf("Failed to read flex plugin: %s", err.Error())
	}

	versionsPaginator := client.Plugin(getResponse.Sid).Versions.NewVersionsPaginatorWithOptions(&versions.VersionsPageOptions{
		PageSize: sdkUtils.Int(5),
	})
	// The twilio api return the latest version as the first element in the array.
	// So there is no need to loop to retrieve all records
	versionsPaginator.Next()

	if versionsPaginator.Error() != nil {
		return diag.Errorf("Failed to read flex plugin versions: %s", versionsPaginator.Error().Error())
	}

	d.SetId(getResponse.Sid)
	d.Set("sid", getResponse.Sid)
	d.Set("account_sid", getResponse.AccountSid)
	d.Set("archived", getResponse.Archived)
	d.Set("description", getResponse.Description)
	d.Set("friendly_name", getResponse.FriendlyName)
	d.Set("unique_name", getResponse.UniqueName)
	d.Set("date_created", getResponse.DateCreated.Format(time.RFC3339))

	if getResponse.DateUpdated != nil {
		d.Set("date_updated", getResponse.DateUpdated.Format(time.RFC3339))
	}

	d.Set("url", getResponse.URL)

	if len(versionsPaginator.Versions) > 0 {
		latestVersion := versionsPaginator.Versions[0]

		d.Set("latest_version_sid", latestVersion.Sid)
		d.Set("changelog", latestVersion.Changelog)
		d.Set("version", latestVersion.Version)
		d.Set("plugin_url", latestVersion.PluginURL)
		d.Set("private", latestVersion.Private)
		d.Set("version_archived", latestVersion.Archived)
	} else {
		log.Printf("[INFO] No flex plugin versions found for plugin (%s)", getResponse.Sid)
	}

	return nil
}
