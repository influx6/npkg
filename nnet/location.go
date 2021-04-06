package nnet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/influx6/npkg/nerror"
	"github.com/influx6/npkg/nunsafe"
	"xojoc.pw/useragent"
)

type Agent struct {
	*useragent.UserAgent
	Browser string
}

func ParseAgent(agentString string) (*Agent, error) {
	var agentInfo = useragent.Parse(agentString)
	if agentInfo == nil {
		return nil, nerror.New("failed to parse agent")
	}
	return &Agent{
		Browser:   agentInfo.Name,
		UserAgent: agentInfo,
	}, nil
}

type Locations []Location

func (l Locations) Find(ip string) (Location, error) {
	for _, lt := range l {
		if lt.IPIsInLocation(ip) {
			return lt, nil
		}
	}
	return Location{}, nerror.New("not found")
}

type Location struct {
	Type          string  `json:"type_"`
	Street        string  `json:"street"`
	ContinentCode string  `json:"continent_code"`
	ContinentName string  `json:"continent_name"`
	City          string  `json:"city"`
	State         string  `json:"state"`
	Postal        string  `json:"postal"`
	Zip           string  `json:"zip"`
	Zipcode       string  `json:"zip_code"`
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	RegionCode    string  `json:"region_code"`
	RegionName    string  `json:"region_name"`
	Timezone      string  `json:"time_zone"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	MetroCode     string  `json:"metro_code"`
	AreaCode      string  `json:"area_code"`
	FromIP        string  `json:"from_ip"`
	ToIP          string  `json:"to_ip"`
	FromIPNumeric string  `json:"from_ip_numeric"`
	ToIPNumeric   string  `json:"to_ip_numeric"`
}

func (l Location) IPIsInLocation(ip string) bool {
	var fromIP = net.ParseIP(l.FromIP)
	if fromIP == nil {
		return false
	}
	var toIP = net.ParseIP(l.ToIP)
	if toIP == nil {
		return false
	}
	var targetIP = net.ParseIP(ip)
	if toIP == nil {
		return false
	}
	return IsTargetBetween(targetIP, fromIP, toIP)
}

type LocationService interface {
	Get(ip string) (Location, error)
}

// DudLocationService returns a default unknown location with
// provided address as ip.
type DudLocationService struct{}

func (f DudLocationService) Get(address string) (Location, error) {
	var lt Location
	lt.City = "Unknown"
	lt.State = "Unknown"
	lt.CountryName = "Unknown"
	lt.RegionCode = "Unknown"
	lt.RegionName = "Unknown"
	lt.CountryCode = "Unknown"
	lt.Zipcode = "00000"
	return lt, nil
}

type CampIPService struct {
	Addr string
}

func (f CampIPService) Get(address string) (Location, error) {
	var lt Location

	// Use campip service to get a JSON response
	var response, err = http.Get(fmt.Sprintf("%s/%s", f.Addr, address))
	if err != nil {
		return lt, nerror.WrapOnly(err)
	}
	defer response.Body.Close()

	// response.Body() is a reader type. We have
	// to use ioutil.ReadAll() to read the data
	// in to a byte slice(string)
	var body, berr = ioutil.ReadAll(response.Body)
	if berr != nil {
		return lt, nerror.WrapOnly(berr)
	}

	// Unmarshal the JSON byte slice to a GeoIP struct
	err = json.Unmarshal(body, &lt)
	if err != nil {
		return lt, nerror.WrapOnly(err).Add("body", nunsafe.Bytes2String(body))
	}

	return lt, nil
}

type IPStackService struct {
	Token string
}

func (f IPStackService) Get(address string) (Location, error) {
	var lt Location

	// Use freegeoip.net to get a JSON response
	// There is also /xml/ and /csv/ formats available
	var response, err = http.Get(fmt.Sprintf("http://api.ipstack.com/%s?access_key=%s", address, f.Token))
	if err != nil {
		return lt, nerror.WrapOnly(err)
	}
	defer response.Body.Close()

	// response.Body() is a reader type. We have
	// to use ioutil.ReadAll() to read the data
	// in to a byte slice(string)
	var body, berr = ioutil.ReadAll(response.Body)
	if berr != nil {
		return lt, nerror.WrapOnly(berr)
	}

	// Unmarshal the JSON byte slice to a GeoIP struct
	err = json.Unmarshal(body, &lt)
	if err != nil {
		return lt, nerror.WrapOnly(err).Add("body", nunsafe.Bytes2String(body))
	}

	return lt, nil
}
