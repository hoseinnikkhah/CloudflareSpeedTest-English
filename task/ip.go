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

// If it is a single IP, add the subnet mask, otherwise, get the subnet mask (r.mask)
func (r *IPRanges) fixIP(ip string) string {
	// If it does not contain '/', it means that it is not an IP segment, but a separate IP, so /32 /128 subnet mask needs to be added
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

// Parse the IP segment to get IP, IP range, subnet mask
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
	minIP = r.firstIP[15] & r.ipNet.Mask[3] // IP 第四段最小值

	// Get the number of hosts according to the subnet mask
	m := net.IPv4Mask(255, 255, 255, 255)
	for i, v := range r.ipNet.Mask {
		m[i] ^= v
	}
	total, _ := strconv.ParseInt(m.String(), 16, 32) // Total available IPs
	if total > 255 {                                 // Correct the number of available IPs in the fourth paragraph
		hosts = 255
		return
	}
	hosts = byte(total)
	return
}

func (r *IPRanges) chooseIPv4() {
	if r.mask == "/32" { // A single IP does not need to be random, just add itself directly
		r.appendIP(r.firstIP)
	} else {
		minIP, hosts := r.getIPRange()    // Return the minimum value and available number of the fourth segment IP
		for r.ipNet.Contains(r.firstIP) { // As long as the IP does not exceed the range of the IP network segment, continue to cycle randomly
			if TestAll { // If it is speed test all IP
				for i := 0; i <= int(hosts); i++ { // Traversing the last segment of IP from the minimum value to the maximum value
					r.appendIPv4(byte(i) + minIP)
				}
			} else { // last segment of random IP 0.0.0.X
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
}

func (r *IPRanges) chooseIPv6() {
	if r.mask == "/128" { // A single IP does not need to be random, just add itself directly
		r.appendIP(r.firstIP)
	} else {
		var tempIP uint8                  // Temporary variable, used to record the value of the previous bit
		for r.ipNet.Contains(r.firstIP) { // As long as the IP does not exceed the range of the IP network segment, continue to cycle randomly
			r.firstIP[15] = randIPEndWith(255) // The last segment of the random IP
			r.firstIP[14] = randIPEndWith(255) // The last segment of the random IP

			targetIP := make([]byte, len(r.firstIP))
			copy(targetIP, r.firstIP)
			r.appendIP(targetIP) // Join the IP address pool

			for i := 13; i >= 0; i-- { // Random from the bottom three to the front
				tempIP = r.firstIP[i]              // save the previous value
				r.firstIP[i] += randIPEndWith(255) // Random 0~255, add to the current bit
				if r.firstIP[i] >= tempIP {        // If the value of the current bit is greater than or equal to the value of the previous bit, it means the random success, you can exit the loop
					break
				}
			}
		}
	}
}

func loadIPRanges() []*net.IPAddr {
	ranges := newIPRanges()
	if IPText != "" { // Get IP segment data from parameters
		IPs := strings.Split(IPText, ",") // comma separated as array and loop through
		for _, IP := range IPs {
			IP = strings.TrimSpace(IP) // Remove leading and trailing whitespace characters (spaces, tabs, newlines, etc.)
			if IP == "" {              // Skip empty ones (that is, the beginning, end, or multiple consecutive ,, cases)
				continue
			}
			ranges.parseCIDR(IP) // Parse the IP segment to get IP, IP range, subnet mask
			if isIPv4(IP) {      // Generate all IPv4/IPv6 addresses to test (single/random/all)
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
		for scanner.Scan() { // Loop through each line of the file
			line := strings.TrimSpace(scanner.Text()) // Remove leading and trailing whitespace characters (spaces, tabs, newlines, etc.)
			if line == "" {                           // skip empty lines
				continue
			}
			ranges.parseCIDR(line) // Parse the IP segment to get IP, IP range, subnet mask
			if isIPv4(line) {      // Generate all IPv4/IPv6 addresses to test (single/random/all)
				ranges.chooseIPv4()
			} else {
				ranges.chooseIPv6()
			}
		}
	}
	return ranges.ips
}
