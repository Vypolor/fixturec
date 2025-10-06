package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToGen_Sum(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Arrange
		f := setUp(t)
		defer f.tearDown()

		f.myType2.EXPECT().Call2().Return(2)
		f.myType1.EXPECT().Call1().Return(1)

		// Act
		got := f.toGen.Sum()

		// Assert
		assert.Equal(t, 3, got)
	})
}
