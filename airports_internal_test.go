package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestPositionCalculation(t *testing.T) {
	lat, lon := calculatePosition("Point(120.210435142 60.38344912)")
	assert.EqualValues(t, "60.38344912", lat)
	assert.EqualValues(t, "120.210435142", lon)
	lat, lon = calculatePosition("")
	assert.EqualValues(t, "", lat)
	assert.EqualValues(t, "", lon)
	lat, lon = calculatePosition("t2177715128")
	assert.EqualValues(t, "", lat)
	assert.EqualValues(t, "", lon)
}
func TestDuplicateCounter(t *testing.T) {
	duplicate := duplicateCounter{}
	file, err := ioutil.TempFile(t.TempDir(), "airports_duplicate_test.db")
	if err != nil {
		t.Fatalf("unable to create temporary file for testing")
	}
	duplicate.duplicateFile = file
	duplicate.checkIcao("XGRT")
	assert.EqualValues(t, 0, duplicate.duplicateCounter)
	duplicate.checkIcao("XXYZ")
	duplicate.checkIcao("XXYZ")
	assert.EqualValues(t, 1, duplicate.duplicateCounter)
	duplicate.checkIcao("XXZZ")
	assert.EqualValues(t, 1, duplicate.duplicateCounter)
	duplicate.checkIcao("XZZZ")
	duplicate.checkIcao("XZZZ")
	duplicate.checkIcao("XZZZ")
	duplicate.checkIcao("XZZZ")
	assert.EqualValues(t, 2, duplicate.duplicateCounter)
	duplicate.checkIcao("ZZZZ")
	duplicate.checkIcao("ZZZZ")
	assert.EqualValues(t, 3, duplicate.duplicateCounter)
	duplicate.checkIcao("YYYY")
	duplicate.checkIcao("YYYY")
	duplicate.checkIcao("YYYY")
	duplicate.checkIcao("YYYY")
	duplicate.checkIcao("YYYY")
	assert.EqualValues(t, 4, duplicate.duplicateCounter)
}
