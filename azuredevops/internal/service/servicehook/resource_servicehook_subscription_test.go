//go:build (all || resource_servicehook_subscription) && !exclude_subscriptions
// +build all resource_servicehook_subscription
// +build !exclude_subscriptions

package servicehook

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/servicehooks"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var testResourceSubscription = []servicehooks.Subscription{
	{
		Id:               &uuid.UUID{},
		ConsumerActionId: converter.String("httpRequest"),
		ConsumerId:       converter.String("webHooks"),
		ConsumerInputs: &map[string]string{
			"url": "https://example.com/webhook",
		},
		EventType:   converter.String("ms.vss-pipelines.run-state-changed-event"),
		PublisherId: converter.String("pipelines"),
		PublisherInputs: &map[string]string{
			"projectId": "myprojectid",
		},
		ResourceVersion: converter.String("5.1-preview.1"),
	},
}

func TestServicehookSubscription_FlattenExpandRoundTrip(t *testing.T) {
	for _, subscription := range testResourceSubscription {
		resourceData := schema.TestResourceDataRaw(t, ResourceServicehookSubscription().Schema, nil)
		flattenServicehookSubscription(resourceData, &subscription)
		subscriptionAfterRoundTrip := expandServicehookSubscription(resourceData)
		subscriptionAfterRoundTrip.Id = subscription.Id

		require.Equal(t, subscription, *subscriptionAfterRoundTrip)
	}
}

func TestServicehookSubscription_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServicehookSubscription()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.Set("publisher_id", "pipelines")
	resourceData.Set("event_type", "ms.vss-pipelines.run-state-changed-event")
	resourceData.Set("consumer_id", "webHooks")
	resourceData.Set("consumer_action_id", "httpRequest")
	resourceData.Set("consumer_inputs", map[string]interface{}{
		"url": "https://example.com/webhook",
	})

	buildClient := azdosdkmocks.NewMockServicehooksClient(ctrl)
	clients := &client.AggregatedClient{ServiceHooksClient: buildClient, Ctx: context.Background()}

	expectedArgs := servicehooks.CreateSubscriptionArgs{Subscription: &servicehooks.Subscription{
		PublisherId:      converter.String("pipelines"),
		EventType:        converter.String("ms.vss-pipelines.run-state-changed-event"),
		ConsumerId:       converter.String("webHooks"),
		ConsumerActionId: converter.String("httpRequest"),
		PublisherInputs:  &map[string]string{},
		ConsumerInputs: &map[string]string{
			"url": "https://example.com/webhook",
		},
		ResourceVersion: converter.String("5.1-preview.1"),
	}}

	buildClient.
		EXPECT().
		CreateSubscription(clients.Ctx, expectedArgs).
		Return(nil, errors.New("CreateSubscription() Failed")).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateSubscription() Failed")
}

func TestServicehookSubscription_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServicehookSubscription()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.SetId("00000000-0000-0000-0000-000000000000")
	resourceData.Set("publisher_id", "pipelines")
	resourceData.Set("event_type", "ms.vss-pipelines.run-state-changed-event")
	resourceData.Set("consumer_id", "webHooks")
	resourceData.Set("consumer_action_id", "httpRequest")
	resourceData.Set("consumer_inputs", map[string]interface{}{
		"url": "https://example.com/webhook",
	})

	buildClient := azdosdkmocks.NewMockServicehooksClient(ctrl)
	clients := &client.AggregatedClient{ServiceHooksClient: buildClient, Ctx: context.Background()}

	expectedArgs := servicehooks.ReplaceSubscriptionArgs{
		SubscriptionId: converter.UUID("00000000-0000-0000-0000-000000000000"),
		Subscription: &servicehooks.Subscription{
			Id:               converter.UUID("00000000-0000-0000-0000-000000000000"),
			PublisherId:      converter.String("pipelines"),
			EventType:        converter.String("ms.vss-pipelines.run-state-changed-event"),
			ConsumerId:       converter.String("webHooks"),
			ConsumerActionId: converter.String("httpRequest"),
			PublisherInputs:  &map[string]string{},
			ConsumerInputs: &map[string]string{
				"url": "https://example.com/webhook",
			},
			ResourceVersion: converter.String("5.1-preview.1"),
		},
	}

	buildClient.
		EXPECT().
		ReplaceSubscription(clients.Ctx, expectedArgs).
		Return(nil, errors.New("ReplaceSubscription() Failed")).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "ReplaceSubscription() Failed")
}

func TestServicehookSubscription_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServicehookSubscription()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.SetId("00000000-0000-0000-0000-000000000000")

	buildClient := azdosdkmocks.NewMockServicehooksClient(ctrl)
	clients := &client.AggregatedClient{ServiceHooksClient: buildClient, Ctx: context.Background()}

	expectedArgs := servicehooks.GetSubscriptionArgs{
		SubscriptionId: converter.UUID("00000000-0000-0000-0000-000000000000"),
	}

	buildClient.
		EXPECT().
		GetSubscription(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetSubscription() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "GetSubscription() Failed")
}

func TestServicehookSubscription_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServicehookSubscription()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	resourceData.SetId("00000000-0000-0000-0000-000000000000")

	buildClient := azdosdkmocks.NewMockServicehooksClient(ctrl)
	clients := &client.AggregatedClient{ServiceHooksClient: buildClient, Ctx: context.Background()}

	expectedArgs := servicehooks.DeleteSubscriptionArgs{
		SubscriptionId: converter.UUID("00000000-0000-0000-0000-000000000000"),
	}

	buildClient.
		EXPECT().
		DeleteSubscription(clients.Ctx, expectedArgs).
		Return(errors.New("DeleteSubscription() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteSubscription() Failed")
}
