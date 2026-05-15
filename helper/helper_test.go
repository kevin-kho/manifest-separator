package helper

import (
	"manifest-seperator/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mb []models.ManifestByte = []models.ManifestByte{
	models.ManifestByte(`
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
---`),
	models.ManifestByte(`
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
---`),
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
