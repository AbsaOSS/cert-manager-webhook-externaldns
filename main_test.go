package main

import (
	"context"
	"os"
	"testing"

	whapi "github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/stretchr/testify/assert"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/external-dns/endpoint"
)

func TestRunsSuite(t *testing.T) {
	// The manifest path should contain a file named config.json that is a
	// snippet of valid configuration that should be included on the
	// ChallengeRequest passed as part of the test cases.
	jsonData, err := os.ReadFile("testdata/externaldns/config.json")
	expectedLabels := map[string]string{"cert-manager": "true", "some-label": "bar"}
	assert.NoError(t, err)

	var ep endpoint.DNSEndpoint
	challenge := &whapi.ChallengeRequest{
		ResourceNamespace: "default",
		DNSName:           "example.com.",
		ResolvedFQDN:      "_acme-challenge.example.com",
		Key:               "foo",
		Config: &extapi.JSON{
			Raw: jsonData,
		},
	}
	objName := fqdnToName(challenge.ResolvedFQDN)
	namespace := challenge.ResourceNamespace
	solver := &externalDNSProviderSolver{}

	stopCh := make(chan struct{})
	crdInstallOpts := envtest.CRDInstallOptions{
		Paths:              []string{"./testdata/crd/crd-manifest.yaml"},
		ErrorIfPathMissing: true,
	}
	envTest := &envtest.Environment{
		CRDInstallOptions: crdInstallOpts,
	}
	conf, err := envTest.Start()
	assert := assert.New(t)
	assert.NoError(err)
	solver.Initialize(conf, stopCh)
	err = solver.Present(challenge)
	assert.NoError(err)
	obj, err := solver.client.Resource(gvr).Namespace(namespace).Get(context.TODO(), objName, metav1.GetOptions{})
	assert.NoError(err)
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &ep)
	assert.NoError(err)
	assert.Equal(ep.Spec.Endpoints[0].Targets[0], challenge.Key)
	assert.Equal(ep.Spec.Endpoints[0].DNSName, challenge.ResolvedFQDN)
	assert.Equal(ep.Spec.Endpoints[0].RecordType, "TXT")
	assert.Equal(ep.ObjectMeta.Labels, expectedLabels)

	err = solver.CleanUp(challenge)
	assert.NoError(err)
	_ = envTest.Stop()

}
