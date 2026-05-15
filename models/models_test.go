package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var validMb ManifestByte = ManifestByte(`
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

	assert.Equal("ClusterRole_ClusterRoleName.yaml", cFileName, "they should be equal")
	assert.Equal("Role_RoleName_Namespace.yaml", nFileName, "they should be equal")

}

func TestIsValidManifest(t *testing.T) {

	assert := assert.New(t)

	var emptyMb ManifestByte
	assert.False(emptyMb.IsValidManifest())

	invalidMb := ManifestByte("---\n# Source: cni/templates/clusterrolebindings.yaml\n---")
	assert.False(invalidMb.IsValidManifest())

	assert.True(validMb.IsValidManifest())

}

func TestGetCmd(t *testing.T) {
	assert := assert.New(t)

	_, err := validMb.GetCmd("invalid")
	assert.Error(err)

	cmd, err := validMb.GetCmd("get")
	assert.NoError(err)
	assert.Equal("kubectl get -f out/ServiceAccount/ServiceAccount_istio-cni_kube-system.yaml -oyaml", cmd)

	cmd, err = validMb.GetCmd("diff")
	assert.NoError(err)
	assert.Equal("kubectl diff -f out/ServiceAccount/ServiceAccount_istio-cni_kube-system.yaml", cmd)

}
