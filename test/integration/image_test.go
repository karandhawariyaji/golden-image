package integration

import (
	"os"
	"os/exec"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestGoldenImageInfraWithOPA(t *testing.T) {
	t.Parallel()

	// -----------------------------
	// ✅ 1. READ ENV VARS
	// -----------------------------
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

	// -----------------------------
	// ✅ 2. TERRAFORM INIT
	// -----------------------------
	terraform.Init(t, terraformOptions)

	// -----------------------------
	// ✅ 3. TERRAFORM PLAN
	// -----------------------------
	terraform.RunTerraformCommand(
		t,
		terraformOptions,
		terraform.FormatArgs(terraformOptions, "plan", "-out=tfplan")...,
	)

	// Convert plan to JSON
	terraform.RunTerraformCommand(
		t,
		terraformOptions,
		"show", "-json", "tfplan", "tfplan.json",
	)

	// -----------------------------
	// ✅ 4. OPA POLICY TEST
	// -----------------------------
	t.Log("Running OPA policy test on Terraform plan...")

	opaCmd := exec.Command(
		"conftest",
		"test",
		"tfplan.json",
		"--policy",
		"../../opa-policies",
	)

	opaCmd.Dir = terraformOptions.TerraformDir
	opaOutput, err := opaCmd.CombinedOutput()

	if err != nil {
		t.Fatalf("❌ OPA policy violation:\n%s", string(opaOutput))
	}

	t.Log("✅ OPA policy validation passed")

	// -----------------------------
	// ✅ 5. TERRAFORM APPLY
	// -----------------------------
	terraform.Apply(t, terraformOptions)

	// -----------------------------
	// ✅ 6. TERRAFORM DESTROY
	// -----------------------------
	defer terraform.Destroy(t, terraformOptions)

	t.Log("✅ Infra created, validated, and will be destroyed")
}














// package integration

// import (
// 	"os"
// 	"testing"

// 	"github.com/gruntwork-io/terratest/modules/terraform"
// )

// func TestApacheGoldenImage(t *testing.T) {
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

// 	// ✅ Only infra validation
// 	terraform.InitAndApply(t, terraformOptions)

// 	// ✅ Always cleanup
// 	defer terraform.Destroy(t, terraformOptions)
// }





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





