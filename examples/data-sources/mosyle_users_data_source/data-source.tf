data "mosyle_users" "all" {}

data "mosyle_users" "some" {
  filter = {
    some = "val"
  }
}