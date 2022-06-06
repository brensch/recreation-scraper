package main

import (
	"context"
	"testing"
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
