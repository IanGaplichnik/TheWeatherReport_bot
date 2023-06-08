package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(event Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
	Location
)

type Event struct {
	Type     Type
	Text     string
	Meta     interface{}
	Location *Coordinates
}

type Coordinates struct {
	Longitude float32
	Latitude  float32
}

type CityData struct {
	CityName  string
	Latitude  float32
	Longitude float32
}
