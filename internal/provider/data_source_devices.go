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

func dataSourceDevices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDevicesRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeMap,
				Description: "Filters to limit API data",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"devices": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Device data from Mosyle, this can be macOS, iOS or tvOS",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deviceudid":                       &schema.Schema{Type: schema.TypeString, Computed: true},
						"total_disk":                       &schema.Schema{Type: schema.TypeString, Computed: true},
						"os":                               &schema.Schema{Type: schema.TypeString, Computed: true},
						"serial_number":                    &schema.Schema{Type: schema.TypeString, Computed: true},
						"device_model_name":                &schema.Schema{Type: schema.TypeString, Computed: true},
						"device_name":                      &schema.Schema{Type: schema.TypeString, Computed: true},
						"device_model":                     &schema.Schema{Type: schema.TypeString, Computed: true},
						"battery":                          &schema.Schema{Type: schema.TypeString, Computed: true},
						"osversion":                        &schema.Schema{Type: schema.TypeString, Computed: true},
						"vpn_status":                       &schema.Schema{Type: schema.TypeString, Computed: true},
						"userid":                           &schema.Schema{Type: schema.TypeString, Computed: true},
						"date_info":                        &schema.Schema{Type: schema.TypeString, Computed: true},
						"carrier":                          &schema.Schema{Type: schema.TypeString, Computed: true},
						"roaming_enabled":                  &schema.Schema{Type: schema.TypeString, Computed: true},
						"isroaming":                        &schema.Schema{Type: schema.TypeString, Computed: true},
						"imei":                             &schema.Schema{Type: schema.TypeString, Computed: true},
						"meid":                             &schema.Schema{Type: schema.TypeString, Computed: true},
						"available_disk":                   &schema.Schema{Type: schema.TypeString, Computed: true},
						"wifi_mac_address":                 &schema.Schema{Type: schema.TypeString, Computed: true},
						"bluetooth_mac_address":            &schema.Schema{Type: schema.TypeString, Computed: true},
						"is_supervised":                    &schema.Schema{Type: schema.TypeBool, Computed: true},
						"date_app_info":                    &schema.Schema{Type: schema.TypeString, Computed: true},
						"date_last_beat":                   &schema.Schema{Type: schema.TypeString, Computed: true},
						"date_last_push":                   &schema.Schema{Type: schema.TypeString, Computed: true},
						"status":                           &schema.Schema{Type: schema.TypeString, Computed: true},
						"isactivationlockenabled":          &schema.Schema{Type: schema.TypeString, Computed: true},
						"isdevicelocatorserviceenabled":    &schema.Schema{Type: schema.TypeString, Computed: true},
						"isdonotdisturbineffect":           &schema.Schema{Type: schema.TypeString, Computed: true},
						"iscloudbackupenabled":             &schema.Schema{Type: schema.TypeString, Computed: true},
						"isnetworktethered":                &schema.Schema{Type: schema.TypeString, Computed: true},
						"needosupdate":                     &schema.Schema{Type: schema.TypeString, Computed: true},
						"productkeyupdate":                 &schema.Schema{Type: schema.TypeString, Computed: true},
						"device_type":                      &schema.Schema{Type: schema.TypeString, Computed: true},
						"lostmode_status":                  &schema.Schema{Type: schema.TypeString, Computed: true},
						"is_muted":                         &schema.Schema{Type: schema.TypeBool, Computed: true},
						"date_muted":                       &schema.Schema{Type: schema.TypeString, Computed: true},
						"activation_bypass":                &schema.Schema{Type: schema.TypeString, Computed: true},
						"date_media_info":                  &schema.Schema{Type: schema.TypeString, Computed: true},
						"tags":                             &schema.Schema{Type: schema.TypeString, Computed: true},
						"is_deleted":                       &schema.Schema{Type: schema.TypeBool, Computed: true},
						"itunesstoreaccounthash":           &schema.Schema{Type: schema.TypeString, Computed: true},
						"itunesstoreaccountisactive":       &schema.Schema{Type: schema.TypeString, Computed: true},
						"date_profiles_info":               &schema.Schema{Type: schema.TypeString, Computed: true},
						"ethernet_mac_address":             &schema.Schema{Type: schema.TypeString, Computed: true},
						"model_name":                       &schema.Schema{Type: schema.TypeString, Computed: true},
						"lastcloudbackupdate":              &schema.Schema{Type: schema.TypeString, Computed: true},
						"systemintegrityprotectionenabled": &schema.Schema{Type: schema.TypeString, Computed: true},
						"buildversion":                     &schema.Schema{Type: schema.TypeString, Computed: true},
						"localhostname":                    &schema.Schema{Type: schema.TypeString, Computed: true},
						"hostname":                         &schema.Schema{Type: schema.TypeString, Computed: true},
						"osupdatesettings":                 &schema.Schema{Type: schema.TypeString, Computed: true},
						"activemanagedusers":               &schema.Schema{Type: schema.TypeString, Computed: true},
						"currentconsolemanageduser":        &schema.Schema{Type: schema.TypeString, Computed: true},
						"date_printers":                    &schema.Schema{Type: schema.TypeString, Computed: true},
						"autosetupadminaccounts":           &schema.Schema{Type: schema.TypeString, Computed: true},
						"appletvid":                        &schema.Schema{Type: schema.TypeString, Computed: true},
						"asset_tag":                        &schema.Schema{Type: schema.TypeString, Computed: true},
						"managementstatus":                 &schema.Schema{Type: schema.TypeString, Computed: true},
						"osupdatestatus":                   &schema.Schema{Type: schema.TypeString, Computed: true},
						"availableosupdates":               &schema.Schema{Type: schema.TypeString, Computed: true},
						"has_password":                     &schema.Schema{Type: schema.TypeString, Computed: true},
						"timezone":                         &schema.Schema{Type: schema.TypeString, Computed: true},
						"activation_bypass_mdm":            &schema.Schema{Type: schema.TypeString, Computed: true},
						"percent_disk":                     &schema.Schema{Type: schema.TypeString, Computed: true},
						"idsharedgroup":                    &schema.Schema{Type: schema.TypeString, Computed: true},
						"enrollment_type":                  &schema.Schema{Type: schema.TypeString, Computed: true},
						"status_login":                     &schema.Schema{Type: schema.TypeString, Computed: true},
						"date_lastlogin":                   &schema.Schema{Type: schema.TypeString, Computed: true},
						"idaccount":                        &schema.Schema{Type: schema.TypeString, Computed: true},
						"date_checkin":                     &schema.Schema{Type: schema.TypeString, Computed: true},
						"date_enroll":                      &schema.Schema{Type: schema.TypeString, Computed: true},
						"date_checkout":                    &schema.Schema{Type: schema.TypeString, Computed: true},
						"date_kinfo":                       &schema.Schema{Type: schema.TypeString, Computed: true},
						"cpu_model":                        &schema.Schema{Type: schema.TypeString, Computed: true},
						"hasvpn":                           &schema.Schema{Type: schema.TypeString, Computed: true},
						"installed_memory":                 &schema.Schema{Type: schema.TypeString, Computed: true},
						"username":                         &schema.Schema{Type: schema.TypeString, Computed: true},
						"usertype":                         &schema.Schema{Type: schema.TypeString, Computed: true},
						"idusermosyle":                     &schema.Schema{Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceDevicesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	filter := d.Get("filter").(map[string]interface{})
	options := make(map[string]interface{})
	for key, value := range filter {
		options[key] = value
	}

	req_data := ListPostBody{
		Operation: "list",
		Options:   options,
	}
	req_body, err := json.Marshal(req_data)
	if err != nil {
		return diag.FromErr(err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/devices", c.HostURL), strings.NewReader(string(req_body)))
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

	if err := d.Set("devices", flattenDevices(response)); err != nil {
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

func flattenDevices(response ListResponse) []interface{} {
	if len(response.Response) < 1 || len(response.Response[0].Devices) < 1 {
		return make([]interface{}, 0)
	}

	ois := make([]interface{}, len(response.Response[0].Devices), len(response.Response[0].Devices))

	for i, device := range response.Response[0].Devices {
		oi := make(map[string]interface{})

		for key, val := range device {
			if strings.HasPrefix(key, "is_") {
				if val == nil {
					oi[strings.ToLower(key)] = false
					continue
				}
				bval, err := strconv.ParseBool(val.(string))
				if err == nil {
					oi[strings.ToLower(key)] = bval
					continue
				}
			}

			if strings.HasPrefix(key, "date_") && val != nil {
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
