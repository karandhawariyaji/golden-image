package integration

import (
  "testing"
  "time"

  "github.com/gruntwork-io/terratest/modules/terraform"
  "github.com/gruntwork-io/terratest/modules/http-helper"
)

func TestGoldenImageVM(t *testing.T) {
  t.Parallel()

  tfOptions := &terraform.Options{
    TerraformDir: "../fixtures/vm",

    Vars: map[string]interface{}{
      "project_id": "your-project-id",
      "zone":       "asia-south1-a",
      "image_name": "your-golden-image-name",
    },
  }

  // ✅ Always delete test VM
  defer terraform.Destroy(t, tfOptions)

  // ✅ Create test VM
  terraform.InitAndApply(t, tfOptions)

  // ✅ Fetch VM Public IP
  vmIP := terraform.Output(t, tfOptions, "vm_ip")

  // ✅ Verify Apache is reachable
  url := "http://" + vmIP

  http_helper.HttpGetWithRetry(
    t,
    url,
    nil,
    200,
    "",
    30,
    5*time.Second,
  )
}
