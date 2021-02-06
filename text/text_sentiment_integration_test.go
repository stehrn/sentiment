package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSentiment(t *testing.T) {

	// happy
	sentiment, err := scoreSentiment("I love things right now, they are great")
	if err != nil {
		t.Fatal("Could not score happy sentiment:", err)
	}

	t.Logf("Score (happy) is: %f", sentiment.GetScore())
	assert.True(t, (sentiment.GetScore() > 0.8), "Expected score greate than 0.8")

	// sad
	sentiment, err = scoreSentiment("I hate things right now, they are bad")
	if err != nil {
		t.Fatal("Could not score bad sentiment:", err)
	}

	t.Logf("Score (bad) is: %f", sentiment.GetScore())
	assert.True(t, (sentiment.GetScore() < -0.7), "Expected score greate than 0.8")
}