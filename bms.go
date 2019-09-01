package gobms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const baseURL = "https://in.bookmyshow.com/serv/getData"
const showtimesURL = "https://in.bookmyshow.com/buytickets/%v-%v/movie-%v-%v/%v"

// NewClient instantiates a Client struct with a Region Name and Region Code
//
// Example:
// 	c := bms.NewClient("MUMBAI", "Mumbai")
func NewClient(RegionCode, RegionName string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	var cookies []*http.Cookie
	cookieStr := fmt.Sprintf("|Code=%v|text=%v|", RegionCode, RegionName)
	urlencodedcookie := url.QueryEscape(cookieStr)
	cookie := &http.Cookie{
		Name:  "Rgn",
		Value: urlencodedcookie,
	}

	cookies = append(cookies, cookie)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	jar.SetCookies(u, cookies)

	httpclient := &http.Client{
		Jar: jar,
	}

	client := &Client{
		RegionCode: RegionCode,
		RegionName: RegionName,
		HTTPClient: httpclient,
	}

	return client, nil

}

func (c *Client) get(endpoint string, params map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}
	parameters := req.URL.Query()
	parameters.Add("cmd", endpoint)

	if params != nil {
		for param, value := range params {
			parameters.Add(param, value)
		}
	}
	req.URL.RawQuery = parameters.Encode()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetQuickbook returns a slice of Event objects
//
// eventType can be:
// 	MT(Movies), CT(Concerts or Misc Events), PL(Plays), SP(Sports)
//
// As of right now, the package only supports MT (Movies)
func (c *Client) GetQuickbook(eventType string) (*QuickBook, error) {
	params := map[string]string{
		"type": eventType,
	}

	resp, err := c.get("QUICKBOOK", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var qb QuickBook
	err = json.Unmarshal([]byte(jsonData), &qb)
	if err != nil {
		return nil, err
	}

	return &qb, nil
}

// GetPreferredCinemas returns a list of popular cinemas in your region
func (c *Client) GetPreferredCinemas() (*PopularCinemas, error) {
	resp, err := c.get("GETPREFERREDCINEMAS", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cinemas PopularCinemas
	json.Unmarshal([]byte(jsonData), &cinemas)

	return &cinemas, nil
}

// GetMovieList returns a slice of movie titles
func (c *Client) GetMovieList() ([]string, error) {
	qb, err := c.GetQuickbook("MT")
	if err != nil {
		return nil, err
	}

	var movieList []string
	for _, movie := range qb.MovieData.BookMyShow.ArrEvents {
		movieList = append(movieList, movie.Title)
	}

	return movieList, nil
}

// GetMovieURL returns a URL to a movie's main BMS page.
//
// Language can be:
// 	Hindi, English, Tamil, Telugu, Kannada, Gujarati or anything else BookMyShow supports
//
// Dimensions can be:
//  2D, 2D 4DX, 3D, 3D 4DX, IMAX 3D or any other that BookMyShow supports
func (c *Client) GetMovieURL(movieName, language, dimension string, date time.Time) (*url.URL, error) {
	// Date shouldn't be string
	// Should return nil if not found
	qb, err := c.GetQuickbook("MT")
	if err != nil {
		return nil, err
	}

	for _, movie := range qb.MovieData.BookMyShow.ArrEvents {
		if movie.Title == movieName {
			for _, childEvent := range movie.ChildEvents {
				if childEvent.Language == language && childEvent.Dimension == dimension {
					movieURL, err := c.buildMovieURL(&childEvent, date)
					if err != nil {
						return nil, err
					}
					return movieURL, nil
				}
			}
		}
	}

	err = goBMSError{
		fmt.Sprintf("%v not found!", movieName),
	}

	return nil, err
}

// GetEventCode returns an Event code for a given movie
func (c *Client) GetEventCode(movieName, language, dimension string) (string, error) {
	qb, err := c.GetQuickbook("MT")
	if err != nil {
		return "", err
	}

	for _, movie := range qb.MovieData.BookMyShow.ArrEvents {
		if movie.Title == movieName {
			for _, childEvent := range movie.ChildEvents {
				if childEvent.Language == language && childEvent.Dimension == dimension {
					return childEvent.Code, nil
				}
			}
		}
	}

	err = goBMSError{
		fmt.Sprintf("Event Code for %v not found!", movieName),
	}

	return "", err
}

// GetShowtimes returns a slice of Showtimes (ArrShows) for a movie at a particular
// venue (theater) on a particular date
func (c *Client) GetShowtimes(eventCode, venueCode string, date time.Time) ([]Show, error) {
	// Date shouldn't be string
	// Should return error/nil if no shows found
	formattedDate := date.Format("20060102")
	params := map[string]string{
		"f":  "json",
		"dc": formattedDate,
		"vc": venueCode,
		"ec": eventCode,
	}

	resp, err := c.get("GETSHOWTIMESBYEVENTANDVENUE", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var showtimes Showtimes
	json.Unmarshal([]byte(jsonData), &showtimes)
	if err != nil {
		return nil, err
	}

	if len(showtimes.BookMyShow.ArrShows) > 0 {
		return showtimes.BookMyShow.ArrShows, nil
	}
	err = goBMSError{
		"No Shows Found!",
	}
	return nil, err
}

func (c *Client) buildMovieURL(event *ChildEvent, date time.Time) (*url.URL, error) {
	formattedDate := date.Format("20060102")
	city := strings.Replace(c.RegionName, " ", "-", -1)
	formattedMovieURL := fmt.Sprintf(showtimesURL, event.URL, city, city, event.Code, formattedDate)
	movieURL, err := url.Parse(formattedMovieURL)
	if err != nil {
		return nil, err
	}
	return movieURL, nil
}

// GetTheaterData returns slice of Cinemas in your region
func (c *Client) GetTheaterData() ([]RegionalCinemaDetail, error) {
	qb, err := c.GetQuickbook("MT")
	if err != nil {
		return nil, err
	}
	return qb.Cinemas.BookMyShow.AIVN, nil
}

// GetShowtimeURL returns booking URL for a specific Show
func GetShowtimeURL(venueCode, sessionID string) (*url.URL, error) {
	formattedURL := fmt.Sprintf("https://in.bookmyshow.com/booktickets/%v/%v", venueCode, sessionID)
	showtimeURL, err := url.Parse(formattedURL)
	if err != nil {
		return nil, err
	}
	return showtimeURL, nil
}

// GetRegionList returns Region list including the name, code and alias of a region
// from GETREGIONS Endpoint
//
// Example:
// 		regionList, err := bms.GetRegionList()
// 		if err != nil {
// 			log.Print(err)
// }
//		for _, rgn := range regionList {
//   		fmt.Println(rgn[0].Name)
//		}
func GetRegionList() (*Regions, error) {
	resp, err := http.Get(baseURL + "?cmd=GETREGIONS")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := strings.SplitN(string(body), "=", 2)[1]
	regionLst := strings.SplitN(data, ";", 2)[0]
	var regions Regions
	json.Unmarshal([]byte(regionLst), &regions)
	return &regions, nil
}
