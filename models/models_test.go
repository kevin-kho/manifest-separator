package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileName(t *testing.T) {
	assert := assert.New(t)
	cm := Manifest{
		ApiVersion: "v1",
		Kind:       "ClusterRole",
		Metadata: Metadata{
			Name:      "ClusterRoleName",
			Namespace: "",
		},
	}

	nm := Manifest{
		ApiVersion: "v1",
		Kind:       "Role",
		Metadata: Metadata{
			Name:      "RoleName",
			Namespace: "Namespace",
		},
	}

	cFileName := cm.GetFileName()
	nFileName := nm.GetFileName()

	assert.Equal(cFileName, "ClusterRole_ClusterRoleName.yaml", "they should be equal")
	assert.Equal(nFileName, "Role_RoleName_Namespace.yaml", "they should be equal")

}

func TestIsValidManifest(t *testing.T) {

	assert := assert.New(t)

	var emptyMb ManifestByte
	assert.False(emptyMb.IsValidManifest())

	invalidMb := ManifestByte("---\n# Source: cni/templates/clusterrolebindings.yaml\n---")
	assert.False(invalidMb.IsValidManifest())

	validMb := ManifestByte(`
---
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
	`)

	assert.True(validMb.IsValidManifest())

}

func TestGetCmd(t *testing.T) {
	assert := assert.New(t)

	mb := ManifestByte(`
---
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
	`)

	_, err := mb.GetCmd("invalid")
	assert.Error(err)

	cmd, err := mb.GetCmd("get")
	assert.NoError(err)
	assert.Equal("kubectl get -f out/ServiceAccount/ServiceAccount_istio-cni_kube-system.yaml -oyaml", cmd)

	cmd, err = mb.GetCmd("diff")
	assert.NoError(err)
	assert.Equal("kubectl diff -f out/ServiceAccount/ServiceAccount_istio-cni_kube-system.yaml", cmd)

}
