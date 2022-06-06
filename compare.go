package main

import (
	"time"

	"github.com/brensch/recreation"
)

type CampsiteDelta struct {
	SiteID   string
	OldState string
	NewState string
	Date     time.Time
}

func FindCampsiteDeltas(oldGround, newGround recreation.Availability) ([]CampsiteDelta, error) {

	var deltas []CampsiteDelta

	// iterate through each field in new and check what the previous value was
	for siteID, newSite := range newGround.Campsites {
		oldSite := oldGround.Campsites[siteID]
		for dateString, availability := range newSite.Availabilities {

			// ignore things that haven't changed.
			// using a map here is nice, i think it's efficient. May try other approaches if i get frisky
			if oldSite.Availabilities[dateString] == availability {
				continue
			}

			date, err := time.Parse(time.RFC3339, dateString)
			if err != nil {
				return nil, err
			}

			deltas = append(deltas, CampsiteDelta{
				SiteID:   siteID,
				OldState: oldSite.Availabilities[dateString],
				NewState: availability,
				Date:     date,
			})
		}

	}

	return deltas, nil

}
