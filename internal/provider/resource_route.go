package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/matthisholleville/terraform-pritunl-provider/internal/pritunl"
)

func resourceRoute() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the server",
			},
			"network": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Network address with subnet to route",
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					return validation.IsCIDR(i, s)
				},
			},
			"comment": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Comment for route",
			},
			"nat": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Description: "NAT vpn traffic destined to this network",
				Computed:    true,
			},
		},
		CreateContext: resourceCreateRoute,
		ReadContext:   resourceReadRoute,
		UpdateContext: resourceUpdateRoute,
		DeleteContext: resourceDeleteRoute,
	}
}

func resourceCreateRoute(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(pritunl.Client)

	routeData := map[string]interface{}{
		"network": d.Get("network"),
		"comment": d.Get("comment"),
		"nat":     d.Get("nat"),
	}

	serverId := d.Get("server_id").(string)

	err := apiClient.StopServer(serverId)
	if err != nil {
		return diag.FromErr(err)
	}

	route, err := apiClient.AddRouteToServer(serverId, pritunl.ConvertMapToRoute(routeData))
	if err != nil {
		return diag.FromErr(err)
	}

	err = apiClient.StartServer(serverId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(route.GetID())

	return resourceReadRoute(ctx, d, meta)
}

func resourceReadRoute(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(pritunl.Client)

	serverId := d.Get("server_id").(string)

	routes, err := apiClient.GetRoutesByServer(serverId)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, route := range routes {
		if route.GetID() == d.Id() {
			d.Set("server_id", serverId)
			d.Set("network", route.Network)
			d.Set("comment", route.Comment)
			d.Set("nat", route.Nat)
		}

		return nil
	}
	return diag.FromErr(errors.New("Unable to find route."))
}

func resourceUpdateRoute(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(pritunl.Client)

	serverId := d.Get("server_id").(string)

	routes, err := apiClient.GetRoutesByServer(serverId)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, route := range routes {
		if route.GetID() == d.Id() {
			if v, ok := d.GetOk("nat"); ok {
				route.Nat = v.(bool)
			}

			if v, ok := d.GetOk("comment"); ok {
				route.Comment = v.(string)
			}

			err = apiClient.StopServer(serverId)
			if err != nil {
				return diag.FromErr(err)
			}

			route, err := apiClient.UpdateRouteOnServer(serverId, route)
			if err != nil {
				return diag.FromErr(err)
			}

			err = apiClient.StartServer(serverId)
			if err != nil {
				return diag.FromErr(err)
			}

			d.SetId(route.GetID())

			return resourceReadRoute(ctx, d, meta)
		}
	}
	return diag.FromErr(errors.New("Unable to find route."))
}

func resourceDeleteRoute(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(pritunl.Client)

	serverId := d.Get("server_id").(string)

	err := apiClient.StopServer(serverId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = apiClient.DeleteRouteFromServer(serverId, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = apiClient.StartServer(serverId)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil

}