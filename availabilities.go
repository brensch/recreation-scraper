package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/brensch/recreation"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CollectionAvailabilities     = "availabilities"
	CollectionAvailabilityDeltas = "availability_deltas_grouped"

	SitesToSyncEachIteration = 25
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
	for i := 0; i < SitesToSyncEachIteration; i++ {
		idToRemove := rand.Intn(len(input))

		// add the chosen id to the array
		randomIDs = append(randomIDs, input[idToRemove])

		// then remove it from the array by replacing and truncating
		input[idToRemove] = input[len(input)-1]
		input = input[:len(input)-1]
	}

	return randomIDs
}

// GetAvailabilityRef gives us a consisten ID to work with for our documents
func GetAvailabilityRef(groundID string, targetTime time.Time) string {
	targetTime = recreation.GetStartOfMonth(targetTime)
	return fmt.Sprintf("%s-%s", groundID, targetTime.Format(time.RFC3339))
}

// CheckForAvailabilityChange gets old and new states of availability and returns the new availability and deltas to old
func CheckForAvailabilityChange(ctx context.Context, rec *recreation.Server, fs *firestore.Client, targetTime time.Time, targetGround string) (recreation.Availability, []CampsiteDelta, error) {

	// get new availability from API
	newAvailability, err := rec.GetAvailability(ctx, targetGround, targetTime)
	if err != nil {
		return recreation.Availability{}, nil, err
	}

	// get old abailability from firestore
	// NotFound errors ignored since the document not existing just results in an empty object, as intended
	oldAvailabilitySnap, err := fs.Collection(CollectionAvailabilities).Doc(GetAvailabilityRef(targetGround, targetTime)).Get(ctx)
	if err != nil && status.Code(err) != codes.NotFound {
		return recreation.Availability{}, nil, err
	}
	var oldAvailability recreation.Availability
	err = oldAvailabilitySnap.DataTo(&oldAvailability)
	if err != nil && status.Code(err) != codes.NotFound {
		return recreation.Availability{}, nil, err
	}

	// compare the old and new availabilities
	deltas, err := FindAvailabilityDeltas(oldAvailability, newAvailability)
	return newAvailability, deltas, err

}

func DoAvailabilitiesSync(ctx context.Context, log *zap.Logger, rec *recreation.Server, fs *firestore.Client, targetTime time.Time) error {

	start := time.Now()

	log = log.With(zap.Time("target_time", targetTime))
	log.Debug("starting availabilities sync")

	targetGrounds, err := GetGroundIDsToScrape(ctx, fs)
	if err != nil {
		log.Error("failed to get ground IDs to scrape", zap.Error(err))
		return err
	}

	var allDeltas []CampsiteDelta

	// iterate over grounds, looking for changes in availability for the given time
	for _, targetGround := range targetGrounds {
		log := log.With(zap.String("ground_id", targetGround))
		log.Debug("checking campground availability")
		newAvailability, deltas, err := CheckForAvailabilityChange(ctx, rec, fs, targetTime, targetGround)
		if err != nil {
			log.Error("failed to check availability change", zap.Error(err))
			return err
		}

		// if there are no deltas, continue
		if len(deltas) == 0 {
			log.Debug("found no deltas, continuing")
			continue
		}

		log.Debug("deltas found", zap.Int("delta_count", len(deltas)))
		allDeltas = append(allDeltas, deltas...)

		// update availabilities
		_, err = fs.Collection(CollectionAvailabilities).Doc(GetAvailabilityRef(targetGround, targetTime)).Set(
			ctx,
			newAvailability,
		)
		if err != nil {
			log.Error("couldn't add availability to firestore", zap.Error(err))
			return err
		}
	}

	// update deltas in firestore.
	// add all deltas to the same document since campsites are globally unique so cheaper to grab the whole
	// document and search through for campsites you want.
	log.Debug("syncing deltas to firestore", zap.Int("delta_count", len(allDeltas)))
	checkDelta := CheckDelta{
		Deltas:    allDeltas,
		CheckTime: start,
	}
	_, _, err = fs.Collection(CollectionAvailabilityDeltas).Add(
		ctx,
		checkDelta,
	)
	if err != nil {
		log.Error("couldn't add availability deltas to firestore", zap.Error(err))
		return err
	}

	log.Info("completed availabilities sync", zap.Duration("duration", time.Since(start)))

	return nil
}
