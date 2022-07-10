data "mosyle_devices" "all" {
  filter = {
    os = "mac"
  }
}