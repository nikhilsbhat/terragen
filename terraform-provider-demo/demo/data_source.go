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
)

func dataSourceDEMO() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDEMORead,

		Schema: map[string]*schema.Schema{},
	}
}

func dataSourceDEMORead(d *schema.ResourceData, meta interface{}) error {
	// Your code goes here
	return nil
}