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
package main

import (
	"terraform-provider-demo/demo"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: demo.Provider})
}