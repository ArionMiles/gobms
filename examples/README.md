# Examples

This file contains a few small snippets demonstrating common functionality of the module.

You must first initialize the `bms.NewClient` with a `RegionCode` and `RegionName`.

## Get Your Region and Theater Information

After installing the CLI:

1. `RegionCode` and `RegionName`:

   ```
   gobms-cli listRegions > regions.txt
   ```

   You can search the output file for your region's name and code.

2. `TheaterCode`:

   From the region info you received in step 1, run:

   ```
   gobms-cli listTheaters --regionName=<your-region-name> --regionCode=<your-region-code>
   ```

```go
import "github.com/ArionMiles/gobms/pkg/bms"
```

```go

client, err := bms.NewClient("Mumbai", "MUMBAI")
if err != nil {
    // Handle Error
}
```

**NOTE:** Go through the [models.go](models.go) to find more attribute information.

## Get a list of movies currently showing

```go
qb, err := client.GetQuickbook("MT")
if err != nil {
    // Handle Error
}

for _, movie := range qb.MovieData.BookMyShow.ArrEvents {
    fmt.Println(movie.Title)
    for _, childEvent := range movie.ChildEvents {
        fmt.Println(childEvent.Language, childEvent.Dimension)
    }
}
```

## List popular theaters in your region

```go
cinemas, err := client.ListPreferredCinemas()
if err != nil {
	// Handle Error
}

for _, venue := range cinemas.Popular {
	fmt.Println(venue.VenueCode, " - ", venue.VenueName)
}
```

## Get all the theaters in your region

```go
theaters, err := client.ListTheaters()
if err != nil {
    // Handle error
}

for _, theater := range theaters {
    fmt.Println(theater.Name)
}
```

## Get a list of movie names currently showing

```go
movies, err := client.ListMovieTitles()
if err != nil {
    // Handle Error
}

for _, movie := range movies {
    fmt.Println(movie)
}
```

## Get a BookMyShow URL for a movie

```go
// Both date examples work
// date := time.Now()
date := time.Date(2019, time.December, 20, 0, 0, 0, 0, time.UTC)
movieURL, err := client.GetMovieURL("Star Wars: The Rise of Skywalker", "English", "IMAX 3D", date)
if err != nil {
    // Handle Error
}
fmt.Println(movieURL)

```

## Get an Event Code for a movie

```go
eventCode, err := obj.GetEventCode("Gunda", "Hindi", "IMAX 3D")
if err != nil {
    // Handle Error
}
fmt.Println(eventCode)
```

## Get showtimes for a movie

Showtimes are all the screenings of a movie in a particular theater (Venue) for a particular date.

Requires a `EventCode`, `VenueCode` and a date for the show

```go
date := time.Now()
shows, err := obj.GetShowtimes("ET00042069", "ABCD", date)
if err != nil {
    // Handle Error
}

for _, show := range shows {
    fmt.Println(show.SessionID)
    for _, category := range show.Categories {
        fmt.Println(category.CurPrice, category.SeatsAvail, category.SeatsMax)
    }
}
```

## Get all regions

```go
regions, err := bms.ListRegions()
if err != nil {
    // Handle Error
}

for _, region := range *regions {
    fmt.Println(region[0].Name, region[0].Code)
}
```
