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

package dms

import (
	"fmt"
)

// LatLonError is used for errors with lat/lon values
type LatLonError struct {
	err string
}

func (e *LatLonError) Error() string {
	return e.err
}

// DMS coordinates
type DMS struct {
	Degree    uint
	Minutes   uint
	Seconds   float64
	Direction string
}

func (d *DMS) String() string {
	if d != nil {
		return fmt.Sprintf(`%d°%d'%.02f" %s`, d.Degree, d.Minutes, d.Seconds, d.Direction)
	}
	return ""
}
func (d *DMS) StringRTL() string {
	return fmt.Sprintf(`%s "%.02f '%d °%d`, d.Direction, d.Seconds, d.Minutes, d.Degree)
}
func (d *DMS) StringPersian() string {
	return fmt.Sprintf(`%d درجه %d دقیقه %.02f ثانیه%s`, d.Degree, d.Minutes, d.Seconds, d.Direction)
}

// NewDMS converts Decimal Degreees to SignLongitudeDMS, Minute, Seconds coordinates
func NewDMS(lat, lon float64) (*DMS, *DMS, error) {
	if lat < 0 || lon < 0 {
		return nil, nil, &LatLonError{"LatitudeDMS or longitude must be positive."}
	}
	if lat > 90 || lon > 180 {
		return nil, nil, &LatLonError{"LatitudeDMS must be less than 90 and longitude must be less than 180."}
	}

	var latDirection string
	var lonDirection string
	if lat > 0 {
		latDirection = "N"
	} else {
		latDirection = "S"
	}

	if lon > 0 {
		lonDirection = "E"
	} else {
		lonDirection = "W"
	}

	latitude := uint(lat)
	latitudeMinutes := uint((lat - float64(latitude)) * 60)
	latitudeSeconds := (lat - float64(latitude) - float64(latitudeMinutes)/60) * 3600

	longitude := uint(lon)
	longitudeMinutes := uint((lon - float64(longitude)) * 60)
	longitudeSeconds := (lon - float64(longitude) - float64(longitudeMinutes)/60) * 3600

	return &DMS{Degree: latitude, Minutes: latitudeMinutes, Seconds: latitudeSeconds, Direction: latDirection},
		&DMS{Degree: longitude, Minutes: longitudeMinutes, Seconds: longitudeSeconds, Direction: lonDirection},
		nil
}

// DecimalToDMS converts Decimal Degrees to SignLongitudeDMS, Minute, Seconds coordinates
func DecimalToDMS(decimalDegree float64) *DMS {
	degree := uint(decimalDegree)
	minutes := uint((decimalDegree - float64(degree)) * 60)
	seconds := (decimalDegree - float64(degree) - float64(minutes)/60) * 3600

	return &DMS{Degree: degree, Minutes: minutes, Seconds: seconds}
}

func DMSToDecimal(dms DMS) (d float64) {
	d = float64(dms.Degree)
	d += float64(dms.Minutes) / 60.0
	d += dms.Seconds / 3600
	return d
}
