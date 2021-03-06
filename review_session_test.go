package leaf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReviewSession(t *testing.T) {
	cards := []CardWithStats{
		{Card{"foo", "foo", []string{"bar"}}, NewStats(SRSSupermemo2PlusCustom)},
		{Card{"bar", "foo", []string{"baz"}}, NewStats(SRSSupermemo2PlusCustom)},
	}

	stats := make(map[string]*Stats)
	s := NewReviewSession(cards, RatingTypeAuto, func(card *CardWithStats) error {
		stats[card.Question] = card.Stats
		return nil
	})

	t.Run("StartedAt", func(t *testing.T) {
		assert.NotNil(t, s.StartedAt())
	})

	t.Run("RatingType", func(t *testing.T) {
		assert.Equal(t, RatingTypeAuto, s.RatingType())
	})

	t.Run("Total", func(t *testing.T) {
		assert.Equal(t, 2, s.Total())
	})

	t.Run("Left", func(t *testing.T) {
		assert.Equal(t, 2, s.Left())
	})

	t.Run("Next", func(t *testing.T) {
		assert.Equal(t, "foo", s.Next())
	})

	t.Run("CorrectAnswer", func(t *testing.T) {
		assert.Equal(t, "bar", s.CorrectAnswer())
	})

	t.Run("Rate - incorrect", func(t *testing.T) {
		require.NoError(t, s.Again())
		assert.Equal(t, 2, s.Left())
		assert.Equal(t, "bar", s.Next())
	})

	t.Run("Rate - correct", func(t *testing.T) {
		require.NoError(t, s.Rate(1))
		assert.Equal(t, 1, s.Left())
		assert.Equal(t, "foo", s.Next())
	})

	t.Run("Rate - multiple incorrect", func(t *testing.T) {
		for i := 0; i < 4; i++ {
			require.NoError(t, s.Again())
		}
		assert.Equal(t, 1, s.Left())
	})

	t.Run("Answer - finish session", func(t *testing.T) {
		require.NoError(t, s.Rate(0))
		assert.Equal(t, 0, s.Left())
	})

	fooStats := stats["foo"].SRSAlgorithm.(*Supermemo2PlusCustom)
	assert.InDelta(t, 0.52, fooStats.Difficulty, 0.01)
	assert.InDelta(t, 0.2, fooStats.Interval, 0.01)

	barStats := stats["bar"].SRSAlgorithm.(*Supermemo2PlusCustom)
	assert.InDelta(t, 0.28, barStats.Difficulty, 0.01)
	assert.InDelta(t, 0.37, barStats.Interval, 0.01)
}
