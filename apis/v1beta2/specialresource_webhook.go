package v1beta2

import ctrl "sigs.k8s.io/controller-runtime"

func (in *SpecialResource) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}
