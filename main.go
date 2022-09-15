package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"

	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/external-dns/endpoint"
)

const GroupName = "externaldns.k8s.io"
const GroupVersion = "v1alpha1"

var gvk = schema.GroupVersionKind{
	Group:   GroupName,
	Version: GroupVersion,
	Kind:    "DNSEndpoint",
}
var gvr = schema.GroupVersionResource{
	Group:    GroupName,
	Version:  GroupVersion,
	Resource: "dnsendpoints",
}

type externalDNSProviderSolver struct {
	client dynamic.Interface
}

type externalDNSProviderConfig struct {
	Labels map[string]string `json:"labels"`
}

func (c *externalDNSProviderSolver) Name() string {
	return "externaldns"
}

func main() {
	fmt.Println("Starting webhook for externaldns")
	cmd.RunWebhookServer(GroupName,
		&externalDNSProviderSolver{},
	)
}
func fqdnToName(fqdn string) string {
	fqdn = strings.TrimSuffix(fqdn, ".")
	fqdn = strings.TrimPrefix(fqdn, "_")
	return strings.ReplaceAll(fqdn, ".", "-")
}

func (c *externalDNSProviderSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return err
	}

	obj := unstructured.Unstructured{}
	fqdn := ch.ResolvedFQDN

	endpointSpec := endpoint.DNSEndpointSpec{}

	endpointSpec.Endpoints = append(endpointSpec.Endpoints, endpoint.NewEndpoint(fqdn, endpoint.RecordTypeTXT, ch.Key))
	obj.SetName(fqdnToName(fqdn))
	obj.SetGroupVersionKind(gvk)
	obj.SetLabels(cfg.Labels)

	epSpec, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&endpointSpec)
	if err != nil {
		return err
	}
	if err := unstructured.SetNestedField(obj.Object, epSpec, "spec"); err != nil {
		return err
	}
	if _, err := c.client.Resource(gvr).Namespace(ch.ResourceNamespace).Create(context.TODO(), &obj, metav1.CreateOptions{}); err != nil {
		return err
	}
	return nil
}

// CleanUp should delete the relevant TXT record from the DNS provider console.
func (c *externalDNSProviderSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	if err := c.client.Resource(gvr).Namespace(ch.ResourceNamespace).Delete(context.TODO(), fqdnToName(ch.ResolvedFQDN), deleteOptions); err != nil {
		return err
	}
	return nil

}

// Initialize will be called when the webhook first starts.
func (c *externalDNSProviderSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	cl, err := dynamic.NewForConfig(kubeClientConfig)
	if err != nil {
		return err
	}

	c.client = cl

	return nil
}

// loadConfig is a small helper function that decodes JSON configuration into
// the typed config struct.
func loadConfig(cfgJSON *extapi.JSON) (*externalDNSProviderConfig, error) {
	cfg := externalDNSProviderConfig{}
	// handle the 'base case' where no configuration has been provided
	if cfgJSON == nil {
		return &cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return &cfg, fmt.Errorf("error decoding solver config: %v", err)
	}

	return &cfg, nil
}
