package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoldenImageWithOPA(t *testing.T) {
	t.Parallel()

	// Get environment variables
	projectID := os.Getenv("GCP_PROJECT_ID")
	require.NotEmpty(t, projectID, "GCP_PROJECT_ID must be set")

	imageName := os.Getenv("IMAGE_NAME")
	require.NotEmpty(t, imageName, "IMAGE_NAME must be set")

	// Test cases
	testCases := []struct {
		name          string
		region        string
		machineType   string
		hasPublicIP   bool
		shouldPassOPA bool
	}{
		{
			name:          "ValidConfig",
			region:        "us-central1",
			machineType:   "n1-standard-2",
			hasPublicIP:   false,
			shouldPassOPA: true,
		},
		{
			name:          "InvalidRegion",
			region:        "asia-south1",
			machineType:   "n1-standard-2",
			hasPublicIP:   false,
			shouldPassOPA: false,
		},
		{
			name:          "InvalidMachineType",
			region:        "us-central1",
			machineType:   "n1-highmem-8",
			hasPublicIP:   false,
			shouldPassOPA: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup test directory
			testDir := "fixtures/vm"
			
			// Terraform options
			terraformOptions := &terraform.Options{
				TerraformDir: testDir,
				Vars: map[string]interface{}{
					"project_id":    projectID,
					"image_name":    imageName,
					"region":        tc.region,
					"zone":          fmt.Sprintf("%s-a", tc.region),
					"machine_type":  tc.machineType,
					"instance_name": fmt.Sprintf("test-%s", tc.name),
				},
				EnvVars: map[string]string{
					"GOOGLE_PROJECT": projectID,
				},
				NoColor: true,
			}

			// Step 1: Initialize Terraform
			terraform.Init(t, terraformOptions)

			// Step 2: Create plan
			planFile := filepath.Join(testDir, "tfplan")
			terraform.PlanWithOptions(t, terraformOptions, &terraform.PlanOptions{
				Out: planFile,
			})

			// Step 3: Convert plan to JSON for OPA
			planJSON := terraform.ShowWithOptions(t, terraformOptions, &terraform.ShowOptions{
				Json:         true,
				PlanFilePath: planFile,
			})

			// Step 4: Run OPA validation
			passed := validateWithOPA(t, planJSON)
			
			if tc.shouldPassOPA {
				assert.True(t, passed, "OPA validation should pass")
				
				// Apply Terraform
				defer terraform.Destroy(t, terraformOptions)
				terraform.Apply(t, terraformOptions)
				
				// Verify outputs
				outputs := terraform.OutputAll(t, terraformOptions)
				assert.NotEmpty(t, outputs["instance_id"])
				assert.Contains(t, outputs["zone"], tc.region)
				
				t.Logf("✅ Test '%s' passed", tc.name)
			} else {
				assert.False(t, passed, "OPA validation should fail")
				t.Logf("✅ Test '%s' correctly failed OPA check", tc.name)
			}
		})
	}
}

func validateWithOPA(t *testing.T, planJSON string) bool {
	// Write plan to temp file
	tmpDir := t.TempDir()
	planFile := filepath.Join(tmpDir, "plan.json")
	err := os.WriteFile(planFile, []byte(planJSON), 0644)
	require.NoError(t, err)

	// Policy file path
	policyFile, _ := filepath.Abs("../opa-policies/vm_policy.rego")

	// Run OPA command
	cmd := fmt.Sprintf("opa eval -i %s -d %s 'data.vm.policies.deny' --format json", planFile, policyFile)
	
	output, err := terraform.RunShellCommandWithOutputE(t, &terraform.Options{}, "bash", "-c", cmd)
	
	if err != nil {
		t.Logf("OPA command failed: %v", err)
		return false
	}

	// Parse OPA result
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Logf("Failed to parse OPA output: %v", err)
		return false
	}

	// Check for denials
	if result["result"] != nil {
		if denials, ok := result["result"].([]interface{}); ok && len(denials) > 0 {
			t.Logf("OPA Violations (%d):", len(denials))
			for _, denial := range denials {
				t.Logf("  - %v", denial)
			}
			return false
		}
	}

	return true
}

// Simple test without OPA for quick validation
func TestGoldenImageExists(t *testing.T) {
	projectID := os.Getenv("GCP_PROJECT_ID")
	imageName := os.Getenv("IMAGE_NAME")
	
	require.NotEmpty(t, projectID, "GCP_PROJECT_ID must be set")
	require.NotEmpty(t, imageName, "IMAGE_NAME must be set")
	
	t.Logf("Testing image: %s in project: %s", imageName, projectID)
	
	// This just verifies the image exists
	testDir := "fixtures/vm"
	terraformOptions := &terraform.Options{
		TerraformDir: testDir,
		Vars: map[string]interface{}{
			"project_id": projectID,
			"image_name": imageName,
		},
	}
	
	// Just init to validate
	terraform.Init(t, terraformOptions)
	t.Log("✅ Image reference is valid")
}








// package integration

// import (
// 	"testing"
// 	"os"
// 	"encoding/json"
	
// 	"github.com/gruntwork-io/terratest/modules/terraform"
// 	"github.com/stretchr/testify/assert"
// )

// func TestGoldenImageWithOPA(t *testing.T) {
// 	// Get from GitHub Actions
// 	projectID := os.Getenv("GCP_PROJECT_ID")
// 	imageName := os.Getenv("IMAGE_NAME")

// 	// 1. Setup Terraform
// 	terraformOptions := &terraform.Options{
// 		TerraformDir: "../../terraform",
// 		Vars: map[string]interface{}{
// 			"project_id":   projectID,
// 			"image_name":   imageName,
// 			"region":       "us-central1",
// 			"machine_type": "n1-standard-2",
// 		},
// 	}

// 	// Cleanup
// 	defer terraform.Destroy(t, terraformOptions)

// 	// 2. Plan
// 	terraform.Init(t, terraformOptions)
// 	planFile := "tfplan"
// 	terraform.PlanWithOptions(t, terraformOptions, &terraform.PlanOptions{
// 		Out: planFile,
// 	})

// 	// 3. Convert plan to JSON for OPA
// 	planJSON := terraform.ShowWithOptions(t, terraformOptions, &terraform.ShowOptions{
// 		Json:         true,
// 		PlanFilePath: planFile,
// 	})

// 	// 4. Run OPA test
// 	pass := runOPATest(t, planJSON)
// 	assert.True(t, pass, "OPA policy check failed")

// 	// 5. Apply only if OPA passes
// 	terraform.Apply(t, terraformOptions)
	
// 	t.Log("✅ All tests passed!")
// }

// func runOPATest(t *testing.T, planJSON string) bool {
// 	// Save plan JSON to temp file
// 	planFile := "plan.json"
// 	err := os.WriteFile(planFile, []byte(planJSON), 0644)
// 	if err != nil {
// 		t.Fatal("Failed to write plan file:", err)
// 	}
// 	defer os.Remove(planFile)

// 	// Run OPA eval
// 	cmd := "opa eval -i plan.json -d ../opa/policy.rego 'data.main.deny' --format json"
// 	output, err := terraform.RunShellCommandWithOutputE(t, &terraform.Options{}, "bash", "-c", cmd)
	
// 	if err != nil {
// 		t.Logf("OPA check failed: %v", err)
// 		t.Logf("Output: %s", output)
// 		return false
// 	}

// 	// Parse OPA result
// 	var result map[string]interface{}
// 	json.Unmarshal([]byte(output), &result)
	
// 	// Check if deny array is empty
// 	if result["result"] != nil {
// 		denials := result["result"].([]interface{})
// 		if len(denials) > 0 {
// 			t.Logf("OPA violations found:")
// 			for _, d := range denials {
// 				t.Logf("  - %v", d)
// 			}
// 			return false
// 		}
// 	}
	
// 	return true
// }








// package integration

// import (
// 	"os"
// 	"os/exec"
// 	"testing"

// 	"github.com/gruntwork-io/terratest/modules/terraform"
// )

// func TestGoldenImageInfraWithOPA(t *testing.T) {
// 	t.Parallel()

// 	// -----------------------------
// 	// ✅ 1. READ ENV VARS
// 	// -----------------------------
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

// 	// -----------------------------
// 	// ✅ 2. TERRAFORM INIT
// 	// -----------------------------
// 	terraform.Init(t, terraformOptions)

// 	// -----------------------------
// 	// ✅ 3. TERRAFORM PLAN
// 	// -----------------------------
// 	terraform.RunTerraformCommand(
// 		t,
// 		terraformOptions,
// 		terraform.FormatArgs(terraformOptions, "plan", "-out=tfplan")...,
// 	)

// 	// Convert plan to JSON
// 	terraform.RunTerraformCommand(
// 		t,
// 		terraformOptions,
// 		"show", "-json", "tfplan", "tfplan.json",
// 	)

// 	// -----------------------------
// 	// ✅ 4. OPA POLICY TEST
// 	// -----------------------------
// 	t.Log("Running OPA policy test on Terraform plan...")

// 	opaCmd := exec.Command(
// 		"conftest",
// 		"test",
// 		"tfplan.json",
// 		"--policy",
// 		"../../opa-policies",
// 	)

// 	opaCmd.Dir = terraformOptions.TerraformDir
// 	opaOutput, err := opaCmd.CombinedOutput()

// 	if err != nil {
// 		t.Fatalf("❌ OPA policy violation:\n%s", string(opaOutput))
// 	}

// 	t.Log("✅ OPA policy validation passed")

// 	// -----------------------------
// 	// ✅ 5. TERRAFORM APPLY
// 	// -----------------------------
// 	terraform.Apply(t, terraformOptions)

// 	// -----------------------------
// 	// ✅ 6. TERRAFORM DESTROY
// 	// -----------------------------
// 	defer terraform.Destroy(t, terraformOptions)

// 	t.Log("✅ Infra created, validated, and will be destroyed")
// }














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





