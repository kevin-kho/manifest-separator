package helper

import (
	"manifest-seperator/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var data []byte = []byte(`---
# Source: cni/templates/serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: istio-cni
  namespace: kube-system
  labels:
    app: istio-cni
    release: istio-cni
    istio.io/rev: default
    install.operator.istio.io/owning-resource: unknown
    operator.istio.io/component: "Cni"
---
# Source: cni/templates/serviceaccount.yaml
apiVersion: v1
kind: Role
metadata:
  name: roleName
  namespace: namespace
  labels:
    app: istio-cni
    release: istio-cni
    istio.io/rev: default
    install.operator.istio.io/owning-resource: unknown
    operator.istio.io/component: "Cni"
---`)

var appSetData []byte = []byte(`- apiVersion: argoproj.io/v1alpha1
  kind: Application
  metadata:
    finalizers:
    - resources-finalizer.argocd.argoproj.io
    name: app-name-0
  spec:
    destination:
      name: destination-name-0
      namespace: namespace-0
    project: default
    sources:
    - path: path-0
      repoURL: repoURL
      targetRevision: HEAD
    syncPolicy:
      automated:
        prune: true
        selfHeal: true
      syncOptions:
      - ServerSideApply=true
- apiVersion: argoproj.io/v1alpha1
  kind: Application
  metadata:
    finalizers:
    - resources-finalizer.argocd.argoproj.io
    name: app-name-1
  spec:
    destination:
      name: destination-name-1
      namespace: namespace-1
    project: default
    sources:
    - path: path
    - repoURL: repoUrl
      targetRevision: HEAD
    syncPolicy:
      automated:
        prune: true
        selfHeal: true
      syncOptions:
      - ServerSideApply=true`)

var mb []models.Manifest = []models.Manifest{
	{
		ApiVersion: "",
		Kind:       "ServiceAccount",
		Metadata:   models.Metadata{},
	},
	{
		ApiVersion: "",
		Kind:       "Role",
		Metadata:   models.Metadata{},
	},
}

var duplicateManifests []models.Manifest = []models.Manifest{
	{
		ApiVersion: "v1",
		Kind:       "ServiceAccount",
		Metadata: models.Metadata{
			Name:      "name",
			Namespace: "namespace",
		},
	},
	{
		ApiVersion: "v1",
		Kind:       "ServiceAccount",
		Metadata: models.Metadata{
			Name:      "name",
			Namespace: "namespace",
		},
	},
}

var uniqueManifests []models.Manifest = []models.Manifest{
	{
		ApiVersion: "v1",
		Kind:       "ServiceAccount",
		Metadata: models.Metadata{
			Name:      "alice",
			Namespace: "namespace",
		},
	},
	{
		ApiVersion: "v1",
		Kind:       "ServiceAccount",
		Metadata: models.Metadata{
			Name:      "bob",
			Namespace: "namespace",
		},
	},
}

func TestGetKinds(t *testing.T) {
	assert := assert.New(t)

	mp := GetKinds(mb)
	assert.Len(mp, 2)
	assert.True(mp["Role"])
	assert.True(mp["ServiceAccount"])
	assert.False(mp["RoleBinding"])

}

func TestSeparateManifests(t *testing.T) {
	assert := assert.New(t)
	m, err := SeparateManifests(data)
	assert.Nil(err)

	assert.Len(m, 2)

}

func TestSeparateAppSet(t *testing.T) {
	assert := assert.New(t)

	m, err := SeparateAppSet(appSetData)

	assert.Nil(err)
	assert.Len(m, 2)
}

func TestContainDupes(t *testing.T) {
	assert := assert.New(t)

	res := ContainsDupes(duplicateManifests)
	assert.True(res)

	res = ContainsDupes(uniqueManifests)
	assert.False(res)

}
