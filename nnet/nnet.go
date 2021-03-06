package nnet

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/influx6/npkg/nerror"
	"github.com/influx6/npkg/nunsafe"
)

var (
	noTime = time.Time{}
)

// TimedConn implements a wrapper around a net.Conn which guards giving connection
// with appropriate read/write timeout.
type TimedConn struct {
	net.Conn
	readTimeout  time.Duration
	writeTimeout time.Duration
}

// NewTimedConn returns a new instance of a TimedConn.
func NewTimedConn(conn net.Conn, rd time.Duration, wd time.Duration) *TimedConn {
	return &TimedConn{
		Conn:         conn,
		readTimeout:  rd,
		writeTimeout: wd,
	}
}

// Write calls the underline connection read with provided timeout.
func (c *TimedConn) Write(b []byte) (int, error) {
	var writeErr = c.Conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	if writeErr != nil {
		return 0, writeErr
	}

	var writeCount, err = c.Conn.Write(b)
	if err != nil {
		return writeCount, err
	}

	var resetErr = c.Conn.SetWriteDeadline(noTime)
	if resetErr != nil {
		return writeCount, resetErr
	}

	return writeCount, nil
}

// Read calls the underline connection read with provided timeout.
func (c *TimedConn) Read(b []byte) (int, error) {
	var readErr = c.Conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	if readErr != nil {
		return 0, readErr
	}

	var readCount, err = c.Conn.Read(b)
	if err != nil {
		return readCount, err
	}

	var resetErr = c.Conn.SetReadDeadline(noTime)
	if resetErr != nil {
		return readCount, resetErr
	}

	return readCount, nil
}

// GetAddr takes the giving address string and if it has no ip or use the
// zeroth ip format, then modifies the ip with the current systems ip.
func GetAddr(addr string) string {
	if addr == "" {
		if real, err := GetMainIP(); err == nil {
			return real + ":0"
		}
	}

	ip, port, err := net.SplitHostPort(addr)
	if err == nil && (ip == "" || ip == "0.0.0.0") {
		if realIP, err := GetMainIP(); err == nil {
			return net.JoinHostPort(realIP, port)
		}
	}

	return addr
}

var (
	zeros = regexp.MustCompile("^0+")
)

// ResolveAddr returns an appropriate address by validating the
// presence of the ip and port, if non is found, it uses the default
// 0.0.0.0 address and assigns a port if non is found.
func ResolveAddr(addr string) string {
	var scheme string
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		uri, err := url.Parse(addr)
		if err != nil {
			return "0.0.0.0:" + strconv.Itoa(FreePort())
		}

		if strings.Contains(uri.Host, ":") {
			sub := strings.Index(uri.Host, ":")
			uri.Host = uri.Host[0:sub]
		}

		scheme = uri.Scheme
		host = uri.Host
		port = uri.Port()
	}

	if host == "" {
		host = "0.0.0.0"
	}

	if port == "" || zeros.MatchString(port) {
		port = strconv.Itoa(FreePort())
	}

	if scheme == "" {
		return host + ":" + port
	}

	return scheme + "://" + host + ":" + port
}

func IPDotNotation2LongNotation(ipAddr string) (uint32, error) {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return 0, errors.New("wrong ipAddr format")
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip), nil
}

func IPLongNotation2IPFromString(ipLong string) (net.IP, error) {
	var parsedVal, parseErr = strconv.ParseUint(ipLong, 10, 32)
	if parseErr != nil {
		return nil, nerror.WrapOnly(parseErr)
	}
	return IPLongNotation2IP(uint32(parsedVal)), nil
}

func IPLongNotation2IP(ipLong uint32) net.IP {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	return net.IP(ipByte)
}

func IPLongNotation2DotNotation(ipLong uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip := net.IP(ipByte)
	return ip.String()
}

func IntToIP(ip uint32) string {
	result := make(net.IP, 4)
	result[0] = byte(ip)
	result[1] = byte(ip >> 8)
	result[2] = byte(ip >> 16)
	result[3] = byte(ip >> 24)
	return result.String()
}

func IsTargetBetweenUsingCDIR(cdir string, to net.IP) (bool, error) {
	if to == nil {
		return false, nerror.New("target cant be nil")
	}

	var _, subnet, subErr = net.ParseCIDR(cdir)
	if subErr != nil {
		return false, nerror.WrapOnly(subErr)
	}

	return subnet.Contains(to), nil
}

func IsTargetBetween(target net.IP, from net.IP, to net.IP) bool {
	if from == nil || to == nil || target == nil {
		return false
	}

	from16 := from.To16()
	to16 := to.To16()
	test16 := target.To16()
	if from16 == nil || to16 == nil || test16 == nil {
		return false
	}

	if bytes.Compare(test16, from16) >= 0 && bytes.Compare(test16, to16) <= 0 {
		return true
	}
	return false
}

// FreePort returns a random free port from the underline system for use in a network socket.
func FreePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
		return 0
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port
}

//==============================================================================

// UpgradeConnToTLS upgrades the giving tcp connection to use a tls based connection
// encrypted by the giving tls.Config.
func UpgradeConnToTLS(conn net.Conn, cm *tls.Config) (net.Conn, error) {
	if cm == nil {
		return conn, nil
	}

	if cm.ServerName == "" {
		h, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
		cm.ServerName = h
	}

	tlsConn := tls.Client(conn, cm)

	if err := tlsConn.Handshake(); err != nil {
		return conn, err
	}

	return tlsConn, nil
}

//================================================================================

//LoadTLS loads a tls.Config from a key and cert file path
func LoadTLS(cert string, key string, ca string) (*tls.Config, error) {
	var config *tls.Config
	config.MinVersion = tls.VersionTLS12
	config.Certificates = make([]tls.Certificate, 1)

	c, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	c.Leaf, err = x509.ParseCertificate(c.Certificate[0])
	if err != nil {
		return nil, err
	}

	if ca != "" {
		rootPEM, err := ioutil.ReadFile(ca)
		if err != nil {
			return nil, err
		}

		if rootPEM == nil {
			return nil, errors.New("Empty perm file")
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(rootPEM) {
			return nil, errors.New("Failed to append perm file data")
		}

		config.RootCAs = pool
	}

	config.Certificates[0] = c
	return config, nil
}

// TLSVersion returns a string version number based on the tls version int.
func TLSVersion(ver uint16) string {
	switch ver {
	case tls.VersionTLS10:
		return "1.0"
	case tls.VersionTLS11:
		return "1.1"
	case tls.VersionTLS12:
		return "1.2"
	}
	return fmt.Sprintf("Unknown [%x]", ver)
}

// TLSCipher returns a cipher string version based on the supplied hex value.
func TLSCipher(cs uint16) string {
	switch cs {
	case 0x0005:
		return "TLS_RSA_WITH_RC4_128_SHA"
	case 0x000a:
		return "TLS_RSA_WITH_3DES_EDE_CBC_SHA"
	case 0x002f:
		return "TLS_RSA_WITH_AES_128_CBC_SHA"
	case 0x0035:
		return "TLS_RSA_WITH_AES_256_CBC_SHA"
	case 0xc007:
		return "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA"
	case 0xc009:
		return "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA"
	case 0xc00a:
		return "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA"
	case 0xc011:
		return "TLS_ECDHE_RSA_WITH_RC4_128_SHA"
	case 0xc012:
		return "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA"
	case 0xc013:
		return "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA"
	case 0xc014:
		return "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA"
	case 0xc02f:
		return "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"
	case 0xc02b:
		return "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"
	case 0xc030:
		return "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
	case 0xc02c:
		return "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"
	}
	return fmt.Sprintf("Unknown [%x]", cs)
}

//==============================================================================

// MakeListener returns a new net.Listener requests.
func MakeListener(protocol string, addr string, conf *tls.Config) (net.Listener, error) {
	var l net.Listener
	var err error

	if conf != nil {
		l, err = tls.Listen(protocol, addr, conf)
	} else {
		l, err = net.Listen(protocol, addr)
	}

	if err != nil {
		return nil, err
	}

	return l, nil
}

// TCPListener returns a new net.Listener requests.
func TCPListener(addr string, conf *tls.Config) (net.Listener, error) {
	var l net.Listener
	var err error

	if conf != nil {
		l, err = tls.Listen("tcp", addr, conf)
	} else {
		l, err = net.Listen("tcp", addr)
	}

	if err != nil {
		return nil, err
	}

	return NewKeepAliveListener(l), nil
}

//NewHTTPServer returns a new http.Server using the provided listener
func NewHTTPServer(l net.Listener, handle http.Handler, c *tls.Config) (*http.Server, net.Listener, error) {
	s := &http.Server{
		Addr:           l.Addr().String(),
		Handler:        handle,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      c,
	}

	log.Printf("Serving http connection on: %+q\n", s.Addr)
	go func() {
		if err := s.Serve(l); err != nil {
			log.Fatal(err)
		}
	}()

	return s, l, nil
}

type keepAliveListener struct {
	net.Listener
}

// NewKeepAliveListener returns a new net.Listener from underline net.TCPListener
// where produced net.Conns respect keep alive regulations.
func NewKeepAliveListener(tl net.Listener) net.Listener {
	return &keepAliveListener{
		Listener: tl,
	}
}

func (kl *keepAliveListener) Accept() (net.Conn, error) {
	var tl, isTL = kl.Listener.(*net.TCPListener)
	if !isTL {
		return kl.Listener.Accept()
	}

	conn, err := tl.AcceptTCP()
	if err != nil {
		return nil, err
	}

	if err := conn.SetKeepAlive(true); err != nil {
		return conn, err
	}
	if err := conn.SetKeepAlivePeriod(2 * time.Minute); err != nil {
		return conn, err
	}

	return conn, nil
}

// NewConn returns a tls.Conn object from the provided parameters.
func NewConn(protocol string, addr string) (net.Conn, error) {
	newConn, err := net.Dial(protocol, addr)
	if err != nil {
		return nil, err
	}

	return newConn, nil
}

// TLSConn returns a tls.Conn object from the provided parameters.
func TLSConn(protocol string, addr string, conf *tls.Config) (*tls.Conn, error) {
	newTLS, err := tls.Dial(protocol, addr, conf)
	if err != nil {
		return nil, err
	}

	return newTLS, nil
}

// TLSFromConn returns a new tls.Conn using the address and the certicates from
// the provided *tls.Conn.
func TLSFromConn(tl *tls.Conn, addr string) (*tls.Conn, error) {
	var conf *tls.Config

	if err := tl.Handshake(); err != nil {
		return nil, err
	}

	state := tl.ConnectionState()
	pool := x509.NewCertPool()

	for _, v := range state.PeerCertificates {
		pool.AddCert(v)
	}

	conf = &tls.Config{RootCAs: pool}
	newTLS, err := tls.Dial("tcp", addr, conf)
	if err != nil {
		return nil, err
	}

	return newTLS, nil
}

// ProxyHTTPRequest copies a http request from a src net.Conn connection to a
// destination net.Conn.
func ProxyHTTPRequest(src net.Conn, dest net.Conn) error {
	reader := bufio.NewReader(src)
	req, err := http.ReadRequest(reader)
	if err != nil {
		return err
	}

	if req == nil {
		return errors.New("No Request Read")
	}

	if err = req.Write(dest); err != nil {
		return err
	}

	resread := bufio.NewReader(dest)
	res, err := http.ReadResponse(resread, req)
	if err != nil {
		return err
	}

	if res != nil {
		return errors.New("No Response Read")
	}

	return res.Write(src)
}

// hop headers, These are removed when sent to the backend
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html.
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

// ConnToHTTP proxies a requests from a net.Conn to a destination request, writing
// the response to provided ResponseWriter.
func ConnToHTTP(src net.Conn, destReq *http.Request, destRes http.ResponseWriter) error {
	reader := bufio.NewReader(src)
	req, err := http.ReadRequest(reader)
	if err != nil {
		return err
	}

	destReq.Method = req.Method

	for key, val := range req.Header {
		destReq.Header.Set(key, strings.Join(val, ","))
	}

	for _, v := range hopHeaders {
		destReq.Header.Del(v)
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return err
	}

	//add us to the proxy list or makeone
	hops, ok := req.Header["X-Forwarded-For"]
	if ok {
		ip = strings.Join(hops, ",") + "," + ip
	}

	destReq.Header.Set("X-Forwarded-For", ip)

	var buf bytes.Buffer
	if req.Body != nil {
		io.Copy(&buf, req.Body)
	}

	if buf.Len() > 0 {
		destReq.Body = ioutil.NopCloser(&buf)
		destReq.ContentLength = int64(buf.Len())
	}

	res, err := http.DefaultClient.Do(destReq)
	if err != nil {
		return err
	}

	for k, v := range res.Header {
		destRes.Header().Add(k, strings.Join(v, ","))
	}

	return res.Write(destRes)
}

// HTTPToConn proxies a src Request to a net.Con connection and writes back
// the response to the src Response.
func HTTPToConn(srcReq *http.Request, srcRes http.ResponseWriter, dest net.Conn) error {
	if err := srcReq.Write(dest); err != nil {
		return err
	}

	resRead := bufio.NewReader(dest)
	res, err := http.ReadResponse(resRead, srcReq)
	if err != nil {
		return err
	}

	for key, val := range res.Header {
		srcRes.Header().Set(key, strings.Join(val, ","))
	}

	srcRes.WriteHeader(res.StatusCode)

	return res.Write(srcRes)
}

//==============================================================================

// GetClustersFriends returns a giving set of routes from the provided port number.
func GetClustersFriends(clusterPort int, routes []*url.URL) ([]*url.URL, error) {
	var cleanRoutes []*url.URL
	cport := strconv.Itoa(clusterPort)

	selfIPs, err := GetInterfaceIPs()
	if err != nil {
		return nil, err
	}

	for _, r := range routes {
		host, port, err := net.SplitHostPort(r.Host)
		if err != nil {
			return nil, err
		}

		ips, err := GetURLIP(host)
		if err != nil {
			return nil, err
		}

		if cport == port && IsIPInList(selfIPs, ips) {
			continue
		}

		cleanRoutes = append(cleanRoutes, r)
	}

	return cleanRoutes, nil
}

// GetMainIP returns the giving system IP by attempting to connect to a imaginary
// ip and returns the giving system ip.
func GetMainIP() (string, error) {
	udp, err := net.DialTimeout("udp", "8.8.8.8:80", 1*time.Millisecond)
	if udp == nil {
		return "", err
	}

	defer udp.Close()

	localAddr := udp.LocalAddr().String()
	ip, _, _ := net.SplitHostPort(localAddr)

	return ip, nil
}

// GetExternalIP returns the actual internal external IP of the
// calling system.
func GetExternalIP() (string, error) {
	var response, err = http.Get("http://ipv4bot.whatismyipaddress.com")
	if err != nil {
		return "", nerror.WrapOnly(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", nerror.WrapOnly(err)
	}

	return nunsafe.Bytes2String(body), nil
}

// GetMainIPByInterface returns the giving ip of the current system by looping
// through all interfaces returning the first ipv4 found that is not on the
// loopback interface.
func GetMainIPByInterface() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()

			if ip == nil {
				continue
			}

			return ip.String(), nil
		}
	}

	return "", errors.New("Error: No network connection found")
}

// IsIPInList returns true/false if the giving ip list were equal.
func IsIPInList(list1 []net.IP, list2 []net.IP) bool {
	for _, ip1 := range list1 {
		for _, ip2 := range list2 {
			if ip1.Equal(ip2) {
				return true
			}
		}
	}
	return false
}

// GetURLIP returns a giving ip addres if the ip string is not an ip address.
func GetURLIP(ipStr string) ([]net.IP, error) {
	ipList := []net.IP{}

	ip := net.ParseIP(ipStr)
	if ip != nil {
		ipList = append(ipList, ip)
		return ipList, nil
	}

	hostAddr, err := net.LookupHost(ipStr)
	if err != nil {
		return nil, err
	}

	for _, addr := range hostAddr {
		ip = net.ParseIP(addr)
		if ip != nil {
			ipList = append(ipList, ip)
		}
	}

	return ipList, nil
}

// GetInterfaceIPs returns the list of IP of the giving interfaces found in the
// system.
func GetInterfaceIPs() ([]net.IP, error) {
	var localIPs []net.IP

	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return nil, errors.New("Error getting self referencing addr")
	}

	for i := 0; i < len(interfaceAddr); i++ {
		interfaceIP, _, _ := net.ParseCIDR(interfaceAddr[i].String())
		if net.ParseIP(interfaceIP.String()) != nil {
			localIPs = append(localIPs, interfaceIP)
		} else {
			err = errors.New("Error getting self referencing addr")
		}
	}

	if err != nil {
		return nil, err
	}

	return localIPs, nil
}

// CopyUDPAddr returns a new copy of a giving UDPAddr.
func CopyUDPAddr(addr *net.UDPAddr) *net.UDPAddr {
	newAddr := new(net.UDPAddr)
	*newAddr = *addr
	newAddr.IP = make(net.IP, len(addr.IP))
	copy(newAddr.IP, addr.IP)
	return newAddr
}

// Helper to move from float seconds to time.Duration
func secondsToDuration(seconds float64) time.Duration {
	ttl := seconds * float64(time.Second)
	return time.Duration(ttl)
}

// Ascii numbers 0-9
const (
	asciiZero = 48
	asciiNine = 57
)

// parseSize expects decimal positive numbers. We
// return -1 to signal error
func parseSize(d []byte) (n int) {
	if len(d) == 0 {
		return -1
	}
	for _, dec := range d {
		if dec < asciiZero || dec > asciiNine {
			return -1
		}
		n = n*10 + (int(dec) - asciiZero)
	}
	return n
}

// parseInt64 expects decimal positive numbers. We
// return -1 to signal error
func parseInt64(d []byte) (n int64) {
	if len(d) == 0 {
		return -1
	}
	for _, dec := range d {
		if dec < asciiZero || dec > asciiNine {
			return -1
		}
		n = n*10 + (int64(dec) - asciiZero)
	}
	return n
}
