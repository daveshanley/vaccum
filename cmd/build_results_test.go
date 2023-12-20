package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildResults(t *testing.T) {
	_, _, err := BuildResults("nuggets", nil, nil, "", 5)
	assert.Error(t, err)
}

func TestBuildResults_SkipCheck(t *testing.T) {
	_, _, err := BuildResultsWithDocCheckSkip("nuggets", nil, nil, "", true, 5)
	assert.Error(t, err)
}
