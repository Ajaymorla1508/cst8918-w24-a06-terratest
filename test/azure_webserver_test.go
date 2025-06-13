package test

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

var subscriptionID string = "ec3ee3c5-b6d0-44db-98c1-49df6c18e116"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"labelPrefix": "ajaymorla",
		},
	}

	defer func() {
		t.Log(" Waiting before terraform destroy to avoid Azure race conditions...")
		time.Sleep(30 * time.Second)

		t.Log("Attempting terraform destroy...")
		_, err := terraform.DestroyE(t, terraformOptions)
		if err != nil {
			t.Logf("Terraform destroy failed: %v", err)
			t.Log("Waiting 60s before retrying destroy...")
			time.Sleep(60 * time.Second)
			t.Log("Retrying terraform destroy...")
			terraform.Destroy(t, terraformOptions)
		} else {
			t.Log("Terraform destroy succeeded.")
		}

		t.Log(" Final wait after destroy to ensure Azure releases all resources...")
		time.Sleep(60 * time.Second)
	}()

	// Run `terraform init` and `terraform apply`
	t.Log("Running terraform init and apply...")
	terraform.InitAndApply(t, terraformOptions)

	// Get output variables from Terraform
	t.Log("Fetching output variables...")
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Validate the VM exists
	t.Logf("Verifying if VM %s exists in resource group %s...", vmName, resourceGroupName)
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID), "Expected VM to exist but it does not.")
}
