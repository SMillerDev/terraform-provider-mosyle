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

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeMap,
				Description: "Filters to limit API data",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"users": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User data from Mosyle",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iduser":     &schema.Schema{Type: schema.TypeString, Computed: true},
						"code":       &schema.Schema{Type: schema.TypeString, Computed: true},
						"name":       &schema.Schema{Type: schema.TypeString, Computed: true},
						"type":       &schema.Schema{Type: schema.TypeString, Computed: true},
						"identifier": &schema.Schema{Type: schema.TypeString, Computed: true},
						"email":      &schema.Schema{Type: schema.TypeString, Computed: true},
						"is_removed": &schema.Schema{Type: schema.TypeBool, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	filter := d.Get("filter").(map[string]interface{})
	req_body, err := json.Marshal(ListPostBody{
		Operation: "list_users",
		Options:   filter,
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

	if err := d.Set("users", flattenUsers(response)); err != nil {
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

func flattenUsers(response ListResponse) []map[string]interface{} {
	if len(response.Response) < 1 || len(response.Response[0].Users) < 1 {
		return make([]map[string]interface{}, 0)
	}

	ois := make([]map[string]interface{}, len(response.Response[0].Users), len(response.Response[0].Users))

	for i, user := range response.Response[0].Users {
		oi := make(map[string]interface{})

		for key, val := range user {
			oi[strings.ToLower(key)] = val
		}

		ois[i] = oi
	}

	return ois
}
