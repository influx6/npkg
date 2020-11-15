package nip

import (
	"bytes"
	"net"
	"net/http"
	"strings"

	"github.com/influx6/npkg/nerror"
)

//IpRange - a structure that holds the start and end of a range of ip addresses
type IpRange struct {
	Start net.IP
	End   net.IP
}

// InRange - check to see if a given ip address is within a range given
func InRange(r IpRange, ipAddress net.IP) bool {
	// strcmp type byte comparison
	if bytes.Compare(ipAddress, r.Start) >= 0 && bytes.Compare(ipAddress, r.End) < 0 {
		return true
	}
	return false
}

type PrivateSubnets []IpRange

func (subnet *PrivateSubnets) ParseRequestIP(r *http.Request) (net.Addr, error) {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || IsPrivateSubnet(*subnet, realIP) {
				// bad address, go to next
				continue
			}

			var parsedIP, parseErr = net.ResolveIPAddr("tcp", realIP.String())
			if parseErr != nil {
				return nil, nerror.WrapOnly(parseErr)
			}
			return parsedIP, nil
		}
	}

	var ip, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, nerror.WrapOnly(err)
	}

	var parsedIP, parseErr = net.ResolveIPAddr("tcp", ip)
	if parseErr != nil {
		return nil, nerror.WrapOnly(parseErr)
	}

	return parsedIP, nil
}

// IsPrivateSubnet - check to see if this ip is in a private subnet
func IsPrivateSubnet(subnet PrivateSubnets, ipAddress net.IP) bool {
	// my use case is only concerned with ipv4 atm
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		// iterate over all our ranges
		for _, r := range subnet {
			// check if this ip is in a private range
			if InRange(r, ipAddress) {
				return true
			}
		}
	}
	return false
}
