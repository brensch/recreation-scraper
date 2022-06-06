package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/brensch/recreation"
)

const (
	CollectionAvailabilities = "availabilities"
)

// GetGroundsToScrape returns everything we want to scrape. Currently it's just selecting 5 random sites
func GetGroundIDsToScrape(ctx context.Context, fs *firestore.Client) ([]string, error) {
	snap, err := fs.Collection(CollectionGroundsSummary).Doc(DocGroundsSummary).Get(ctx)
	if err != nil {
		return nil, err
	}

	var summary GroundSummary
	err = snap.DataTo(&summary)
	if err != nil {
		return nil, err
	}

	ids := summary.GroundIDs
	return SelectRandomIDs(ids), nil
}

// SelectRandomIDs picks 5 sites at random, ensuring they're all different
func SelectRandomIDs(input []string) []string {
	var randomIDs []string
	rand.Seed(time.Now().UnixMicro())
	for i := 0; i < 5; i++ {
		idToRemove := rand.Intn(len(input))

		// add the chosen id to the array
		randomIDs = append(randomIDs, input[idToRemove])

		// then remove it from the array by replacing and truncating
		input[idToRemove] = input[len(input)-1]
		input = input[:len(input)-1]
	}

	return randomIDs
}

func GetAvailabilityRef(groundID string, targetTime time.Time) string {
	targetTime = recreation.GetStartOfMonth(targetTime)
	return fmt.Sprintf("%s-%s", groundID, targetTime.Format(time.RFC3339))
}

// CheckForAvailabilityChange gets old and new states of availability and returns the deltas
func CheckForAvailabilityChange(ctx context.Context, client *recreation.Obfuscator, fs *firestore.Client, targetTime time.Time, targetGround string) ([]CampsiteDelta, error) {

	newAvailability, err := recreation.GetAvailability(ctx, client, targetGround, targetTime)
	if err != nil {
		return nil, err
	}

	oldAvailabilitySnap, err := fs.Collection(CollectionAvailabilities).Doc(GetAvailabilityRef(targetGround, targetTime)).Get(ctx)
	if err != nil {
		return nil, err
	}

	var oldAvailability recreation.Availability
	err = oldAvailabilitySnap.DataTo(&oldAvailability)
	if err != nil {
		return nil, err
	}

	return FindCampsiteDeltas(oldAvailability, newAvailability)

}

func DoAvailabilitiesSync(ctx context.Context, client *recreation.Obfuscator, fs *firestore.Client, targetTime time.Time) error {

	targetGrounds, err := GetGroundIDsToScrape(ctx, fs)
	if err != nil {
		return err
	}

	// iterate over grounds
	for _, targetGround := range targetGrounds {
		deltas, err := CheckForAvailabilityChange(ctx, client, fs, targetTime, targetGround)
		if err != nil {
			return err
		}

		fmt.Println(deltas)
	}

	return nil
}
