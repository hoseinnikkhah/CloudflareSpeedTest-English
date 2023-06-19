package task

import (
	//"crypto/tls"
	//"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	Httping           bool
	HttpingStatusCode int
	HttpingCFColo     string
	HttpingCFColomap  *sync.Map
	OutRegexp         = regexp.MustCompile(`[A-Z]{3}`)
)

// pingReceived pingTotalTime
func (p *Ping) httping(ip *net.IPAddr) (int, time.Duration) {
	hc := http.Client{
		Timeout: time.Second * 2,
		Transport: &http.Transport{
			DialContext: getDialContext(ip),
			//TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Skip certificate verification
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // prevent redirection
		},
	}

	// Visit once to get HTTP status code and Cloudflare Colo
	{
		requ, err := http.NewRequest(http.MethodHead, URL, nil)
		if err != nil {
			return 0, 0
		}
		requ.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
		resp, err := hc.Do(requ)
		if err != nil {
			return 0, 0
		}
		defer resp.Body.Close()

		//fmt.Println("IP:", ip, "StatusCode:", resp.StatusCode, resp.Request.URL)
		// If the HTTP status code is not specified, or the specified status code is not compliant, only 200, 301, 302 are considered as HTTPing by default
		if HttpingStatusCode == 0 || HttpingStatusCode < 100 && HttpingStatusCode > 599 {
			if resp.StatusCode != 200 && resp.StatusCode != 301 && resp.StatusCode != 302 {
				return 0, 0
			}
		} else {
			if resp.StatusCode != HttpingStatusCode {
				return 0, 0
			}
		}

		io.Copy(io.Discard, resp.Body)

		// only match the three-character code of the airport if the region is specified
		if HttpingCFColo != "" {
			// Determine whether it is Cloudflare or AWS CloudFront through the header Server value and set cfRay to the complete content of the three-character code of the respective airport
			cfRay := func() string {
				if resp.Header.Get("Server") == "cloudflare" {
					return resp.Header.Get("CF-RAY") // Example cf-ray: 7bd32409eda7b020-SJC
				}
				return resp.Header.Get("x-amz-cf-pop") // Example X-Amz-Cf-Pop: SIN52-P1
			}()
			colo := p.getColo(cfRay)
			if colo == "" { // If no three-character code is matched or the specified region is not met, the IP test will be ended directly
				return 0, 0
			}
		}

	}

	// Loop tacho calculation delay
	success := 0
	var delay time.Duration
	for i := 0; i < PingTimes; i++ {
		requ, err := http.NewRequest(http.MethodHead, URL, nil)
		if err != nil {
			log.Fatal("Unexpected error, report:", err)
			return 0, 0
		}
		requ.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
		if i == PingTimes-1 {
			requ.Header.Set("Connection", "close")
		}
		startTime := time.Now()
		resp, err := hc.Do(requ)
		if err != nil {
			continue
		}
		success++
		io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		duration := time.Since(startTime)
		delay += duration

	}

	return success, delay

}

func MapColoMap() *sync.Map {
	if HttpingCFColo == "" {
		return nil
	}
	// Convert the three-character code of the region specified by the parameter to uppercase and format it
	colos := strings.Split(strings.ToUpper(HttpingCFColo), ",")
	colomap := &sync.Map{}
	for _, colo := range colos {
		colomap.Store(colo, colo)
	}
	return colomap
}

func (p *Ping) getColo(b string) string {
	if b == "" {
		return ""
	}
	// Regular match and return the three-character code of the airport
	out := OutRegexp.FindString(b)

	if HttpingCFColomap == nil {
		return out
	}
	// Match whether the three-character code of the airport is the specified area
	_, ok := HttpingCFColomap.Load(out)
	if ok {
		return out
	}

	return ""
}
