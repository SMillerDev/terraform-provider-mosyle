package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"username": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Username used to log in to Mosyle",
					DefaultFunc: schema.EnvDefaultFunc("MOSYLE_USERNAME", nil),
				},
				"password": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					Description: "Password used to log in to Mosyle",
					DefaultFunc: schema.EnvDefaultFunc("MOSYLE_PASSWORD", nil),
				},
				"accesstoken": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					Description: "Access Token from the Mosyle API integration",
					DefaultFunc: schema.EnvDefaultFunc("MOSYLE_TOKEN", nil),
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"mosyle_user": resourceUser(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"mosyle_devices":      dataSourceDevices(),
				"mosyle_devicegroups": dataSourceDeviceGroups(),
				"mosyle_users":        dataSourceUsers(),
				"mosyle_usergroups":   dataSourceUserGroups(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	// Add whatever fields, client or connection info, etc. here
	// you would need to setup to communicate with the upstream
	// API.
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		username := d.Get("username").(string)
		password := d.Get("password").(string)
		accesstoken := d.Get("accesstoken").(string)

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		if (username != "") && (password != "") && (accesstoken != "") {
			c, err := MosyleClient(version, &username, &password, &accesstoken)

			return c, diag.FromErr(err)
		}

		c, err := MosyleClient(version, nil, nil, nil)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return c, diags
	}
}
