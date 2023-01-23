package dl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ad-8/gobox/net"
	timex "github.com/ad-8/gobox/time"
	"github.com/ad-8/strava-dl-json/model"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

const (
	stravaOAuth              = "https://www.strava.com/oauth/token"
	stravaActivitiesEndpoint = "https://www.strava.com/api/v3/athlete/activities"
	activitiesPerPage        = 200 // 200 activities per page is a Strava limit
	currentlyFullPages       = 10  // TODO improve logic (e.g. store in file and update dynamically)
)

// TokenInfo represents the response that contains information about the Strava access token.
type TokenInfo struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	ExpiresHours int    `json:"-"`
	ExpiresMin   int    `json:"-"`
	ExpiresSec   int    `json:"-"`
}

// ParseTime parses a duration in seconds to a more human-readable representation in hours, minutes and seconds.
func (t *TokenInfo) ParseTime() {
	simpleTime, err := timex.NewSimpleTime(t.ExpiresIn)
	if err != nil {
		log.Fatal(err)
	}
	t.ExpiresHours = simpleTime.H
	t.ExpiresMin = simpleTime.M
	t.ExpiresSec = simpleTime.S
}

// Print prints when the token will expire.
func (t *TokenInfo) Print() {
	fmt.Printf("the token expires in %02d:%02d:%02d (will be automatically refreshed)\n\n",
		t.ExpiresHours, t.ExpiresMin, t.ExpiresSec)
}

// NewTokenInfo gets information about the access token - because the access token expires every 6 hours - and returns
// a *TokenInfo and nil if successful. Returns nil and the error if one occurs.
func NewTokenInfo(clientId, clientSecret, refreshToken string) (*TokenInfo, error) {
	params := map[string]any{
		"client_id":     clientId,
		"client_secret": clientSecret,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	}

	body, statusCode, err := net.MakePOSTRequest(stravaOAuth, params)
	if err != nil {
		return nil, fmt.Errorf("NewTokenInfo: error making request: %w", err)
	}

	tokenInfo := new(TokenInfo)
	if err := json.Unmarshal(body, tokenInfo); err != nil {
		return nil, errors.New(
			fmt.Sprintf("Error: %v\nstatus code is %d. cannot unmarshal this response:\n%v\n",
				err, statusCode, string(body)))
	}

	tokenInfo.ParseTime()

	return tokenInfo, nil
}

// StopFlag is a simple boolean flag safe for concurrent use.
type StopFlag struct {
	mu  sync.Mutex
	val bool
}

// true sets the flag to true.
func (f *StopFlag) true() {
	f.mu.Lock()
	if f.val == false {
		f.val = true
	}
	f.mu.Unlock()
}

// value returns the current value of the flag.
func (f *StopFlag) value() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.val
}

// AllActivities queries the Strava API for all activities of a user using an access token
// and returns all activities, sorted by date (the newest activity is the first element).
// The access token must have read_all scope.
func AllActivities(info TokenInfo) ([]model.StravaActivity, error) {
	allActivities := model.SafeMap{V: make(map[int][]model.StravaActivity)}
	var stop StopFlag

	var wg sync.WaitGroup
	for pageNum := 1; ; pageNum++ {
		// TODO improve logic (maybe use channels)
		// without time.Sleep there will be ~3k goroutines making requests before stop is set to true
		if pageNum > currentlyFullPages {
			time.Sleep(250 * time.Millisecond)
		}
		if stop.value() == true {
			break
		}
		wg.Add(1)
		go getPage(info.AccessToken, pageNum, &allActivities, &wg, &stop)
	}
	wg.Wait()

	return sortActivities(&allActivities), nil
}

// getPage queries the Strava API for all activities on the specified page.
func getPage(accessToken string, pageNum int, m *model.SafeMap, wg *sync.WaitGroup, stop *StopFlag) {
	var activitiesOnPage []model.StravaActivity
	body, err := requestActivitiesFromPage(accessToken, pageNum)
	if err != nil {
		log.Fatalf("getPage: %v", err)
	}

	if string(body) == "[]" {
		stop.true()
		wg.Done()
		return
	}

	if err := json.Unmarshal(body, &activitiesOnPage); err != nil {
		fmt.Println(string(body))
		log.Fatal(err)
	}

	m.Add(activitiesOnPage, pageNum)

	wg.Done()
}

// requestActivitiesFromPage makes an HTTP GET request to get all user activities from one page and
// returns the data and nil if successful. Because only a maximum of 200 activities can be requested
// at once, one may need to run this function multiple times while incrementing pageNum from 1 to n, until the response
// data equals "[]", so the un-marshaled slice of type StravaActivity is empty.
func requestActivitiesFromPage(accessToken string, pageNum int) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, stravaActivitiesEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + accessToken},
	}
	q := req.URL.Query()
	q.Add("page", strconv.Itoa(pageNum))
	q.Add("per_page", strconv.Itoa(activitiesPerPage))
	req.URL.RawQuery = q.Encode()

	resp, _, err := net.MakeGETRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// sortActivities first flattens and then sorts (by date, the newest activity first) all activities in m.V.
func sortActivities(m *model.SafeMap) []model.StravaActivity {
	var all []model.StravaActivity

	for _, page := range m.V {
		for _, activity := range page {
			all = append(all, activity)
		}
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].StartDateLocal.After(all[j].StartDateLocal)
	})

	return all
}
