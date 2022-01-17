package v1beta1

import (
	"github.com/openshift-psap/special-resource-operator/apis/v1beta2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertFrom converts from the Hub version (v1beta1) to this version.
func (in *SpecialResource) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta2.SpecialResource)

	// Metadata

	in.ObjectMeta = src.ObjectMeta

	// Spec

	in.Spec.Chart = src.Spec.Chart
	in.Spec.Namespace = src.Spec.Namespace
	in.Spec.Debug = src.Spec.Debug
	in.Spec.Set = src.Spec.Set
	in.Spec.NodeSelector = src.Spec.NodeSelector
	in.Spec.Dependencies = make([]SpecialResourceDependency, len(src.Spec.Dependencies))

	for i, d := range src.Spec.Dependencies {
		in.Spec.Dependencies[i] = SpecialResourceDependency(d)
	}

	// Status

	in.Status = SpecialResourceStatus{State: src.Status.State}

	return nil
}

// ConvertTo converts this SpecialResource to the Hub version (v1beta1).
func (in *SpecialResource) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta2.SpecialResource)

	// Metadata

	dst.ObjectMeta = in.ObjectMeta

	// Spec

	dst.Spec.Chart = in.Spec.Chart
	dst.Spec.Namespace = in.Spec.Namespace
	dst.Spec.Debug = in.Spec.Debug
	dst.Spec.Set = in.Spec.Set
	dst.Spec.NodeSelector = in.Spec.NodeSelector

	dst.Spec.Dependencies = make([]v1beta2.SpecialResourceDependency, len(in.Spec.Dependencies))

	for i, d := range in.Spec.Dependencies {
		dst.Spec.Dependencies[i] = v1beta2.SpecialResourceDependency(d)
	}

	// Status

	dst.Status.State = in.Status.State

	return nil
}
