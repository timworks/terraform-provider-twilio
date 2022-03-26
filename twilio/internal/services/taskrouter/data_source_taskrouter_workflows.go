package taskrouter

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/timworks/terraform-provider-twilio/twilio/common"
	"github.com/timworks/terraform-provider-twilio/twilio/utils"
	"github.com/timworks/twilio-sdk-go/service/taskrouter/v1/workspace/workflows"
)

func dataSourceTaskRouterWorkflows() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaskRouterWorkflowsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"workspace_sid": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.TaskRouterWorkspaceSidValidation(),
			},
			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workflows": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"friendly_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fallback_assignment_callback_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assignment_callback_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_reservation_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"document_content_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"configuration": {
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
				},
			},
		},
	}
}

func dataSourceTaskRouterWorkflowsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	twilioClient := meta.(*common.TwilioClient)
	client := twilioClient.TaskRouter

	options := &workflows.WorkflowsPageOptions{
		FriendlyName: utils.OptionalString(d, "friendly_name"),
	}

	workspaceSid := d.Get("workspace_sid").(string)
	paginator := client.Workspace(workspaceSid).Workflows.NewWorkflowsPaginatorWithOptions(options)
	for paginator.NextWithContext(ctx) {
	}

	err := paginator.Error()
	if err != nil {
		if utils.IsNotFoundError(err) {
			return diag.Errorf("No workflows were found for taskrouter workspace with sid (%s)", workspaceSid)
		}
		return diag.Errorf("Failed to list workflows: %s", err.Error())
	}

	d.SetId(workspaceSid)
	d.Set("workspace_sid", workspaceSid)
	d.Set("account_sid", twilioClient.AccountSid)

	workflows := make([]interface{}, 0)

	for _, workflow := range paginator.Workflows {
		workflowsMap := make(map[string]interface{})

		workflowsMap["sid"] = workflow.Sid
		workflowsMap["friendly_name"] = workflow.FriendlyName
		workflowsMap["fallback_assignment_callback_url"] = workflow.FallbackAssignmentCallbackURL
		workflowsMap["assignment_callback_url"] = workflow.AssignmentCallbackURL
		workflowsMap["task_reservation_timeout"] = workflow.TaskReservationTimeout
		workflowsMap["document_content_type"] = workflow.DocumentContentType
		workflowsMap["configuration"] = workflow.Configuration
		workflowsMap["date_created"] = workflow.DateCreated.Format(time.RFC3339)

		if workflow.DateUpdated != nil {
			workflowsMap["date_updated"] = workflow.DateUpdated.Format(time.RFC3339)
		}

		workflowsMap["url"] = workflow.URL

		workflows = append(workflows, workflowsMap)
	}

	d.Set("workflows", &workflows)

	return nil
}
