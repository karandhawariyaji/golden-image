package integration

import (
	"testing"
	"time"
	"os"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/http-helper"
)

func TestApacheGoldenImage(t *testing.T) {
	t.Parallel()

	imageName := os.Getenv("IMAGE_NAME")
	if imageName == "" {
		t.Fatal("IMAGE_NAME is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../../terraform/test-vm",
		Vars: map[string]interface{}{
			"image_name": imageName,
		},
	}

	// Create Test VM
	terraform.InitAndApply(t, terraformOptions)

	// Destroy VM After Test
	defer terraform.Destroy(t, terraformOptions)

	vmIP := terraform.Output(t, terraformOptions, "vm_ip")

	url := "http://" + vmIP

	// Retry for 3 minutes
	http_helper.HttpGetWithRetry(
		t,
		url,
		nil,
		200,
		3*time.Minute,
		10*time.Second,
	)
}
