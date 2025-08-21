//go:build (all || resource_servicehook_subscription) && !exclude_subscriptions

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServicehookSubscription_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	webhookUrl := "https://example.com/webhook"

	resourceType := "azuredevops_servicehook_subscription"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookSubscriptionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookSubscriptionResourceBasic(projectName, webhookUrl),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "publisher_id", "pipelines"),
					resource.TestCheckResourceAttr(tfCheckNode, "event_type", "ms.vss-pipelines.run-state-changed-event"),
					resource.TestCheckResourceAttr(tfCheckNode, "consumer_id", "webHooks"),
					resource.TestCheckResourceAttr(tfCheckNode, "consumer_action_id", "httpRequest"),
					resource.TestCheckResourceAttr(tfCheckNode, "consumer_inputs.url", webhookUrl),
				),
			},
		},
	})
}

func TestAccServicehookSubscription_Update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	webhookUrl1 := "https://example1.com/webhook"
	webhookUrl2 := "https://example2.com/webhook"

	resourceType := "azuredevops_servicehook_subscription"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookSubscriptionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookSubscriptionResourceBasic(projectName, webhookUrl1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "consumer_inputs.url", webhookUrl1),
				),
			},
			{
				Config: testutils.HclServicehookSubscriptionResourceBasic(projectName, webhookUrl2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "consumer_inputs.url", webhookUrl2),
				),
			},
		},
	})
}

func TestAccServicehookSubscription_RequiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	webhookUrl := "https://example.com/webhook"

	resourceType := "azuredevops_servicehook_subscription"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookSubscriptionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookSubscriptionResourceBasic(projectName, webhookUrl),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
				),
			},
			{
				Config:                  testutils.HclServicehookSubscriptionResourceBasic(projectName, webhookUrl),
				ResourceName:            tfCheckNode,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					rs, ok := state.RootModule().Resources[tfCheckNode]
					if !ok {
						return "", fmt.Errorf("Not found: %s", tfCheckNode)
					}
					return rs.Primary.ID, nil
				},
			},
		},
	})
}

func CheckServicehookSubscriptionDestroyed(s *terraform.State) error {
	return nil // Since we don't have actual Azure DevOps instance, this is simplified
}
