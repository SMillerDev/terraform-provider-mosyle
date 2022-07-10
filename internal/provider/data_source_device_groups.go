package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDeviceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeviceGroupsRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeMap,
				Description: "Filters to limit API data",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"groups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Device group data from Mosyle",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":             &schema.Schema{Type: schema.TypeString, Computed: true},
						"name":           &schema.Schema{Type: schema.TypeString, Computed: true},
						"device_numbers": &schema.Schema{Type: schema.TypeInt, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceDeviceGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	filter := d.Get("filter").(map[string]interface{})
	options := make(map[string]interface{})
	for key, value := range filter {
		options[key] = value
	}

	req_data := ListPostBody{
		Operation: "list_devicegroup",
		Options:   options,
	}
	req_body, err := json.Marshal(req_data)
	if err != nil {
		return diag.FromErr(err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/devicegroups", c.HostURL), strings.NewReader(string(req_body)))
	if err != nil {
		return diag.FromErr(err)
	}

	b, err := c.doBaseRequest(req)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
			Detail:   fmt.Sprint(string(req_body)),
		})
	}

	response := DeviceGroupListResponse{}
	err = json.Unmarshal(b, &response)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to Decode",
			Detail:   err.Error(),
		})
	}

	if response.Status != "OK" {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Non succesful API call",
			Detail:   string(b),
		})
	}

	if err := d.Set("groups", flattenDeviceGroups(response)); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to transfer data",
			Detail:   err.Error(),
		})
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenDeviceGroups(response DeviceGroupListResponse) []interface{} {
	if len(response.Response.DeviceGroups) < 1 {
		return make([]interface{}, 0)
	}

	ois := make([]interface{}, len(response.Response.DeviceGroups), len(response.Response.DeviceGroups))

	for i, device := range response.Response.DeviceGroups {
		oi := make(map[string]interface{})

		for key, val := range device {
			if key == "device_numbers" {
				ival, err := strconv.ParseInt(val.(string), 10, 64)
				if err == nil {
					oi[strings.ToLower(key)] = ival
					continue
				}
			}

			oi[strings.ToLower(key)] = val
		}

		ois[i] = oi
	}

	return ois
}
