// Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validator

import (
	extensionspredicate "github.com/gardener/gardener/extensions/pkg/predicate"
	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/pkg/apis/core"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/gardener/gardener-extension-provider-aws/pkg/aws"
)

const (
	// Name is a name for a validation webhook.
	Name = "validator"
	// SecretsValidatorName is the name of the secrets validator.
	SecretsValidatorName = "secrets." + Name
)

var logger = log.Log.WithName("aws-validator-webhook")

// New creates a new webhook that validates Shoot and CloudProfile resources.
func New(mgr manager.Manager) (*extensionswebhook.Webhook, error) {
	logger.Info("Setting up webhook", "name", Name)

	return extensionswebhook.New(mgr, extensionswebhook.Args{
		Provider:   aws.Type,
		Name:       Name,
		Path:       "/webhooks/validate",
		Predicates: []predicate.Predicate{extensionspredicate.GardenCoreProviderType(aws.Type)},
		Validators: map[extensionswebhook.Validator][]extensionswebhook.Type{
			NewShootValidator(mgr):         {{Obj: &core.Shoot{}}},
			NewCloudProfileValidator(mgr):  {{Obj: &core.CloudProfile{}}},
			NewSecretBindingValidator(mgr): {{Obj: &core.SecretBinding{}}},
		},
		Target: extensionswebhook.TargetSeed,
		ObjectSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"provider.extensions.gardener.cloud/aws": "true"},
		},
	})
}

// NewSecretsWebhook creates a new validation webhook for Secrets.
func NewSecretsWebhook(mgr manager.Manager) (*extensionswebhook.Webhook, error) {
	logger.Info("Setting up webhook", "name", SecretsValidatorName)

	return extensionswebhook.New(mgr, extensionswebhook.Args{
		Provider: aws.Type,
		Name:     SecretsValidatorName,
		Path:     "/webhooks/validate/secrets",
		Validators: map[extensionswebhook.Validator][]extensionswebhook.Type{
			NewSecretValidator(): {{Obj: &corev1.Secret{}}},
		},
		Target: extensionswebhook.TargetSeed,
		ObjectSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"provider.shoot.gardener.cloud/aws": "true"},
		},
	})
}
