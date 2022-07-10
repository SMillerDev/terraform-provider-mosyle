data "mosyle_usergroups" "all" {}

data "mosyle_usergroups" "some" {
  filter = {
    some = "val"
  }
}