// Copyright 2021 Mohammad Shafiee and The DMS Authors
//
// Licensed under the GNU General Public License, Version 3.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.gnu.org/licenses/gpl-3.0.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Copyright notice.

package dms

import (
	"errors"
	"fmt"
	"math"
)

// DMS represents a geographical coordinate in Degrees, Minutes, and Seconds format.
type DMS struct {
	Degree    uint    // Degree part of the coordinate.
	Minutes   uint    // Minute part of the coordinate.
	Seconds   float64 // Second part of the coordinate.
	Direction string  // Represents the cardinal direction (N, S, E, W).
}

// String Representations

// String returns the DMS format in a LTR representation.
func (d *DMS) String() string {
	return fmt.Sprintf(`%d°%d'%.02f" %s`, d.Degree, d.Minutes, d.Seconds, d.Direction)
}

// StringRTL returns the DMS format in a RTL representation.
func (d *DMS) StringRTL() string {
	return fmt.Sprintf(`%s "%.02f '%d °%d`, d.Direction, d.Seconds, d.Minutes, d.Degree)
}

// StringPersian returns the DMS format in Persian language representation.
func (d *DMS) StringPersian() string {
	return fmt.Sprintf(`%d درجه %d دقیقه %.02f ثانیه %s`, d.Degree, d.Minutes, d.Seconds, d.Direction)
}

// Rounding methods

// RoundToMinute rounds the coordinate value to the nearest minute.
func (d *DMS) RoundToMinute() {
	d.Seconds = roundToWholeNumber(d.Seconds)
	// Update minutes and degrees if needed after rounding.
	d.updateAfterRounding()
}

// RoundToSecond rounds the coordinate value to the nearest second.
func (d *DMS) RoundToSecond() {
	d.Seconds = roundToWholeNumber(d.Seconds + 0.5)
	// Update minutes and degrees if needed after rounding.
	d.updateAfterRounding()
}

// RoundToDegree rounds the coordinate value to the nearest degree.
func (d *DMS) RoundToDegree() {
	if d.Minutes >= 30 || (d.Minutes == 29 && d.Seconds >= 30) {
		d.Degree++
	}
	d.Minutes = 0
	d.Seconds = 0
}

// updateAfterRounding adjusts the Degree, Minute, and Second values after rounding.
func (d *DMS) updateAfterRounding() {
	if d.Seconds >= 60 {
		d.Seconds -= 60
		d.Minutes++
	}
	if d.Minutes >= 60 {
		d.Minutes -= 60
		d.Degree++
	}
}

// Factory functions

// NewDMS creates new DMS structures for given latitude and longitude.
func NewDMS(lat, lon float64) (*DMS, *DMS, error) {
	// Validate the input latitude and longitude.
	if math.Abs(lat) > 90 || math.Abs(lon) > 180 {
		return nil, nil, errors.New("Invalid latitude or longitude value")
	}
	latDMS := DecimalToDMS(lat, "N", "S")
	lonDMS := DecimalToDMS(lon, "E", "W")
	return latDMS, lonDMS, nil
}

// DecimalToDMS converts a decimal coordinate to DMS format.
func DecimalToDMS(decimalDegree float64, positiveIndicator, negativeIndicator string) *DMS {
	degree, minutes, seconds := decimalToDMSComponents(math.Abs(decimalDegree))
	direction := getDirectionForCoordinate(decimalDegree, positiveIndicator, negativeIndicator)
	return &DMS{Degree: degree, Minutes: minutes, Seconds: seconds, Direction: direction}
}

// DMSToDecimal converts a DMS format coordinate to its decimal representation.
func DMSToDecimal(dms DMS) float64 {
	return float64(dms.Degree) + float64(dms.Minutes)/60.0 + dms.Seconds/3600.0
}

// RoundDecimalToMinute rounds a decimal degree to its nearest minute.
func RoundDecimalToMinute(decimalDegree float64) float64 {
	degree := math.Floor(decimalDegree)
	decimalMinutes := (decimalDegree - degree) * 60
	roundedMinutes := roundToWholeNumber(decimalMinutes)
	return degree + (roundedMinutes / 60)
}

// RoundDecimalToSecond rounds a decimal degree to its nearest second.
func RoundDecimalToSecond(decimalDegree float64) float64 {
	degree := math.Floor(decimalDegree)
	minutes := math.Floor((decimalDegree - degree) * 60)
	decimalSeconds := (decimalDegree - degree - (minutes / 60)) * 3600
	roundedSeconds := roundToWholeNumber(decimalSeconds)
	return degree + (minutes+(roundedSeconds/60))/60
}

// Utility functions

// getDirectionForCoordinate determines the direction (N, S, E, W) of a given value.
func getDirectionForCoordinate(value float64, positiveIndicator, negativeIndicator string) string {
	if value >= 0 {
		return positiveIndicator
	}
	return negativeIndicator
}

// decimalToDMSComponents breaks down a decimal degree into its D, M, S components.
func decimalToDMSComponents(decimalDegree float64) (degree uint, minutes uint, seconds float64) {
	degree = uint(decimalDegree)
	minutes = uint((decimalDegree - float64(degree)) * 60)
	seconds = (decimalDegree - float64(degree) - float64(minutes)/60) * 3600
	return
}

// roundToWholeNumber rounds a float to its nearest whole number.
func roundToWholeNumber(value float64) float64 {
	return math.Round(value)
}
