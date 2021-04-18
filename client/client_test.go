package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BuildRequest(t *testing.T) {
	prediction := NewPrediction()
	prediction.WithCsv("test").WithLabeled().WithMetrics([]string{"rmse"}).WithLabeled()
	assert.Equal(t, prediction.format, "csv")
}
