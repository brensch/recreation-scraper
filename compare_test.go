package main

import (
	"testing"

	"github.com/brensch/recreation"
)

func TestCompareCampgroundStates(t *testing.T) {

	old := recreation.Availability{
		Campsites: map[string]recreation.Campsite{
			"site_1": {
				Availabilities: map[string]string{
					"2022-06-01T00:00:00Z": "Not Reservable",
					"2022-06-02T00:00:00Z": "Not Reservable",
					"2022-06-03T00:00:00Z": "Available",
					"2022-06-04T00:00:00Z": "Available",
					"2022-06-05T00:00:00Z": "Available",
					"2022-06-06T00:00:00Z": "Not Reservable",
					"2022-06-07T00:00:00Z": "Available",
					"2022-06-08T00:00:00Z": "Available",
					"2022-06-09T00:00:00Z": "Available",
					"2022-06-10T00:00:00Z": "Available",
				},
			},
		},
	}

	new := recreation.Availability{
		Campsites: map[string]recreation.Campsite{
			"site_1": {
				Availabilities: map[string]string{
					"2022-06-01T00:00:00Z": "Not Reservable",
					"2022-06-02T00:00:00Z": "Not Reservable",
					"2022-06-03T00:00:00Z": "Available",
					"2022-06-04T00:00:00Z": "Available",
					"2022-06-05T00:00:00Z": "Available",
					"2022-06-06T00:00:00Z": "Available",
					"2022-06-07T00:00:00Z": "Available",
					"2022-06-08T00:00:00Z": "Available",
					"2022-06-09T00:00:00Z": "Available",
					"2022-06-10T00:00:00Z": "Not Reservable",
				},
			},
		},
	}

	delta, err := FindAvailabilityDeltas(old, new)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(delta)

}

func BenchmarkCompareCampgroundStates(b *testing.B) {

	old := recreation.Availability{
		Campsites: map[string]recreation.Campsite{
			"site_1": {
				Availabilities: map[string]string{
					"2022-06-01T00:00:00Z": "Not Reservable",
					"2022-06-02T00:00:00Z": "Not Reservable",
					"2022-06-03T00:00:00Z": "Available",
					"2022-06-04T00:00:00Z": "Available",
					"2022-06-05T00:00:00Z": "Available",
					"2022-06-06T00:00:00Z": "Not Reservable",
					"2022-06-07T00:00:00Z": "Available",
					"2022-06-08T00:00:00Z": "Available",
					"2022-06-09T00:00:00Z": "Available",
					"2022-06-10T00:00:00Z": "Available",
				},
			},
		},
	}

	new := recreation.Availability{
		Campsites: map[string]recreation.Campsite{
			"site_1": {
				Availabilities: map[string]string{
					"2022-06-01T00:00:00Z": "Not Reservable",
					"2022-06-02T00:00:00Z": "Not Reservable",
					"2022-06-03T00:00:00Z": "Available",
					"2022-06-04T00:00:00Z": "Available",
					"2022-06-05T00:00:00Z": "Available",
					"2022-06-06T00:00:00Z": "Available",
					"2022-06-07T00:00:00Z": "Available",
					"2022-06-08T00:00:00Z": "Available",
					"2022-06-09T00:00:00Z": "Available",
					"2022-06-10T00:00:00Z": "Not Reservable",
				},
			},
		},
	}

	for i := 0; i < b.N; i++ {
		FindAvailabilityDeltas(old, new)
	}
}
