package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/chart"
)

func TestGenerateName(t *testing.T) {
	t.Cleanup(func() { CurrentName = "" })

	f := &chart.File{Name: "/path/to/test.json"}

	GenerateName(f, "some-sr")

	assert.Equal(t, "specialresource.openshift.io/state-some-sr-test", CurrentName)
}
