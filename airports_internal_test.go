package main
import "testing"
import "github.com/stretchr/testify/assert"

func TestPositionCalculation(t *testing.T){
	lat,lon := calculatePosition("Point(120.210435142 60.38344912)")
	assert.EqualValues(t,"60.38344912",lat)
	assert.EqualValues(t,"120.210435142",lon)
	lat,lon = calculatePosition("")
	assert.EqualValues(t,"",lat)
	assert.EqualValues(t,"",lon)
	lat,lon = calculatePosition("t2177715128")
	assert.EqualValues(t,"",lat)
	assert.EqualValues(t,"",lon)
}