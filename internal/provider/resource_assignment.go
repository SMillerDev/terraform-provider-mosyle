package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssignmentCreate,
		ReadContext:   resourceAssignmentRead,
		UpdateContext: resourceAssignmentUpdate,
		DeleteContext: resourceAssignmentDelete,
		Description:   "Assignment data",
		Schema: map[string]*schema.Schema{
			"os":            &schema.Schema{Type: schema.TypeString, Required: true, Description: "Assignment device os"},
			"user_id":       &schema.Schema{Type: schema.TypeString, Required: true, Description: "Assignment user identifier"},
			"device_serial": &schema.Schema{Type: schema.TypeString, Required: true, Description: "Assigned device"},
			"device_udid":   &schema.Schema{Type: schema.TypeString, Computed: true, Description: "Device UDID"},
		},
	}
}

func resourceAssignmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Get("user_id").(string)
	serial := d.Get("device_serial").(string)

	assign_data := map[string]string{"iduser": id, "serialnumber": serial}
	data := map[string]interface{}{"operation": "assign_device_user", "assign": [...]map[string]string{assign_data}}
	req_body, err := json.Marshal(data)
	if err != nil {
		return diag.FromErr(err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/devices", c.HostURL), strings.NewReader(string(req_body)))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.doRequest(req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serial)

	resourceAssignmentRead(ctx, d, m)

	return diags
}

func resourceAssignmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	serial := d.Get("device_serial").(string)
	os := d.Get("os").(string)

	options := map[string]interface{}{"os": os, "serial_numbers": [...]string{serial}}
	data := map[string]interface{}{"operation": "list", "options": options}
	req_body, err := json.Marshal(data)
	if err != nil {
		return diag.FromErr(err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/devices", c.HostURL), strings.NewReader(string(req_body)))
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := c.doRequest(req)
	if err != nil {
		return diag.FromErr(err)
	}

	Assignment := flattenAssignment(response)
	if Assignment == nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "No such Assignment",
			Detail:   "Assignment does not exist",
		})
	}
	for key, val := range Assignment {
		if err := d.Set(key, val); err != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to transfer data",
				Detail:   err.Error(),
			})
		}
	}

	return diags
}

func resourceAssignmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceAssignmentRead(ctx, d, m)
}

func resourceAssignmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	device := d.Get("device_udid").(string)
	data := map[string]interface{}{"operation": "change_to_limbo", "devices": [...]string{device}}
	req_body, err := json.Marshal(data)
	if err != nil {
		return diag.FromErr(err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/devices", c.HostURL), strings.NewReader(string(req_body)))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.doRequest(req)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func flattenAssignment(response ListResponse) map[string]interface{} {
	if len(response.Response) < 1 || len(response.Response[0].Devices) < 1 {
		return make(map[string]interface{}, 0)
	}

	oi := make(map[string]interface{})
	oi["os"] = response.Response[0].Devices[0]["os"]
	oi["user_id"] = response.Response[0].Devices[0]["idusermosyle"]
	oi["device_serial"] = response.Response[0].Devices[0]["serial_number"]
	oi["device_udid"] = response.Response[0].Devices[0]["deviceudid"]

	return oi
}
