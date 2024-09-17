package data

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestWifPrincipalDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testWifPrincipalDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.wif_principal.actions", "id", "wif_principal:repo:reMarkable/actions"),
					resource.TestCheckResourceAttr("data.wif_principal.actions", "url", "principal://iam.googleapis.com/projects/1976/locations/global/workloadIdentityPools/pool/subject/repo:reMarkable/actions"),
				),
			},
		},
	})
}

const testWifPrincipalDataSourceConfig = `
provider "wif" {
  project_id = 1976
  pool_id = "pool"
}
data "wif_principal" "actions" {
  subject = "repo:reMarkable/actions"
}
`
