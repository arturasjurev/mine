package factory_test

import (
	"testing"
	"time"

	"github.com/sheirys/mine/factory"
	"github.com/stretchr/testify/assert"
)

func TestCalculateProcessTime(t *testing.T) {
	testTable := []struct {
		Challenge, Power int
		ExpectedDuration time.Duration
	}{
		{
			Challenge:        1000,
			Power:            500,
			ExpectedDuration: 2 * time.Second,
		},
		{
			Challenge:        500,
			Power:            1000,
			ExpectedDuration: 500 * time.Millisecond,
		},
		{
			Challenge:        2500,
			Power:            100,
			ExpectedDuration: 25 * time.Second,
		},
		{
			Challenge:        500,
			Power:            2000,
			ExpectedDuration: 250 * time.Millisecond,
		},
		{
			Challenge:        1000,
			Power:            4000,
			ExpectedDuration: 250 * time.Millisecond,
		},
		{
			Challenge:        5000,
			Power:            0,
			ExpectedDuration: 0,
		},
	}

	for i, v := range testTable {
		result := factory.CalculateProcessTime(v.Challenge, v.Power)
		assert.Equal(t, v.ExpectedDuration, result, "case %d", i)
	}
}
