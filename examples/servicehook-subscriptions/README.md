# Service Hook Subscriptions Example

This example demonstrates how to use the `azuredevops_servicehook_subscription` resource to create various types of service hook subscriptions in Azure DevOps.

## Usage

1. Set up your Azure DevOps credentials:
   ```bash
   export AZDO_ORG_SERVICE_URL="https://dev.azure.com/yourorgname"
   export AZDO_PERSONAL_ACCESS_TOKEN="your-pat-here"
   ```

2. Update the webhook URLs in `main.tf` to point to your actual webhook endpoints.

3. Run Terraform:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## What This Creates

This example creates:

1. **Azure DevOps Project**: A new project to host the service hooks
2. **Git Push Webhook**: Triggers on pushes to the main branch
3. **Work Item Created Webhook**: Triggers when new Bug work items are created
4. **Build Complete Webhook**: Organization-level webhook for successful builds

## Service Hook Types

The `azuredevops_servicehook_subscription` resource supports many different combinations of publishers, events, and consumers:

### Common Publishers
- `tfs` - Team Foundation Server (Git, Work Items, Builds, etc.)
- `pipelines` - Azure Pipelines
- `boards` - Azure Boards

### Common Event Types
- `git.push` - Git repository push
- `git.pullrequest.created` - Pull request created
- `build.complete` - Build completed
- `workitem.created` - Work item created
- `workitem.updated` - Work item updated
- `ms.vss-pipelines.run-state-changed-event` - Pipeline run state changed

### Common Consumers
- `webHooks` - HTTP webhooks
- `azureServiceBus` - Azure Service Bus
- `azureStorageQueue` - Azure Storage Queue

## Cleanup

To remove all resources:

```bash
terraform destroy
```