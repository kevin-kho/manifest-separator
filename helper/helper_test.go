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

var mb []models.ManifestByte = []models.ManifestByte{
	models.ManifestByte(`---
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
---`),
	models.ManifestByte(`---
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
---`),
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

	mp, err := GetKinds(mb)
	assert.Nil(err)
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

func TestContainDupes(t *testing.T) {
	assert := assert.New(t)

	res := ContainsDupes(duplicateManifests)
	assert.True(res)

	res = ContainsDupes(uniqueManifests)
	assert.False(res)

}
