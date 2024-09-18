# Terraform provider for workload identity federation

This is a WIP Terraform provider for generating [Workload Identity Federation](https://cloud.google.com/iam/docs/workload-identity-federation) principals.

This will attempt some minimal validation of source_expression being a valid CEL expression using [go-cel](https://github.com/google/cel-go)

It can be used like this:

```hcl
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

data "wif_principal" "actions" {
  subject = "repo:reMarkable/actions"
}
```

Note that this provider has not yet been uploaded to a registry.
