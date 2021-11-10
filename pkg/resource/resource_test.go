package resource

import (
	"testing"

	buildv1 "github.com/openshift/api/build/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestIsNamespaced(t *testing.T) {
	t.Parallel()

	assert.True(t, IsNamespaced("Pod"))

	namespacedTypes := []string{
		"Namespace",
		"ClusterRole",
		"ClusterRoleBinding",
		"SecurityContextConstraint",
		"SpecialResource",
	}

	for _, nt := range namespacedTypes {
		assert.False(t, IsNamespaced(nt))
	}
}

func TestIsNotUpdateable(t *testing.T) {
	t.Parallel()

	assert.False(t, IsNotUpdateable("Deployment"))

	for _, nt := range []string{"ServiceAccount", "Pod"} {
		assert.True(t, IsNotUpdateable(nt))
	}
}

func TestNeedsResourceVersionUpdate(t *testing.T) {
	t.Parallel()

	assert.False(t, NeedsResourceVersionUpdate("Pod"))

	resourceTypes := []string{
		"SecurityContextConstraints",
		"Service",
		"ServiceMonitor",
		"Route",
		"Build",
		"BuildRun",
		"BuildConfig",
		"ImageStream",
		"PrometheusRule",
		"CSIDriver",
		"Issuer",
		"CustomResourceDefinition",
		"Certificate",
		"SpecialResource",
		"OperatorGroup",
		"CertManager",
		"MutatingWebhookConfiguration",
		"ValidatingWebhookConfiguration",
		"Deployment",
		"ImagePolicy",
	}

	for _, rt := range resourceTypes {
		assert.True(t, NeedsResourceVersionUpdate(rt))
	}

}

func TestUpdateResourceVersion(t *testing.T) {
	t.Parallel()

	t.Run("Pod: nothing happens", func(t *testing.T) {
		foundPod := v1.Pod{
			TypeMeta: metav1.TypeMeta{Kind: "Pod"},
		}

		foundMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&foundPod)
		require.NoError(t, err)

		err = UpdateResourceVersion(&unstructured.Unstructured{}, &unstructured.Unstructured{Object: foundMap})
		require.NoError(t, err)
	})

	t.Run("Service with no resourceVersion", func(t *testing.T) {
		foundSvc := v1.Service{
			TypeMeta: metav1.TypeMeta{Kind: "Service"},
		}

		foundMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&foundSvc)
		require.NoError(t, err)

		err = UpdateResourceVersion(&unstructured.Unstructured{}, &unstructured.Unstructured{Object: foundMap})
		require.Error(t, err)
	})

	t.Run("Service with no clusterIP", func(t *testing.T) {
		foundSvc := v1.Service{
			TypeMeta:   metav1.TypeMeta{Kind: "Service"},
			ObjectMeta: metav1.ObjectMeta{ResourceVersion: "123"},
		}

		foundMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&foundSvc)
		require.NoError(t, err)

		reqUnstructured := unstructured.Unstructured{
			Object: make(map[string]interface{}),
		}

		err = UpdateResourceVersion(&reqUnstructured, &unstructured.Unstructured{Object: foundMap})
		require.Error(t, err)
	})

	t.Run("Service with clusterIP", func(t *testing.T) {
		const (
			clusterIP       = "1.2.3.4"
			resourceVersion = "123"
		)

		foundSvc := v1.Service{
			TypeMeta:   metav1.TypeMeta{Kind: "Service"},
			ObjectMeta: metav1.ObjectMeta{ResourceVersion: resourceVersion},
			Spec:       v1.ServiceSpec{ClusterIP: clusterIP},
		}

		foundMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&foundSvc)
		require.NoError(t, err)

		reqUnstructured := unstructured.Unstructured{
			Object: make(map[string]interface{}),
		}

		err = UpdateResourceVersion(&reqUnstructured, &unstructured.Unstructured{Object: foundMap})
		require.NoError(t, err)

		reqSvc := v1.Service{}

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(reqUnstructured.Object, &reqSvc)
		require.NoError(t, err)

		assert.Equal(t, resourceVersion, reqSvc.GetResourceVersion())
		assert.Equal(t, clusterIP, reqSvc.Spec.ClusterIP)
	})
}

func TestSetNodeSelectorTerms(t *testing.T) {
	t.Parallel()

	t.Run("DaemonSet", func(t *testing.T) {
		d := appsv1.DaemonSet{
			TypeMeta: metav1.TypeMeta{Kind: "DaemonSet"},
		}

		m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&d)
		require.NoError(t, err)

		err = unstructured.SetNestedStringMap(m, make(map[string]string), "spec", "template", "spec", "nodeSelector")
		require.NoError(t, err)

		terms := map[string]string{"key": "value"}
		uo := unstructured.Unstructured{Object: m}

		err = SetNodeSelectorTerms(&uo, terms)
		require.NoError(t, err)

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(uo.Object, &d)
		require.NoError(t, err)

		assert.Equal(t, terms, d.Spec.Template.Spec.NodeSelector)
	})

	t.Run("Deployment", func(t *testing.T) {
		d := appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{Kind: "Deployment"},
		}

		m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&d)
		require.NoError(t, err)

		err = unstructured.SetNestedStringMap(m, make(map[string]string), "spec", "template", "spec", "nodeSelector")
		require.NoError(t, err)

		terms := map[string]string{"key": "value"}
		uo := unstructured.Unstructured{Object: m}

		err = SetNodeSelectorTerms(&uo, terms)
		require.NoError(t, err)

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(uo.Object, &d)
		require.NoError(t, err)

		assert.Equal(t, terms, d.Spec.Template.Spec.NodeSelector)
	})

	// TODO(qbarrand) this bugs because the code checks if the kind is Statefulset (no capital S)
	//t.Run("StatefulSet", func(t *testing.T) {
	//	statefulSet := appsv1.StatefulSet{
	//		TypeMeta: metav1.TypeMeta{Kind: "StatefulSet"},
	//	}
	//
	//	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&statefulSet)
	//	require.NoError(t, err)
	//
	//	terms := map[string]string{"key": "value"}
	//	uo := unstructured.Unstructured{Object: m}
	//
	//	err = SetNodeSelectorTerms(&uo, terms)
	//	require.NoError(t, err)
	//
	//	err = runtime.DefaultUnstructuredConverter.FromUnstructured(uo.Object, &statefulSet)
	//	require.NoError(t, err)
	//
	//	assert.Equal(t, terms, statefulSet.Spec.Template.Spec.NodeSelector)
	//})

	t.Run("Pod", func(t *testing.T) {
		p := v1.Pod{
			TypeMeta: metav1.TypeMeta{Kind: "Pod"},
		}

		m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&p)
		require.NoError(t, err)

		err = unstructured.SetNestedStringMap(m, make(map[string]string), "spec", "nodeSelector")
		require.NoError(t, err)

		terms := map[string]string{"key": "value"}
		uo := unstructured.Unstructured{Object: m}

		err = SetNodeSelectorTerms(&uo, terms)
		require.NoError(t, err)

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(uo.Object, &p)
		require.NoError(t, err)

		assert.Equal(t, terms, p.Spec.NodeSelector)
	})

	t.Run("BuildConfig", func(t *testing.T) {
		d := buildv1.BuildConfig{
			TypeMeta: metav1.TypeMeta{Kind: "BuildConfig"},
		}

		m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&d)
		require.NoError(t, err)

		err = unstructured.SetNestedStringMap(m, make(map[string]string), "spec", "nodeSelector")
		require.NoError(t, err)

		terms := map[string]string{"key": "value"}
		uo := unstructured.Unstructured{Object: m}

		err = SetNodeSelectorTerms(&uo, terms)
		require.NoError(t, err)

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(uo.Object, &d)
		require.NoError(t, err)

		assert.Equal(t, buildv1.OptionalNodeSelector(terms), d.Spec.NodeSelector)
	})
}

// TODO(qbarrand) test CreateFromYAML()

func TestIsOneTimer(t *testing.T) {
	t.Run("Service", func(t *testing.T) {
		svc := v1.Service{
			TypeMeta: metav1.TypeMeta{Kind: "Service"},
		}

		m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&svc)
		require.NoError(t, err)

		ret, err := IsOneTimer(&unstructured.Unstructured{Object: m})
		require.NoError(t, err)
		assert.False(t, ret)
	})

	t.Run("Pod, restartPolicy undefined", func(t *testing.T) {
		pod := v1.Pod{
			TypeMeta: metav1.TypeMeta{Kind: "Pod"},
		}

		m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&pod)
		require.NoError(t, err)

		_, err = IsOneTimer(&unstructured.Unstructured{Object: m})
		require.Error(t, err)
	})

	t.Run("Pod, restartPolicy Never", func(t *testing.T) {
		pod := v1.Pod{
			TypeMeta: metav1.TypeMeta{Kind: "Pod"},
			Spec:     v1.PodSpec{RestartPolicy: "Never"},
		}

		m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&pod)
		require.NoError(t, err)

		ret, err := IsOneTimer(&unstructured.Unstructured{Object: m})
		require.NoError(t, err)

		assert.True(t, ret)
	})
}

// TODO(qbarrand) test CRUD()

func TestSetMetaData(t *testing.T) {
	uo := unstructured.Unstructured{Object: make(map[string]interface{})}

	const (
		name      = "test-name"
		namespace = "test-namespace"
	)

	SetMetaData(&uo, name, namespace)

	assert.Equal(t, name, uo.GetAnnotations()["meta.helm.sh/release-name"])
	assert.Equal(t, namespace, uo.GetAnnotations()["meta.helm.sh/release-namespace"])
	assert.Equal(t, "Helm", uo.GetLabels()["app.kubernetes.io/managed-by"])
}

// TODO(qbarrand) test BeforeCRUD()

// TODO(qbarrand) test AfterCRUD()
