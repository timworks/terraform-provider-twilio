package autopilot

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
	"github.com/timworks/twilio-sdk-go/service/autopilot/v1/assistant/task/fields"
)

func resourceAutopilotTaskField() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAutopilotTaskFieldCreate,
		ReadContext:   resourceAutopilotTaskFieldRead,
		DeleteContext: resourceAutopilotTaskFieldDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				format := "/Assistants/(.*)/Tasks/(.*)/Fields/(.*)"
				regex := regexp.MustCompile(format)
				match := regex.FindStringSubmatch(d.Id())

				if len(match) != 4 {
					return nil, fmt.Errorf("The imported ID (%s) does not match the format (%s)", d.Id(), format)
				}

				d.Set("assistant_sid", match[1])
				d.Set("task_sid", match[2])
				d.Set("sid", match[3])
				d.SetId(match[3])
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
			"assistant_sid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: utils.AutopilotAssistantSidValidation(),
			},
			"task_sid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: utils.AutopilotTaskSidValidation(),
			},
			"unique_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"field_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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

func resourceAutopilotTaskFieldCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Autopilot

	createInput := &fields.CreateFieldInput{
		UniqueName: d.Get("unique_name").(string),
		FieldType:  d.Get("field_type").(string),
	}

	createResult, err := client.Assistant(d.Get("assistant_sid").(string)).Task(d.Get("task_sid").(string)).Fields.CreateWithContext(ctx, createInput)
	if err != nil {
		return diag.Errorf("Failed to create autopilot task field: %s", err.Error())
	}

	d.SetId(createResult.Sid)
	return resourceAutopilotTaskFieldRead(ctx, d, meta)
}

func resourceAutopilotTaskFieldRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Autopilot

	getResponse, err := client.Assistant(d.Get("assistant_sid").(string)).Task(d.Get("task_sid").(string)).Field(d.Id()).FetchWithContext(ctx)
	if err != nil {
		if utils.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Failed to read autopilot task field: %s", err.Error())
	}

	d.Set("sid", getResponse.Sid)
	d.Set("account_sid", getResponse.AccountSid)
	d.Set("assistant_sid", getResponse.AssistantSid)
	d.Set("task_sid", getResponse.TaskSid)
	d.Set("unique_name", getResponse.UniqueName)
	d.Set("field_type", getResponse.FieldType)
	d.Set("date_created", getResponse.DateCreated.Format(time.RFC3339))

	if getResponse.DateUpdated != nil {
		d.Set("date_updated", getResponse.DateUpdated.Format(time.RFC3339))
	}

	d.Set("url", getResponse.URL)
	return nil
}

func resourceAutopilotTaskFieldDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Autopilot

	if err := client.Assistant(d.Get("assistant_sid").(string)).Task(d.Get("task_sid").(string)).Field(d.Id()).DeleteWithContext(ctx); err != nil {
		return diag.Errorf("Failed to delete autopilot task field: %s", err.Error())
	}
	d.SetId("")
	return nil
}
