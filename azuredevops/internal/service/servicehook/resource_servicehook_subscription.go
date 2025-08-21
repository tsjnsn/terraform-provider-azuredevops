package servicehook

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/servicehooks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceServicehookSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceServicehookSubscriptionCreate,
		Read:   resourceServicehookSubscriptionRead,
		Update: resourceServicehookSubscriptionUpdate,
		Delete: resourceServicehookSubscriptionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The ID of the project. Leave empty for organization-level subscriptions.",
			},
			"publisher_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the publisher (e.g., 'pipelines', 'git', 'workitems').",
			},
			"event_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The event type (e.g., 'ms.vss-pipelines.run-state-changed-event').",
			},
			"consumer_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the consumer (e.g., 'webHooks', 'azureStorageQueue').",
			},
			"consumer_action_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The consumer action ID (e.g., 'httpRequest', 'enqueue').",
			},
			"publisher_inputs": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Publisher-specific input values.",
			},
			"consumer_inputs": {
				Type:        schema.TypeMap,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Consumer-specific input values.",
			},
			"resource_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "5.1-preview.1",
				Description: "The resource version for the subscription.",
			},
		},
	}
}

func resourceServicehookSubscriptionCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscription := expandServicehookSubscription(d)

	createdSubscription, err := createSubscription(clients, subscription)
	if err != nil {
		return err
	}

	d.SetId(createdSubscription.Id.String())
	return resourceServicehookSubscriptionRead(d, m)
}

func resourceServicehookSubscriptionRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscriptionId := converter.UUID(d.Id())

	subscription, err := getSubscription(clients, subscriptionId)
	if err != nil {
		return err
	}

	flattenServicehookSubscription(d, subscription)
	return nil
}

func resourceServicehookSubscriptionUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscription := expandServicehookSubscription(d)
	subscriptionId := converter.UUID(d.Id())
	subscription.Id = subscriptionId

	_, err := updateSubscription(clients, subscription)
	if err != nil {
		return err
	}

	return resourceServicehookSubscriptionRead(d, m)
}

func resourceServicehookSubscriptionDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	subscriptionId := converter.UUID(d.Id())

	return clients.ServiceHooksClient.DeleteSubscription(
		clients.Ctx,
		servicehooks.DeleteSubscriptionArgs{
			SubscriptionId: subscriptionId,
		})
}

func expandServicehookSubscription(d *schema.ResourceData) *servicehooks.Subscription {
	publisherInputs := make(map[string]string)

	// Add project ID to publisher inputs if specified
	if projectId := d.Get("project_id").(string); projectId != "" {
		publisherInputs["projectId"] = projectId
	}

	// Add additional publisher inputs
	if inputs, ok := d.GetOk("publisher_inputs"); ok {
		for key, value := range inputs.(map[string]interface{}) {
			publisherInputs[key] = value.(string)
		}
	}

	consumerInputs := make(map[string]string)
	for key, value := range d.Get("consumer_inputs").(map[string]interface{}) {
		consumerInputs[key] = value.(string)
	}

	return &servicehooks.Subscription{
		PublisherId:      converter.String(d.Get("publisher_id").(string)),
		EventType:        converter.String(d.Get("event_type").(string)),
		ConsumerId:       converter.String(d.Get("consumer_id").(string)),
		ConsumerActionId: converter.String(d.Get("consumer_action_id").(string)),
		PublisherInputs:  &publisherInputs,
		ConsumerInputs:   &consumerInputs,
		ResourceVersion:  converter.String(d.Get("resource_version").(string)),
	}
}

func flattenServicehookSubscription(d *schema.ResourceData, subscription *servicehooks.Subscription) {
	d.Set("publisher_id", *subscription.PublisherId)
	d.Set("event_type", *subscription.EventType)
	d.Set("consumer_id", *subscription.ConsumerId)
	d.Set("consumer_action_id", *subscription.ConsumerActionId)
	d.Set("resource_version", *subscription.ResourceVersion)

	// Extract project_id from publisher inputs if present
	if subscription.PublisherInputs != nil {
		publisherInputs := make(map[string]interface{})
		for key, value := range *subscription.PublisherInputs {
			if key == "projectId" {
				d.Set("project_id", value)
			} else {
				publisherInputs[key] = value
			}
		}
		d.Set("publisher_inputs", publisherInputs)
	}

	// Set consumer inputs
	if subscription.ConsumerInputs != nil {
		consumerInputs := make(map[string]interface{})
		for key, value := range *subscription.ConsumerInputs {
			consumerInputs[key] = value
		}
		d.Set("consumer_inputs", consumerInputs)
	}
}
