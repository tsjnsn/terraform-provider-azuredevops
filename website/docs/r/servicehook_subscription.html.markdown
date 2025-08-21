---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_servicehook_subscription"
description: |-
  Manages a Service Hook Subscription.
---

# azuredevops_servicehook_subscription

Manages a Service Hook Subscription in Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

# Webhook subscription for pipeline events
resource "azuredevops_servicehook_subscription" "pipeline_webhook" {
  project_id         = azuredevops_project.example.id
  publisher_id       = "pipelines"
  event_type         = "ms.vss-pipelines.run-state-changed-event"
  consumer_id        = "webHooks"
  consumer_action_id = "httpRequest"
  consumer_inputs = {
    url = "https://example.com/webhook"
  }
  publisher_inputs = {
    runStateId = "Completed"
  }
}

# Azure Storage Queue subscription for build events  
resource "azuredevops_servicehook_subscription" "build_storage_queue" {
  project_id         = azuredevops_project.example.id
  publisher_id       = "tfs"
  event_type         = "build.complete"
  consumer_id        = "azureStorageQueue"
  consumer_action_id = "enqueue"
  consumer_inputs = {
    accountName = "mystorageaccount"
    accountKey  = "myaccountkey"
    queueName   = "myqueue"
  }
}

# Organization-level subscription (no project_id)
resource "azuredevops_servicehook_subscription" "org_webhook" {
  publisher_id       = "tfs"
  event_type         = "workitem.created"
  consumer_id        = "webHooks"
  consumer_action_id = "httpRequest"
  consumer_inputs = {
    url = "https://example.com/org-webhook"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `publisher_id` - (Required) The ID of the publisher (e.g., 'pipelines', 'tfs', 'git'). Changing this forces a new resource to be created.

* `event_type` - (Required) The event type identifier (e.g., 'ms.vss-pipelines.run-state-changed-event', 'build.complete'). Changing this forces a new resource to be created.

* `consumer_id` - (Required) The ID of the consumer (e.g., 'webHooks', 'azureStorageQueue'). Changing this forces a new resource to be created.

* `consumer_action_id` - (Required) The consumer action identifier (e.g., 'httpRequest', 'enqueue'). Changing this forces a new resource to be created.

* `consumer_inputs` - (Required) A map of consumer-specific input values.

---

* `project_id` - (Optional) The ID of the project. If not specified, creates an organization-level subscription.

* `publisher_inputs` - (Optional) A map of publisher-specific input values.

* `resource_version` - (Optional) The resource version for the subscription. Defaults to `5.1-preview.1`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Hook Subscription.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Service Hook Subscription.
* `read` - (Defaults to 5 minute) Used when retrieving the Service Hook Subscription.
* `update` - (Defaults to 10 minutes) Used when updating the Service Hook Subscription.
* `delete` - (Defaults to 10 minutes) Used when deleting the Service Hook Subscription.

## Import

Service Hook Subscriptions can be imported using the subscription ID:

```sh
terraform import azuredevops_servicehook_subscription.example 00000000-0000-0000-0000-000000000000
```

## Relevant Links

* [Azure DevOps Service REST API 7.0 - Service Hooks](https://docs.microsoft.com/en-us/rest/api/azure/devops/hooks/?view=azure-devops-rest-7.0)