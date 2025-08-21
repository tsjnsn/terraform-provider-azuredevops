terraform {
  required_providers {
    azuredevops = {
      source  = "microsoft/azuredevops"
      version = ">=0.9.0"
    }
  }
}

# Configure the Azure DevOps Provider
provider "azuredevops" {
  # Configuration can be set via environment variables:
  # AZDO_ORG_SERVICE_URL - Organization URL
  # AZDO_PERSONAL_ACCESS_TOKEN - Personal Access Token
}

# Create a project
resource "azuredevops_project" "example" {
  name               = "ServiceHook-Example"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Example project for Service Hook subscriptions"
}

# Create a webhook subscription for Git push events
resource "azuredevops_servicehook_subscription" "git_push_webhook" {
  project_id         = azuredevops_project.example.id
  publisher_id       = "tfs"
  event_type         = "git.push"
  consumer_id        = "webHooks"
  consumer_action_id = "httpRequest"

  publisher_inputs = {
    # Filter for pushes to main branch only
    branch = "refs/heads/main"
  }

  consumer_inputs = {
    # Replace with your webhook URL
    url = "https://webhook.example.com/git-push"
  }

  resource_version = "1.0"
  status          = "enabled"
}

# Create a webhook subscription for work item creation
resource "azuredevops_servicehook_subscription" "workitem_created_webhook" {
  project_id         = azuredevops_project.example.id
  publisher_id       = "tfs"
  event_type         = "workitem.created"
  consumer_id        = "webHooks"
  consumer_action_id = "httpRequest"

  publisher_inputs = {
    # Filter for Bug work items only
    workItemType = "Bug"
  }

  consumer_inputs = {
    # Replace with your webhook URL
    url = "https://webhook.example.com/workitem-created"
  }

  resource_version = "1.0"
  status          = "enabled"
}

# Create a webhook subscription for build completion (organization-level)
resource "azuredevops_servicehook_subscription" "build_complete_webhook" {
  # No project_id specified - organization-level subscription
  publisher_id       = "tfs"
  event_type         = "build.complete"
  consumer_id        = "webHooks"
  consumer_action_id = "httpRequest"

  publisher_inputs = {
    # Filter for successful builds only
    buildStatus = "Succeeded"
  }

  consumer_inputs = {
    # Replace with your webhook URL
    url = "https://webhook.example.com/build-complete"
  }

  resource_version = "1.0"
  status          = "enabled"
}

# Output the subscription IDs for reference
output "git_push_subscription_id" {
  value = azuredevops_servicehook_subscription.git_push_webhook.id
}

output "workitem_created_subscription_id" {
  value = azuredevops_servicehook_subscription.workitem_created_webhook.id
}

output "build_complete_subscription_id" {
  value = azuredevops_servicehook_subscription.build_complete_webhook.id
}