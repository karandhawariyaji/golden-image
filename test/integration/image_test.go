package integration

import (
	"os"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestApacheGoldenImage(t *testing.T) {
	t.Parallel()

	imageName := os.Getenv("IMAGE_NAME")
	if imageName == "" {
		t.Fatal("IMAGE_NAME is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../fixtures/vm",
		Vars: map[string]interface{}{
			"image_name":  imageName,
			"project_id":  os.Getenv("GCP_PROJECT"),
		},
	}

	terraform.InitAndApply(t, terraformOptions)
	defer terraform.Destroy(t, terraformOptions)

	vmIP := terraform.Output(t, terraformOptions, "vm_ip")

	url := "http://" + vmIP

	http_helper.HttpGetWithRetry(
		t,
		url,
		nil,
		200,
		"Apache",
		20,
		10*time.Second,
	)
}


