package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseBadInterval(t *testing.T) {
	_, err := ParseTimeIntervals("9-")
	if err == nil {
		t.Errorf("Expected error")
	}

	_, err = ParseTimeIntervals("-12")
	if err == nil {
		t.Errorf("Expected error")
	}

	_, err = ParseTimeIntervals("9-12,10")
	if err == nil {
		t.Errorf("Expected error")
	}

	_, err = ParseTimeIntervals(".30-12")
	if err == nil {
		t.Errorf("Expected error")
	}

	_, err = ParseTimeIntervals("10.-12")
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestBadStartEnd(t *testing.T) {
	_, err := ParseTimeIntervals("9-8")
	if err == nil {
		t.Errorf("Expected error")
	}

	_, err = ParseTimeIntervals("9-9")
	if err == nil {
		t.Errorf("Expected error")
	}
}

func TestGoodSingleInterval(t *testing.T) {

	assert := assert.New(t)

	res, err := ParseTimeIntervals("5-12")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	assert.Len(res, 1)
	assert.Equal(5, res["5-12"][0].Hour())
	assert.Equal(12, res["5-12"][1].Hour())

	res, err = ParseTimeIntervals("5.10-12.40")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	assert.Len(res, 1)
	assert.Equal(5, res["5.10-12.40"][0].Hour())
	assert.Equal(10, res["5.10-12.40"][0].Minute())
	assert.Equal(12, res["5.10-12.40"][1].Hour())
	assert.Equal(40, res["5.10-12.40"][1].Minute())
}

func TestGoodMultipleIntervals(t *testing.T) {

	assert := assert.New(t)

	res, err := ParseTimeIntervals("9-12,13-16,14-18.30")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	assert.Len(res, 3)
	assert.Equal(9, res["9-12"][0].Hour())
	assert.Equal(12, res["9-12"][1].Hour())
	assert.Equal(13, res["13-16"][0].Hour())
	assert.Equal(16, res["13-16"][1].Hour())
	assert.Equal(14, res["14-18.30"][0].Hour())
	assert.Equal(18, res["14-18.30"][1].Hour())
	assert.Equal(30, res["14-18.30"][1].Minute())

	res, err = ParseTimeIntervals("5.10-12.40")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	assert.Len(res, 1)
}
