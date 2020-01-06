// ----------------------------------------------------------------------------
//
//     ***     TERRAGEN GENERATED CODE    ***    TERRAGEN GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file was generated automatically by Terragen.
//     It has to be enhanced further to make it fully working terraform-provider.
//
//     Get more information on how terragen works.
//     https://github.com/nikhilsbhat/terragen
//
// ----------------------------------------------------------------------------
package demo

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},

		ResourcesMap: map[string]*schema.Resource{
			"null_resource": resourceDEMO(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"null_data_source": dataSourceDEMO(),
		},
	}
}