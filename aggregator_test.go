package health

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAggregator(t *testing.T) {
	a := newAggregator(time.Minute, time.Minute*5)
	assert.Equal(t, time.Minute, a.intervalDuration)
	assert.Equal(t, time.Minute*5, a.retain)
	assert.Equal(t, 5, a.maxIntervals)
	assert.Equal(t, 0, len(a.intervalAggregations))
	assert.NotNil(t, a.intervalAggregations)
}
