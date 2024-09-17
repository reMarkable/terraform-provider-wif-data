package data

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestWifPrincipalSetDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testWifPrincipalSetDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.wif_principal_set.actions", "id", "wif_principal_set:repo:assertion.arn.contains(\":instance-profile/Production\")"),
					resource.TestCheckResourceAttr("data.wif_principal_set.actions", "url", "principalSet://iam.googleapis.com/projects/1976/locations/global/workloadIdentityPools/mypool/repo/assertion.arn.contains(\":instance-profile/Production\")"),
					resource.TestCheckResourceAttr("data.wif_principal_set.mygroup", "url", "principalSet://iam.googleapis.com/projects/1976/locations/global/workloadIdentityPools/mypool/attribute.group/mygroup"),
				),
			},
		},
	})
}

const testWifPrincipalSetDataSourceConfig = `
provider "wif" {
  project_id = 1976
  pool_id = "mypool"
}
data "wif_principal_set" "actions" {
  target = "repo"
  source_expression = "assertion.arn.contains(\":instance-profile/Production\")"
}
data "wif_principal_set" "mygroup" {
  target = "attribute.group"
  source_expression = "mygroup"
}
`
