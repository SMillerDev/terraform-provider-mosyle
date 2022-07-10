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

func dataSourceUserGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserGroupsRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeMap,
				Description: "Filters to limit API data",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"groups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User group data from Mosyle",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"idusergroup":        &schema.Schema{Type: schema.TypeString, Computed: true},
						"identifier":         &schema.Schema{Type: schema.TypeString, Computed: true},
						"name":               &schema.Schema{Type: schema.TypeString, Computed: true},
						"idusergroup_parent": &schema.Schema{Type: schema.TypeString, Computed: true},
						"date_created":       &schema.Schema{Type: schema.TypeString, Computed: true},
						"date_modified":      &schema.Schema{Type: schema.TypeString, Computed: true},
						"is_removed":         &schema.Schema{Type: schema.TypeBool, Computed: true},
						"idusers_primary": &schema.Schema{Type: schema.TypeList, Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							}},
					},
				},
			},
		},
	}
}

func dataSourceUserGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	filter := d.Get("filter").(map[string]interface{})
	options := make(map[string]interface{})
	for key, value := range filter {
		options[key] = value
	}

	req_data := ListPostBody{
		Operation: "list_usergroup",
		Options:   options,
	}
	req_body, err := json.Marshal(req_data)
	if err != nil {
		return diag.FromErr(err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/usergroups", c.HostURL), strings.NewReader(string(req_body)))
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := c.doRequest(req)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
			Detail:   fmt.Sprint(string(req_body)),
		})
	}

	if err := d.Set("groups", flattenUserGroups(response)); err != nil {
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

func flattenUserGroups(response ListResponse) []interface{} {
	if len(response.Response) < 1 || len(response.Response[0].UserGroups) < 1 {
		return make([]interface{}, 0)
	}

	ois := make([]interface{}, len(response.Response[0].UserGroups), len(response.Response[0].UserGroups))

	for i, user := range response.Response[0].UserGroups {
		oi := make(map[string]interface{})

		for key, val := range user {
			if strings.HasPrefix(key, "date_") {
				ival, err := strconv.ParseInt(val.(string), 10, 64)
				if err == nil {
					t := time.Unix(ival, 0)
					strDate := t.Format(time.RFC3339)
					oi[strings.ToLower(key)] = strDate
					continue
				}
			}

			oi[strings.ToLower(key)] = val
		}

		ois[i] = oi
	}

	return ois
}
