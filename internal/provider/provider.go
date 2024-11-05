package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/matthisholleville/terraform-pritunl-provider/internal/pritunl"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PRITUNL_URL", ""),
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PRITUNL_TOKEN", ""),
			},
			"secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PRITUNL_SECRET", ""),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PRITUNL_INSECURE", false),
			},
			"connection_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PRITUNL_CONNECTION_CHECK", true),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"pritunl_route":        resourceRoute(),
			"pritunl_server":       resourceServer(),
			"pritunl_organization": resourceOrganization(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"pritunl_route":   dataSourceRoute(),
			"pritunl_server": dataSourceServer(),
			"pritunl_host":   dataSourceHost(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := d.Get("url").(string)
	token := d.Get("token").(string)
	secret := d.Get("secret").(string)
	insecure := d.Get("insecure").(bool)
	connectionCheck := d.Get("connection_check").(bool)

	apiClient := pritunl.NewClient(url, token, secret, insecure)

	if connectionCheck {
		// execute test api call to ensure that provided credentials are valid and pritunl api works
		err := apiClient.TestApiCall()
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	return apiClient, nil
}
