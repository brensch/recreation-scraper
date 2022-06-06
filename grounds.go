package main

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/brensch/recreation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	CollectionGroundsMeta    = "grounds_meta"
	CollectionGroundsSummary = "grounds_summary"
	DocGroundsSummary        = "grounds_summary"
)

func ScrapeGrounds(ctx context.Context, client *recreation.Obfuscator) (recreation.SearchResults, error) {
	return recreation.DoSearchGeo(ctx, client, 37.3859, -122.0882)
}

func UpdateGrounds(ctx context.Context, fs *firestore.Client, grounds []recreation.CampGround) error {

	// batch all results and submit all updates at once.
	b := fs.Batch()
	for _, ground := range grounds {
		ref := fs.Collection(CollectionGroundsMeta).Doc(ground.EntityID)
		b.Create(ref, ground)
	}

	_, err := b.Commit(ctx)
	return err
}

// GroundSummary is just a list of all ground IDs to allow us to easily pull them all and save on firestore calls
type GroundSummary struct {
	GroundIDs []string `json:"ground_ids,omitempty" firestore:"ground_ids,omitempty"`
}

func InitGroundsSummary(ctx context.Context, fs *firestore.Client) error {
	_, err := fs.Collection(CollectionGroundsSummary).Doc(DocGroundsSummary).Create(ctx, GroundSummary{})
	if err != nil && status.Code(err) != codes.AlreadyExists {
		return err
	}
	return nil
}

func UpdateGroundsSummary(ctx context.Context, fs *firestore.Client, grounds []recreation.CampGround) error {

	// using interface because arrayunion requires it
	var groundIDS []interface{}
	for _, ground := range grounds {
		groundIDS = append(groundIDS, ground.EntityID)
	}

	_, err := fs.Collection(CollectionGroundsSummary).Doc(DocGroundsSummary).Update(ctx, []firestore.Update{
		{Path: "ground_ids", Value: firestore.ArrayUnion(groundIDS...)},
	})

	return err

}

// DoGroundsSync does the full routine of scraping from the website, syncing detailed data, then updating summary
func DoGroundsSync(ctx context.Context, client *recreation.Obfuscator, fs *firestore.Client) error {

	res, err := ScrapeGrounds(ctx, client)
	if err != nil {
		return err
	}

	err = UpdateGrounds(ctx, fs, res.Results)
	if err != nil {
		return err
	}

	err = UpdateGroundsSummary(ctx, fs, res.Results)
	if err != nil {
		return err
	}

	return nil

}
