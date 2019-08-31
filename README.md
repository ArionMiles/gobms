# GoBMS

An unofficial API wrapper for BookMyShow.com written in Go

## Installation
This API wrapper has been tested compatible with Go 1.11. 

It relies only on the Go standard library and has no external dependencies.

```
go get github.com/ArionMiles/gobms
```

## Usage

You must first initialize the `gobms.NewClient` with a `RegionCode` and `RegionName`.

`gobms.GetRegionList` returns a list of Regions which you can inspect to find the region name and code for your particular city.

```go
client := gobms.NewClient("Mumbai", "MUMBAI")
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

## Get list of popular theaters in your region
```go
cinemas, err := client.GetPreferredCinemas()
if err != nil {
	// Handle Error
}

for _, venue := range cinemas.Popular {
	fmt.Println(venue.VenueCode, " - ", venue.VenueName)
}
```

## Get all the theaters in your region
```go
theaters, err := client.GetTheaterData()
if err != nil {
    // Handle error
}

for _, theater := range theaters {
    fmt.Println(theater.Name)
}
```

## Get a list of movie names currently showing
```go
movies, err := client.GetMovieList()
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
regions, err := gobms.GetRegionList()
if err != nil {
    // Handle Error
}

for _, region := range *regions {
    fmt.Println(region[0].Name, region[0].Code)
}
```

## Contributions
For making feature requests, bug reports, open an issue.

All PRs are welcome, though, do create an issue first before making a PR.

## License

[MIT License](LICENSE)