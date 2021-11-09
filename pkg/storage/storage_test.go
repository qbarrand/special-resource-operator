package storage

import (
	"context"
	"testing"

	"github.com/openshift-psap/special-resource-operator/pkg/clients"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	namespaceName = "test-ns"
	resourceName  = "test-resource"
)

var nsn = types.NamespacedName{Namespace: namespaceName, Name: resourceName}

func getConfigmap(t *testing.T, client client.Client) *v1.ConfigMap {
	t.Helper()

	cm := v1.ConfigMap{}

	err := client.Get(context.TODO(), nsn, &cm)
	require.NoError(t, err)

	return &cm
}

func resetGlobals() {
	clients.Interface = nil
}

func TestCheckConfigMapEntry(t *testing.T) {
	const key = "test-key"

	t.Run("no ConfigMap", func(t *testing.T) {
		t.Cleanup(resetGlobals)

		clients.Interface = &clients.ClientsInterface{Client: fake.NewClientBuilder().Build()}

		_, err := CheckConfigMapEntry(key, nsn)
		assert.Error(t, err)
	})

	t.Run("ConfigMap with key absent", func(t *testing.T) {
		t.Cleanup(resetGlobals)

		cm := v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: resourceName,
			},
			Data: make(map[string]string),
		}

		clients.Interface = &clients.ClientsInterface{
			Client: fake.NewClientBuilder().WithObjects(&cm).Build(),
		}

		_, err := CheckConfigMapEntry(key, nsn)
		assert.Error(t, err)
	})

	t.Run("ConfigMap with key present", func(t *testing.T) {
		t.Cleanup(resetGlobals)

		const data = "test-data"

		cm := v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: namespaceName,
			},
			Data: map[string]string{key: data},
		}

		clients.Interface = &clients.ClientsInterface{
			Client: fake.NewClientBuilder().WithObjects(&cm).Build(),
		}

		v, err := CheckConfigMapEntry(key, nsn)
		require.NoError(t, err)

		assert.Equal(t, data, v)
	})
}

func TestGetConfigMap(t *testing.T) {
	t.Run("no ConfigMap", func(t *testing.T) {
		t.Cleanup(resetGlobals)

		clients.Interface = &clients.ClientsInterface{Client: fake.NewClientBuilder().Build()}

		_, err := GetConfigMap(namespaceName, resourceName)
		assert.Error(t, err)
	})

	t.Run("ConfigMap present", func(t *testing.T) {
		t.Cleanup(resetGlobals)

		data := map[string]string{"key": "value"}

		cm := v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: namespaceName,
			},
			Data: data,
		}

		clients.Interface = &clients.ClientsInterface{
			Client: fake.NewClientBuilder().WithObjects(&cm).Build(),
		}

		res, err := GetConfigMap(namespaceName, resourceName)
		require.NoError(t, err)

		resData, found, err := unstructured.NestedStringMap(res.Object, "data")
		require.NoError(t, err)

		assert.True(t, found)
		assert.Equal(t, data, resData)
	})
}

func TestUpdateConfigMapEntry(t *testing.T) {
	t.Run("no ConfigMap", func(t *testing.T) {
		t.Cleanup(resetGlobals)

		clients.Interface = &clients.ClientsInterface{Client: fake.NewClientBuilder().Build()}

		err := UpdateConfigMapEntry("any-key", "any-value", nsn)
		assert.Error(t, err)
	})

	t.Run("key does not exist", func(t *testing.T) {
		t.Cleanup(resetGlobals)

		const (
			key   = "key"
			value = "value"
		)

		cm := v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: namespaceName,
			},
			Data: make(map[string]string),
		}

		kubeClient := fake.NewClientBuilder().WithObjects(&cm).Build()

		clients.Interface = &clients.ClientsInterface{Client: kubeClient}

		err := UpdateConfigMapEntry(key, value, nsn)
		assert.NoError(t, err)

		// check that the data was indeed updated wby querying Kubernetes
		assert.Equal(t, value, getConfigmap(t, kubeClient).Data[key])
	})

	t.Run("key exists", func(t *testing.T) {
		t.Cleanup(resetGlobals)

		const (
			key      = "key"
			newValue = "new-value"
		)

		cm := v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: namespaceName,
			},
			Data: map[string]string{key: "old-value"},
		}

		kubeClient := fake.NewClientBuilder().WithObjects(&cm).Build()

		clients.Interface = &clients.ClientsInterface{Client: kubeClient}

		err := UpdateConfigMapEntry(key, newValue, nsn)
		assert.NoError(t, err)

		// check that the data was indeed updated wby querying Kubernetes
		assert.Equal(t, newValue, getConfigmap(t, kubeClient).Data[key])
	})
}

func TestDeleteConfigMapEntry(t *testing.T) {
	t.Run("no ConfigMap", func(t *testing.T) {
		t.Cleanup(resetGlobals)

		clients.Interface = &clients.ClientsInterface{Client: fake.NewClientBuilder().Build()}

		err := DeleteConfigMapEntry("any-key", nsn)
		assert.Error(t, err)
	})

	t.Run("key does not exist", func(t *testing.T) {
		t.Cleanup(resetGlobals)

		data := map[string]string{"key": "value"}

		cm := v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: namespaceName,
			},
			Data: data,
		}

		kubeClient := fake.NewClientBuilder().WithObjects(&cm).Build()

		clients.Interface = &clients.ClientsInterface{Client: kubeClient}

		err := DeleteConfigMapEntry("some-other-key", nsn)
		assert.NoError(t, err)

		// check that the data is still the same Kubernetes
		assert.Equal(t, data, getConfigmap(t, kubeClient).Data)
	})

	t.Run("key exists", func(t *testing.T) {
		t.Cleanup(resetGlobals)

		const (
			key      = "key"
			otherKey = "other-key"
			value    = "value"
		)

		data := map[string]string{key: value, otherKey: "other-value"}

		cm := v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: namespaceName,
			},
			Data: data,
		}

		kubeClient := fake.NewClientBuilder().WithObjects(&cm).Build()

		clients.Interface = &clients.ClientsInterface{Client: kubeClient}

		err := DeleteConfigMapEntry(otherKey, nsn)
		assert.NoError(t, err)

		// check that the data was indeed deleted wby querying Kubernetes
		assert.Equal(t, map[string]string{key: value}, getConfigmap(t, kubeClient).Data)
	})
}
