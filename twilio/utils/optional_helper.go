package utils

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	sdkUtils "github.com/timworks/twilio-sdk-go/utils"
)

func OptionalString(d *schema.ResourceData, key string) *string {
	if v, ok := d.GetOk(key); ok {
		return sdkUtils.String(v.(string))
	}
	return nil
}

func OptionalStringWithEmptyStringOnChange(d *schema.ResourceData, key string) *string {
	if v, ok := d.GetOk(key); ok {
		return sdkUtils.String(v.(string))
	}
	return defaultToEmptyStringIfChanged(d, key)
}

func OptionalJSONString(d *schema.ResourceData, key string) *string {
	if v, ok := d.GetOk(key); ok {
		// error not handled as it is assumed stringIsJSON validation is applied to the resource
		normalizedJSON, _ := structure.NormalizeJsonString(v.(string))
		return sdkUtils.String(normalizedJSON)
	}
	return nil
}

func OptionalJSONStringWithEmptyStringOnChange(d *schema.ResourceData, key string) *string {
	if v, ok := d.GetOk(key); ok {
		// error not handled as it is assumed stringIsJSON validation is applied to the resource
		normalizedJSON, _ := structure.NormalizeJsonString(v.(string))
		return sdkUtils.String(normalizedJSON)
	}
	return defaultToEmptyStringIfChanged(d, key)
}

func OptionalSeperatedString(d *schema.ResourceData, key string, separator string) *string {
	if v, ok := d.GetOk(key); ok {
		return sdkUtils.String(ConvertSliceToSeperatedString(v.([]interface{}), separator))
	}
	return nil
}

func OptionalSeperatedStringWithEmptyStringOnChange(d *schema.ResourceData, key string, separator string) *string {
	if v, ok := d.GetOk(key); ok {
		return sdkUtils.String(ConvertSliceToSeperatedString(v.([]interface{}), separator))
	}
	return defaultToEmptyStringIfChanged(d, key)
}

func OptionalStringSlice(d *schema.ResourceData, key string) *[]string {
	if v, ok := d.GetOk(key); ok {
		stringSlice := ConvertToStringSlice(v.([]interface{}))
		return &stringSlice
	}
	return nil
}

func OptionalInt(d *schema.ResourceData, key string) *int {
	if v, ok := d.GetOk(key); ok {
		return sdkUtils.Int(v.(int))
	}
	return nil
}

func OptionalIntWith0OnChange(d *schema.ResourceData, key string) *int {
	if v, ok := d.GetOk(key); ok {
		return sdkUtils.Int(v.(int))
	}
	if ok := d.HasChange(key); ok {
		return sdkUtils.Int(0)
	}
	return nil
}

func OptionalBool(d *schema.ResourceData, key string) *bool {
	if v, ok := d.GetOkExists(key); ok {
		return sdkUtils.Bool(v.(bool))
	}
	return nil
}

// defaultToEmptyStringIfChanged caters for the scenario where Terraform previously had a value
// but it has since been removed for the terraform configuration, so setting this to an empty
// string to force the value to be unset in Twilio
func defaultToEmptyStringIfChanged(d *schema.ResourceData, key string) *string {
	if ok := d.HasChange(key); ok {
		return sdkUtils.String("")
	}
	return nil
}
