package task

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultInputFile = "ip.txt"

var (
	// TestAll test all ip
	TestAll = false
	// IPFile is the filename of IP Rangs
	IPFile = defaultInputFile
	IPText string
)

func InitRandSeed() {
	rand.Seed(time.Now().UnixNano())
}

func isIPv4(ip string) bool {
	return strings.Contains(ip, ".")
}

func randIPEndWith(num byte) byte {
	if num == 0 { // For a single IP like /32
		return byte(0)
	}
	return byte(rand.Intn(int(num)))
}

type IPRanges struct {
	ips     []*net.IPAddr
	mask    string
	firstIP net.IP
	ipNet   *net.IPNet
}

func newIPRanges() *IPRanges {
	return &IPRanges{
		ips: make([]*net.IPAddr, 0),
	}
}

func (r *IPRanges) fixIP(ip string) string {
	// if it does not contain '/' It means that it is not an IP segment, but a single IP，So you need to add /32 /128 subnet mask
	if i := strings.IndexByte(ip, '/'); i < 0 {
		if isIPv4(ip) {
			r.mask = "/32"
		} else {
			r.mask = "/128"
		}
		ip += r.mask
	} else {
		r.mask = ip[i:]
	}
	return ip
}

func (r *IPRanges) parseCIDR(ip string) {
	var err error
	if r.firstIP, r.ipNet, err = net.ParseCIDR(r.fixIP(ip)); err != nil {
		log.Fatalln("ParseCIDR err", err)
	}
}

func (r *IPRanges) appendIPv4(d byte) {
	r.appendIP(net.IPv4(r.firstIP[12], r.firstIP[13], r.firstIP[14], d))
}

func (r *IPRanges) appendIP(ip net.IP) {
	r.ips = append(r.ips, &net.IPAddr{IP: ip})
}

// Return the minimum value and available number of the fourth segment ip
func (r *IPRanges) getIPRange() (minIP, hosts byte) {
	minIP = r.firstIP[15] & r.ipNet.Mask[3] // The minimum value of the fourth segment of IP

	// Get the number of hosts according to the subnet mask
	m := net.IPv4Mask(255, 255, 255, 255)
	for i, v := range r.ipNet.Mask {
		m[i] ^= v
	}
	total, _ := strconv.ParseInt(m.String(), 16, 32) // Total Available IPs
	if total > 255 {                                 // Correct the number of available IPs in the fourth paragraph
		hosts = 255
		return
	}
	hosts = byte(total)
	return
}

func (r *IPRanges) chooseIPv4() {
	minIP, hosts := r.getIPRange()
	for r.ipNet.Contains(r.firstIP) {
		if TestAll { // If it is speed test of all IPs
			for i := 0; i <= int(hosts); i++ { // Traversing the last segment of the IP from the minimum value to the maximum value
				r.appendIPv4(byte(i) + minIP)
			}
		} else { // Last segment of random IP 0.0.0.X
			r.appendIPv4(minIP + randIPEndWith(hosts))
		}
		r.firstIP[14]++ // 0.0.(X+1).X
		if r.firstIP[14] == 0 {
			r.firstIP[13]++ // 0.(X+1).X.X
			if r.firstIP[13] == 0 {
				r.firstIP[12]++ // (X+1).X.X.X
			}
		}
	}
}

func (r *IPRanges) chooseIPv6() {
	var tempIP uint8
	for r.ipNet.Contains(r.firstIP) {
		if r.mask != "/128" {
			r.firstIP[15] = randIPEndWith(255) // The last segment of the random IP
			r.firstIP[14] = randIPEndWith(255) // The last segment of the random IP
		}
		targetIP := make([]byte, len(r.firstIP))
		copy(targetIP, r.firstIP)
		r.appendIP(targetIP)
		for i := 13; i >= 0; i-- {
			tempIP = r.firstIP[i]
			r.firstIP[i] += randIPEndWith(255)
			if r.firstIP[i] >= tempIP {
				break
			}
		}
	}
}

func loadIPRanges() []*net.IPAddr {
	ranges := newIPRanges()
	if IPText != "" { // Get IP segment data from parameters
		IPs := strings.Split(IPText, ",")
		for _, IP := range IPs {
			ranges.parseCIDR(IP)
			if isIPv4(IP) {
				ranges.chooseIPv4()
			} else {
				ranges.chooseIPv6()
			}
		}
	} else { // Get IP segment data from file
		if IPFile == "" {
			IPFile = defaultInputFile
		}
		file, err := os.Open(IPFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			ranges.parseCIDR(scanner.Text())
			if isIPv4(scanner.Text()) {
				ranges.chooseIPv4()
			} else {
				ranges.chooseIPv6()
			}
		}
	}
	return ranges.ips
}
