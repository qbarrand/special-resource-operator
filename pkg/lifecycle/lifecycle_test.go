package lifecycle

import (
	"context"
	"os"
	"testing"

	"github.com/openshift-psap/special-resource-operator/pkg/clients"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	name      = "test"
	namespace = "ns"
)

var labels = map[string]string{"key": "value"}

func makePod(name string) *v1.Pod {
	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
	}
}

func resetGlobals() {
	clients.Interface = nil
}

func TestGetPodFromDaemonSet(t *testing.T) {
	nsn := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}

	t.Run("DaemonSet does not exist", func(t *testing.T) {
		t.Cleanup(resetGlobals)

		clients.Interface = &clients.ClientsInterface{
			Client: fake.NewClientBuilder().Build(),
		}

		pl := GetPodFromDaemonSet(nsn)

		assert.Empty(t, pl.Items)
	})

	t.Run("DaemonSet has 2 pods", func(t *testing.T) {
		t.Cleanup(resetGlobals)

		ds := appsv1.DaemonSet{
			TypeMeta: metav1.TypeMeta{Kind: "DaemonSet"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: appsv1.DaemonSetSpec{
				Selector: &metav1.LabelSelector{MatchLabels: labels},
			},
		}

		clients.Interface = &clients.ClientsInterface{
			Client: fake.NewClientBuilder().WithObjects(&ds, makePod("p1"), makePod("p2")).Build(),
		}

		pl := GetPodFromDaemonSet(types.NamespacedName{Namespace: "ns", Name: "test"})

		assert.Len(t, pl.Items, 2)
	})
}

func TestUpdateDaemonSetPods(t *testing.T) {
	t.Cleanup(resetGlobals)

	const namespaceEnvVar = "OPERATOR_NAMESPACE"

	err := os.Setenv(namespaceEnvVar, namespace)
	require.NoError(t, err)

	t.Cleanup(func() {
		// TODO(qbarrand) when we move to Go 1.17, use https://pkg.go.dev/testing#T.Setenv
		err = os.Unsetenv(namespaceEnvVar)
		assert.NoError(t, err)
	})

	ds := appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{Kind: "DaemonSet"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: labels},
		},
	}

	const cmName = "special-resource-lifecycle"

	cm := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmName,
			Namespace: namespace,
		},
	}

	kubeClient := fake.NewClientBuilder().WithObjects(
		&ds,
		&cm,
		makePod("p1"),
		makePod("p2"),
	).
		Build()

	clients.Interface = &clients.ClientsInterface{Client: kubeClient}

	err = UpdateDaemonSetPods(&ds)
	require.NoError(t, err)

	err = kubeClient.Get(context.TODO(), types.NamespacedName{Name: cmName, Namespace: namespace}, &cm)
	require.NoError(t, err)

	assert.Equal(t, "*v1.Pod", cm.Data["49cebbba48498f09"])
	assert.Equal(t, "*v1.Pod", cm.Data["49ceb8ba484989f0"])
}
