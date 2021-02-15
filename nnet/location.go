package nnet

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/influx6/npkg/nerror"
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
	IP            string `json:"ip"`
	Street        string `json:"street"`
	City          string `json:"city"`
	State         string `json:"state"`
	Postal        string `json:"postal"`
	CountryCode   string `json:"country_code"`
	CountryName   string `json:"country_name"`
	RegionCode    string `json:"region_code"`
	RegionName    string `json:"region_name"`
	Zipcode       string `json:"zip_code"`
	Lat           string `json:"lat"`
	Long          string `json:"long"`
	MetroCode     string `json:"metro_code"`
	Timezone      string `json:"time_zone"`
	AreaCode      string `json:"area_code"`
	FromIP        string `json:"from_ip"`
	ToIP          string `json:"to_ip"`
	FromIPNumeric string `json:"from_ip_numeric"`
	ToIPNumeric   string `json:"to_ip_numeric"`
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

type FreeGeoipLocationService struct{}

func (f FreeGeoipLocationService) Get(address string) (Location, error) {
	var lt Location

	// Use freegeoip.net to get a JSON response
	// There is also /xml/ and /csv/ formats available
	var response, err = http.Get("https://freegeoip.net/json/" + address)
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
		return lt, nerror.WrapOnly(err)
	}

	return lt, nil
}
