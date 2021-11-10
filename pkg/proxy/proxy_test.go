package proxy

import (
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func resetProxy() {
	ProxyConfiguration = Configuration{}
}

func TestSetup(t *testing.T) {
	t.Run("Pod with empty spec", func(t *testing.T) {
		pod := v1.Pod{
			TypeMeta: metav1.TypeMeta{Kind: "Pod"},
		}

		m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&pod)
		require.NoError(t, err)

		uo := unstructured.Unstructured{Object: m}

		err = Setup(&uo)
		require.Error(t, err)
	})

	t.Run("Pod with one container", func(t *testing.T) {
		t.Cleanup(resetProxy)

		const (
			httpProxy  = "http-host-with-proxy"
			httpsProxy = "https-host-with-proxy"
			noProxy    = "host-without-proxy"
		)

		ProxyConfiguration = Configuration{
			HttpProxy:  httpProxy,
			HttpsProxy: httpsProxy,
			NoProxy:    noProxy,
		}

		pod := v1.Pod{
			TypeMeta: metav1.TypeMeta{Kind: "Pod"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name: "test",
						Env:  make([]v1.EnvVar, 0),
					},
				},
			},
		}

		m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&pod)
		require.NoError(t, err)

		uo := unstructured.Unstructured{Object: m}

		err = Setup(&uo)
		require.NoError(t, err)

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(uo.Object, &pod)
		require.NoError(t, err)

		// TODO(qbarrand) fix the method and then uncomment.
		// SetupPod does not set the resulting containers slice with unstructured.SetNestedSlice
		//env := pod.Spec.Containers[0].Env

		//assert.Contains(t, env, v1.EnvVar{Name: "HTTP_PROXY", Value: httpProxy})
		//assert.Contains(t, env, v1.EnvVar{Name: "HTTPS_PROXY", Value: httpsProxy})
		//assert.Contains(t, env, v1.EnvVar{Name: "NO_PROXY", Value: noProxy})
	})

	t.Run("DaemonSet with empty spec", func(t *testing.T) {
		ds := appsv1.DaemonSet{
			TypeMeta: metav1.TypeMeta{Kind: "DaemonSet"},
		}

		m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&ds)
		require.NoError(t, err)

		uo := unstructured.Unstructured{Object: m}

		err = Setup(&uo)
		require.Error(t, err)
	})

	t.Run("DaemonSet with one container template", func(t *testing.T) {
		t.Cleanup(resetProxy)

		const (
			httpProxy  = "http-host-with-proxy"
			httpsProxy = "https-host-with-proxy"
			noProxy    = "host-without-proxy"
		)

		ProxyConfiguration = Configuration{
			HttpProxy:  httpProxy,
			HttpsProxy: httpsProxy,
			NoProxy:    noProxy,
		}

		ds := appsv1.DaemonSet{
			TypeMeta: metav1.TypeMeta{Kind: "DaemonSet"},
			Spec: appsv1.DaemonSetSpec{
				Template: v1.PodTemplateSpec{
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name: "test",
								Env:  make([]v1.EnvVar, 0),
							},
						},
					},
				},
			},
		}

		m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&ds)
		require.NoError(t, err)

		uo := unstructured.Unstructured{Object: m}

		err = Setup(&uo)
		require.NoError(t, err)

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(uo.Object, &ds)
		require.NoError(t, err)

		// TODO(qbarrand) fix the method and then uncomment.
		// SetupDaemonSet does not set the resulting containers slice with unstructured.SetNestedSlice
		//env := ds.Spec.Template.Spec.Containers[0].Env

		//assert.Contains(t, env, v1.EnvVar{Name: "HTTP_PROXY", Value: httpProxy})
		//assert.Contains(t, env, v1.EnvVar{Name: "HTTPS_PROXY", Value: httpsProxy})
		//assert.Contains(t, env, v1.EnvVar{Name: "NO_PROXY", Value: noProxy})
	})
}

func TestClusterConfiguration(t *testing.T) {
	// TODO(qbarrand) make the DiscoveryClient in clients.HasResource injectable, so we can mock it.
}
