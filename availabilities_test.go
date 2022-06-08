package main

import (
	"context"
	"testing"
	"time"

	"github.com/brensch/recreation"
	"go.uber.org/zap"
)

func TestSelectRandomIDs(t *testing.T) {

	// since it's random, do it a lot of times
	for i := 0; i < 1000; i++ {

		IDs := []string{"1", "2", "3", "4", "5", "6", "7"}
		randomIDs := SelectRandomIDs(IDs)

		// validate got the right number of ids
		if len(randomIDs) != 5 {
			t.Log("got wrong number of IDs:", len(randomIDs))
			t.Fail()
		}

		// check there are no duplicates
		var checkedIDs []string
		for _, randomID := range randomIDs {
			for _, checkedID := range checkedIDs {
				if randomID == checkedID {
					t.Log("got duplicate ID", randomID, checkedID)
					t.FailNow()
				}
			}
			checkedIDs = append(checkedIDs, randomID)
		}
	}
}

func TestGetGroundIDsToScrape(t *testing.T) {

	ctx := context.Background()
	// TODO: use a local firestore instance for this
	fs, err := InitFirestore(ctx)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer fs.Close()

	ids, err := GetGroundIDsToScrape(ctx, fs)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(ids)
}

func TestCheckForAvailabilityChange(t *testing.T) {

	ctx := context.Background()
	// TODO: use a local firestore instance for this
	fs, err := InitFirestore(ctx)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer fs.Close()

	log, _ := zap.NewProduction()
	rec := recreation.InitServer(ctx, log, 0)

	newAvailabilities, deltas, err := CheckForAvailabilityChange(ctx, rec, fs, time.Now(), "232784")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// TODO: since this is an integration test with firebase need to figure out what to expect
	t.Log(newAvailabilities, deltas)
}

func TestDoAvailabilitiesSync(t *testing.T) {

	ctx := context.Background()
	// TODO: use a local firestore instance for this
	fs, err := InitFirestore(ctx)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer fs.Close()

	log, _ := zap.NewProduction()
	rec := recreation.InitServer(ctx, log, 0)

	err = DoAvailabilitiesSync(ctx, log, rec, fs, time.Now())
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// TODO: since this is an integration test with firebase need to figure out what to expect

}
