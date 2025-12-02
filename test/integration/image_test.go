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
		TerraformDir: "../fixtures/vm",
		Vars: map[string]interface{}{
			"image_name": imageName,
			"project_id": os.Getenv("GCP_PROJECT"),
		},
	}

	terraform.InitAndApply(t, terraformOptions)
	defer terraform.Destroy(t, terraformOptions)

	vmIP := terraform.Output(t, terraformOptions, "vm_ip")

	// ✅ Wait for VM boot & Apache startup
	time.Sleep(30 * time.Second)

	url := "http://" + vmIP

	http_helper.HttpGetWithRetry(
		t,
		url,
		nil,
		200,
		"Apache Server from Golden Image",
		30,
		10*time.Second,
	)
}




















// package integration

// import (
// 	"os"
// 	"testing"
// 	"time"

// 	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
// 	"github.com/gruntwork-io/terratest/modules/terraform"
// )

// func TestApacheGoldenImage(t *testing.T) {
// 	t.Parallel()

// 	// ✅ Ensure IMAGE_NAME is passed from GitHub Actions
// 	imageName := os.Getenv("IMAGE_NAME")
// 	if imageName == "" {
// 		t.Fatal("❌ IMAGE_NAME is not set")
// 	}

// 	// ✅ Ensure GCP_PROJECT is passed
// 	projectID := os.Getenv("GCP_PROJECT")
// 	if projectID == "" {
// 		t.Fatal("❌ GCP_PROJECT is not set")
// 	}

// 	terraformOptions := &terraform.Options{
// 		TerraformDir: "../fixtures/vm",
// 		Vars: map[string]interface{}{
// 			"image_name": imageName,
// 			"project_id": projectID,
// 		},
// 	}

// 	// ✅ Create VM from Golden Image
// 	terraform.InitAndApply(t, terraformOptions)

// 	// ✅ Always destroy VM after test
// 	defer terraform.Destroy(t, terraformOptions)

// 	// ✅ Fetch VM Public IP
// 	vmIP := terraform.Output(t, terraformOptions, "vm_ip")

// 	url := "http://" + vmIP

// 	// ✅ Validate Apache Response (Plain Text)
// 	http_helper.HttpGetWithRetry(
// 		t,
// 		url,
// 		nil,
// 		200,
// 		"Apache Server from Golden Image", // MUST match index.html exactly
// 		15,                                // retries
// 		10*time.Second,                   // wait between retries
// 	)
// }


// package integration

// import (
// 	"os"
// 	"testing"
// 	"time"

// 	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
// 	"github.com/gruntwork-io/terratest/modules/terraform"
// )

// func TestApacheGoldenImage(t *testing.T) {
// 	t.Parallel()

// 	imageName := os.Getenv("IMAGE_NAME")
// 	if imageName == "" {
// 		t.Fatal("IMAGE_NAME is not set")
// 	}

// 	terraformOptions := &terraform.Options{
// 		TerraformDir: "../fixtures/vm", 
// 		Vars: map[string]interface{}{
// 			"image_name": imageName,
// 			"project_id": os.Getenv("GCP_PROJECT"),
// 		},
// 	}

// 	terraform.InitAndApply(t, terraformOptions)
// 	defer terraform.Destroy(t, terraformOptions)

// 	vmIP := terraform.Output(t, terraformOptions, "vm_ip")

// 	url := "http://" + vmIP

// 	http_helper.HttpGetWithRetry(
// 		t,
// 		url,
// 		nil,
// 		200,
// 		"Apache Server from Golden Image",
// 		20,
// 		10*time.Second,
// 	)
// }


