package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/chart"
)

func TestFind(t *testing.T) {
	t.Parallel()

	s := []string{"a", "b", "c", "d"}

	assert.Equal(t, 2, Find(s, "c"))
	assert.Equal(t, len(s), Find(s, "z"))
}

func TestContains(t *testing.T) {
	t.Parallel()

	s := []string{"a", "b", "c", "d"}

	assert.True(t, Contains(s, "a"))
	assert.False(t, Contains(s, "z"))
}

func TestFindCRFile(t *testing.T) {
	t.Parallel()

	files := []*chart.File{
		{Name: "chart0.yaml"},
		{Name: "chart1.yaml"},
	}

	assert.Equal(t, 1, FindCRFile(files, "chart1"))
	assert.Equal(t, -1, FindCRFile(files, "chart99"))
}

func TestInsert(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		[]string{"c", "a", "b"},
		Insert([]string{"a", "b"}, 0, "c"),
	)

	assert.Equal(t,
		[]string{"a", "c", "b"},
		Insert([]string{"a", "b"}, 1, "c"),
	)

	assert.Equal(t,
		[]string{"a", "b", "c"},
		Insert([]string{"a", "b"}, 2, "c"),
	)
}
