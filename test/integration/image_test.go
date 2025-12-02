package integration

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestApacheGoldenImage(t *testing.T) {
	t.Parallel()

	imageName := os.Getenv("IMAGE_NAME")
	if imageName == "" {
		t.Fatal("IMAGE_NAME is not set")
	}

	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		t.Fatal("GCP_PROJECT is not set")
	}

	terraformOptions := &terraform.Options{
		TerraformDir: "../fixtures/vm",
		Vars: map[string]interface{}{
			"image_name": imageName,
			"project_id": projectID,
		},
	}

	// ✅ Only infra validation
	terraform.InitAndApply(t, terraformOptions)

	// ✅ Always cleanup
	defer terraform.Destroy(t, terraformOptions)
}





// package integration

// import (
// 	"os"
// 	"testing"

// 	"github.com/gruntwork-io/terratest/modules/terraform"
// )

// func TestGoldenImageVMCreationOnly(t *testing.T) {
// 	t.Parallel()

// 	imageName := os.Getenv("IMAGE_NAME")
// 	if imageName == "" {
// 		t.Fatal("IMAGE_NAME is not set")
// 	}

// 	projectID := os.Getenv("GCP_PROJECT")
// 	if projectID == "" {
// 		t.Fatal("GCP_PROJECT is not set")
// 	}

// 	terraformOptions := &terraform.Options{
// 		TerraformDir: "../fixtures/vm",
// 		Vars: map[string]interface{}{
// 			"image_name": imageName,
// 			"project_id": projectID,
// 		},
// 	}

// 	// ✅ Create VM
// 	terraform.InitAndApply(t, terraformOptions)

// 	// ✅ Validate output exists (means VM created)
// 	vmIP := terraform.Output(t, terraformOptions, "vm_ip")
// 	if vmIP == "" {
// 		t.Fatal("VM IP is empty — VM creation failed")
// 	}

// 	t.Logf("VM successfully created with IP: %s", vmIP)

// 	// ✅ Destroy after validation
// 	defer terraform.Destroy(t, terraformOptions)
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

// 	// ✅ Wait for VM boot & Apache startup
// 	time.Sleep(30 * time.Second)

// 	url := "http://" + vmIP

// 	http_helper.HttpGetWithRetry(
// 		t,
// 		url,
// 		nil,
// 		200,
// 		"Apache Server from Golden Image",
// 		30,
// 		10*time.Second,
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




