package scheme

import (
	"testing"

	buildV1 "github.com/openshift/api/build/v1"
	imageV1 "github.com/openshift/api/image/v1"
	routev1 "github.com/openshift/api/route/v1"
	secv1 "github.com/openshift/api/security/v1"
	monitoringV1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestAddToScheme(t *testing.T) {
	t.Parallel()

	s := runtime.NewScheme()

	err := AddToScheme(s)
	require.NoError(t, err)

	gv := []schema.GroupVersion{
		routev1.GroupVersion,
		secv1.GroupVersion,
		buildV1.GroupVersion,
		imageV1.GroupVersion,
		monitoringV1.SchemeGroupVersion,
	}

	for _, g := range gv {
		assert.True(t, s.IsVersionRegistered(g))
	}
}
