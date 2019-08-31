package gobms

import "net/http"

// Client is a general struct
type Client struct {
	RegionCode string
	RegionName string
	HTTPClient *http.Client
}

// QuickBook is the response from QUICKBOOK endpoint
type QuickBook struct {
	MovieData *MoviesData `json:"moviesData"`
	Cinemas   *Cinemas    `json:"cinemas"`
}

// MoviesData is the "moviesData" field from QUICKBOOK response
type MoviesData struct {
	BookMyShow struct {
		ArrEvents []Event `json:"arrEvents"`
	} `json:"BookMyShow"`
}

// Cinemas is the "cinemas" field from QUICKBOOK response
type Cinemas struct {
	BookMyShow struct {
		AIVN []RegionalCinemaDetail
	} `json:"BookMyShow"`
}

// RegionalCinemaDetail contains information regarding Theaters in a region
type RegionalCinemaDetail struct {
	IsATMOSEnabled string `json:"IsATMOSEnable"`
	Address        string `json:"VenueAddress"`
	Code           string `json:"VenueCode"`
	Latitude       string `json:"VenueLatitude"`
	Longitude      string `json:"VenueLongitude"`
	Name           string `json:"VenueName"`
	SubRegionCode  string `json:"VenueSubRegionCode"`
	SubRegionName  string `json:"VenueSubRegionName"`
}

// Event is a movie object
type Event struct {
	Title       string       `json:"EventTitle"`
	Code        string       `json:"EventCode"`
	ChildEvents []ChildEvent `json:"ChildEvents"`
}

// ChildEvent can be multiple children of an Event (Movie)
// containing different languages/dimensions events
type ChildEvent struct {
	Name      string `json:"EventName"`
	Language  string `json:"EventLanguage"`
	Dimension string `json:"EventDimension"`
	Code      string `json:"EventCode"`
	Synopsis  string `json:"EventSynopsis"`
	URL       string `json:"EventURL"`
}

// PopularCinemas is the GETPREFFEREDCINEMAS endpoint response
type PopularCinemas struct {
	Popular map[string]PrefferedCinemaDetail `json:"popular"`
}

// PrefferedCinemaDetail contains information regarding Theaters
type PrefferedCinemaDetail struct {
	VenueCode  string `json:"venueCode"`
	VenueName  string `json:"venueName"`
	RegionCode string `json:"regionCode"`
	RegionName string `json:"regionName"`
}

// Showtimes contains lists of shows returned by
// GETSHOWTIMESBYEVENTANDVENUE endpoint
type Showtimes struct {
	BookMyShow struct {
		ArrShows []Show `json:"arrShows"`
	} `json:"BookMyShow"`
}

// Show contains info for a particular movie show
type Show struct {
	SessionID       string `json:"SessionId"`
	ShowDateCode    string
	ShowTimeDisplay string
	Categories      []Categories `json:"Categories"`
}

// Categories for a show
type Categories struct {
	AreaCatCode           string
	PriceDesc             string
	CurPrice              string
	SeatsAvail            string
	SeatLayout            string
	MaxSeats              string
	PercentAvail          string
	PriceCode             string
	intCategoryMaxTickets string
	CategoryRange         string
}

// Regions from GETREGIONS endpoint
type Regions map[string][]RegionDetails

type RegionDetails struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}
