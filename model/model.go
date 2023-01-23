package model

import (
	"sync"
	"time"
)

// StravaActivity represents https://developers.strava.com/docs/reference/#api-models-SummaryActivity
type StravaActivity struct {
	ResourceState int `json:"resource_state"`
	Athlete       struct {
		ID            int `json:"id"`
		ResourceState int `json:"resource_state"`
	} `json:"athlete"`
	Name               string    `json:"name"`
	Distance           float64   `json:"distance"`
	MovingTime         int       `json:"moving_time"`
	ElapsedTime        int       `json:"elapsed_time"`
	TotalElevationGain float64   `json:"total_elevation_gain"`
	Type               string    `json:"type"`
	WorkoutType        any       `json:"workout_type"`
	ID                 int64     `json:"id"`
	ExternalID         string    `json:"external_id"`
	UploadID           int64     `json:"upload_id"`
	StartDate          time.Time `json:"start_date"`
	StartDateLocal     time.Time `json:"start_date_local"`
	Timezone           string    `json:"timezone"`
	UtcOffset          float64   `json:"utc_offset"`
	StartLatlng        any       `json:"start_latlng"`
	EndLatlng          any       `json:"end_latlng"`
	LocationCity       any       `json:"location_city"`
	LocationState      any       `json:"location_state"`
	LocationCountry    string    `json:"location_country"`
	AchievementCount   int       `json:"achievement_count"`
	KudosCount         int       `json:"kudos_count"`
	CommentCount       int       `json:"comment_count"`
	AthleteCount       int       `json:"athlete_count"`
	PhotoCount         int       `json:"photo_count"`
	Map                struct {
		ID              string `json:"id"`
		SummaryPolyline any    `json:"summary_polyline"`
		ResourceState   int    `json:"resource_state"`
	} `json:"map"`
	Trainer              bool    `json:"trainer"`
	Commute              bool    `json:"commute"`
	Manual               bool    `json:"manual"`
	Private              bool    `json:"private"`
	Flagged              bool    `json:"flagged"`
	GearID               string  `json:"gear_id"`
	FromAcceptedTag      bool    `json:"from_accepted_tag"`
	AverageSpeed         float64 `json:"average_speed"`
	MaxSpeed             float64 `json:"max_speed"`
	AverageCadence       float64 `json:"average_cadence"`
	AverageWatts         float64 `json:"average_watts"`
	WeightedAverageWatts int     `json:"weighted_average_watts"`
	Kilojoules           float64 `json:"kilojoules"`
	DeviceWatts          bool    `json:"device_watts"`
	HasHeartrate         bool    `json:"has_heartrate"`
	AverageHeartrate     float64 `json:"average_heartrate"`
	MaxHeartrate         float64 `json:"max_heartrate"`
	MaxWatts             int     `json:"max_watts"`
	PrCount              int     `json:"pr_count"`
	TotalPhotoCount      int     `json:"total_photo_count"`
	HasKudoed            bool    `json:"has_kudoed"`
	SufferScore          int     `json:"suffer_score"`
}

// SafeMap stores all activities. The number of each page is used as a key for the map V,
// so all activities from a single page can be easily retrieved. It is safe to use concurrently.
type SafeMap struct {
	mu sync.Mutex
	V  map[int][]StravaActivity
}

// Add adds the activities from a page to m.
func (m *SafeMap) Add(a []StravaActivity, pageNum int) {
	m.mu.Lock()
	m.V[pageNum] = a
	m.mu.Unlock()
}

// GetTotal returns the total number of activities found in m.
func (m *SafeMap) GetTotal() int {
	total := 0
	for _, activitiesOnPage := range m.V {
		total += len(activitiesOnPage)
	}
	return total
}
