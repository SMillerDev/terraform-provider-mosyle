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

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Description:   "User data",
		Schema: map[string]*schema.Schema{
			"name":       &schema.Schema{Type: schema.TypeString, Required: true, Description: "User name",},
			"identifier": &schema.Schema{Type: schema.TypeString, Required: true, Description: "User identifier, set by admin",},
			"email":      &schema.Schema{Type: schema.TypeString, Optional: true, Description: "User email",},
			"type":       &schema.Schema{Type: schema.TypeString, Optional: true, Default: "ENDUSER", Description: "User type, one of (ENDUSER|GROUP_ADMIN|ADMIN) default: ENDUSER",},
			"iduser":     &schema.Schema{Type: schema.TypeString, Computed: true, Description: "User id from mosyle",},
			"code":       &schema.Schema{Type: schema.TypeString, Computed: true, Description: "User code",},
			"is_removed": &schema.Schema{Type: schema.TypeBool, Computed: true, Description: "User is removed",},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	id := d.Get("identifier").(string)
	user_type := d.Get("type").(string)

	data := map[string]string{"operation": "create_user", "user_id": id, "type": user_type, "name": name}
	req_body, err := json.Marshal(data)
	if err != nil {
		return diag.FromErr(err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users", c.HostURL), strings.NewReader(string(req_body)))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.doRequest(req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	resourceUserRead(ctx, d, m)

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req_body, err := json.Marshal(ListUserBody{
		Operation: "list_users",
		Options:   map[string][]string{"identifiers": []string{d.Id()}},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users", c.HostURL), strings.NewReader(string(req_body)))
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := c.doRequest(req)
	if err != nil {
		return diag.FromErr(err)
	}

	users := flattenUsers(response)
	if len(users) < 1 {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "No such user",
			Detail:   "User does not exist",
		})
	}
	for key, val := range users[0] {
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

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
