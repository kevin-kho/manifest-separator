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
