package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/brensch/recreation"
	"google.golang.org/api/option"
)

// var (
// 	Park     = "232770"
// 	Campsite = "86766"
// )

// func init() {
// 	flag.StringVar(&Park, "park", Park, "Which park to check")
// 	flag.StringVar(&Campsite, "campsite", Campsite, "Which campsite to check")

// }

func main() {
	flag.Parse()

	ctx := context.Background()
	client := recreation.InitObfuscator(ctx)
	_ = client

	fs, err := InitFirestore(ctx)
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	// init collection for ground summary
	err = InitGroundsSummary(ctx, fs)
	if err != nil {
		panic(err)
	}

	fmt.Println("starting")

	log.Print("starting server...")
	http.HandleFunc("/avails", HandleScrapeAvailabilities(client, fs))
	http.HandleFunc("/grounds", HandleScrapeGrounds(client, fs))

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

	// GetAllAvailabilities(ctx, client)

}

func InitFirestore(ctx context.Context) (*firestore.Client, error) {
	opt := option.WithCredentialsFile("creds.json")
	conf := &firebase.Config{ProjectID: "campr-app"}

	app, err := firebase.NewApp(context.Background(), conf, opt)
	if err != nil {
		return nil, err
	}

	return app.Firestore(ctx)

}

// Ground is a silly name meaning a campground in the format that i want it.
// the idea is that i will group every site into the different states present to help doing deltas
// type Ground struct {
// 	GroundID string
// 	Sites    map[string]Site
// }

type siteState int

const (
	stateAvailable siteState = iota
	stateReserved
	stateNotReservableManagement
)

// TODO: map states to enums
// var (
// 	stateMappings =
// )

// type Site struct {
// 	// ID string
// 	// this is actually a date as the key
// 	// todo make values sitestates
// 	Availabilities map[string]string
// }

// func MergeAvailabilities(avails []recreation.Availability) recreation.Availability {
// 	s := Ground{
// 		GroundID: groundID,
// 		Sites:    make(map[string]Site),
// 	}

// 	for _, avail := range avails {
// 		for _, site := range avail.Campsites {
// 			// TODO: trim times for space savings
// 			// probably not necessary for a while and without it there's the possibility to go back to time
// 			_, ok := s.Sites[site.CampsiteID]

// 			// add new sites
// 			if !ok {
// 				s.Sites[site.CampsiteID] = Site{
// 					Availabilities: site.Availabilities,
// 				}
// 				continue
// 			}

// 			// add all new dates if site already exists
// 			for date, state := range site.Availabilities {
// 				s.Sites[site.CampsiteID].Availabilities[date] = state
// 			}
// 		}
// 	}

// 	return s
// }

// func GetAllAvailabilities(ctx context.Context, client recreation.HTTPClient) {

// 	res, err := recreation.DoSearchGeo(ctx, client, 37.3859, -122.0882)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Println(len(res.Results))

// 	for _, campground := range res.Results[0:50] {
// 		fmt.Println(campground.EntityID, campground.Name, campground.City)

// 		daysFree := 0
// 		for i := 0; i < 6; i++ {

// 			targetTime := time.Now()
// 			targetTime = time.Date(targetTime.Year(), targetTime.Month()+time.Month(i), 1, 0, 0, 0, 0, time.UTC)

// 			availability, err := recreation.GetAvailability(ctx, client, Park, targetTime)
// 			if err != nil {
// 				fmt.Println(err)
// 				return
// 			}
// 			for _, camp := range availability.Campsites {
// 				for _, avail := range camp.Availabilities {
// 					if avail == recreation.StateAvailable {
// 						daysFree++
// 					}
// 				}
// 			}
// 		}
// 		fmt.Println(daysFree)
// 	}
// }

// func GetDailyAvailabilities(ctx context.Context, client recreation.HTTPClient) {
// 	allAvailabilities := make(map[string][]time.Time)

// 	for i := 0; i < 6; i++ {

// 		targetTime := time.Now()
// 		targetTime = time.Date(targetTime.Year(), targetTime.Month()+time.Month(i), 1, 0, 0, 0, 0, time.UTC)

// 		availability, err := recreation.GetAvailability(ctx, client, Park, targetTime)
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}

// 		for campsiteID, camp := range availability.Campsites {
// 			for date, avail := range camp.Availabilities {
// 				if avail == recreation.StateAvailable {
// 					allAvailabilities[campsiteID] = append(allAvailabilities[campsiteID], date)
// 				}
// 			}

// 		}
// 	}
// 	campsiteAvailability, ok := allAvailabilities[Campsite]
// 	if !ok {
// 		fmt.Println("didn't find the campsite you specified, getrekt")
// 	}

// 	fmt.Printf("availabilities for the next five months at %s:%s\n", Park, Campsite)
// 	sort.Slice(campsiteAvailability, func(i, j int) bool {
// 		return campsiteAvailability[i].Before(campsiteAvailability[j])
// 	})
// 	for _, availableDate := range campsiteAvailability {
// 		fmt.Println(availableDate.Format("2006-01-02 Monday"))
// 	}
// }
