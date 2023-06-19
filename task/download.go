package task

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/XIU2/CloudflareSpeedTest/utils"

	"github.com/VividCortex/ewma"
)

const (
	bufferSize                     = 1024
	defaultURL                     = "https://cloudflare.cdn.openbsd.org/pub/OpenBSD/7.1/alpha/install71.iso"
	defaultTimeout                 = 10 * time.Second
	defaultDisableDownload         = false
	defaultTestNum                 = 10
	defaultMinSpeed        float64 = 0.0
)

var (
	URL     = defaultURL
	Timeout = defaultTimeout
	Disable = defaultDisableDownload

	TestCount = defaultTestNum
	MinSpeed  = defaultMinSpeed
)

func checkDownloadDefault() {
	if URL == "" {
		URL = defaultURL
	}
	if Timeout <= 0 {
		Timeout = defaultTimeout
	}
	if TestCount <= 0 {
		TestCount = defaultTestNum
	}
	if MinSpeed <= 0.0 {
		MinSpeed = defaultMinSpeed
	}
}

func TestDownloadSpeed(ipSet utils.PingDelaySet) (speedSet utils.DownloadSpeedSet) {
	checkDownloadDefault()
	if Disable {
		return utils.DownloadSpeedSet(ipSet)
	}
	if len(ipSet) <= 0 { // Only when the IP array length (number of IPs) is greater than 0 will the download speed test continue
		fmt.Println("\n[INFO] Delay speed test result IP number is 0, skip download speed test.")
		return
	}
	testNum := TestCount
	if len(ipSet) < TestCount || MinSpeed > 0 { // If the length of the IP array (number of IPs) is less than the number of download speed measurements (-dn), the number of times is corrected to the number of IPs
		testNum = len(ipSet)
	}
	if testNum < TestCount {
		TestCount = testNum
	}

	fmt.Printf("Start download speed measurement (download speed lower limit:%.2f MB/s，Number of download speed tests：%d，Download speed test queue：%d）：\n", MinSpeed, TestCount, testNum)
	// Control the download speed measurement progress bar to be the same length as the delayed speed measurement progress bar (obsessive-compulsive disorder)
	bar_a := len(strconv.Itoa(len(ipSet)))
	bar_b := "     "
	for i := 0; i < bar_a; i++ {
		bar_b += " "
	}
	bar := utils.NewBar(TestCount, bar_b, "")
	for i := 0; i < testNum; i++ {
		speed := downloadHandler(ipSet[i].IP)
		ipSet[i].DownloadSpeed = speed
		// After each IP download speed test, filter the results by the [download speed lower limit] condition
		if speed >= MinSpeed*1024*1024 {
			bar.Grow(1, "")
			speedSet = append(speedSet, ipSet[i]) // When the download speed is higher than the lower limit, add it to a new array
			if len(speedSet) == TestCount {       // When there are enough IPs that meet the conditions (download speed measurement number -dn), it will jump out of the loop
				break
			}
		}
	}
	bar.Done()
	if len(speedSet) == 0 { // There is no data that meets the speed limit, return all test data
		speedSet = utils.DownloadSpeedSet(ipSet)
	}
	// Sort by speed
	sort.Sort(speedSet)
	return
}

func getDialContext(ip *net.IPAddr) func(ctx context.Context, network, address string) (net.Conn, error) {
	var fakeSourceAddr string
	if isIPv4(ip.String()) {
		fakeSourceAddr = fmt.Sprintf("%s:%d", ip.String(), TCPPort)
	} else {
		fakeSourceAddr = fmt.Sprintf("[%s]:%d", ip.String(), TCPPort)
	}
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, network, fakeSourceAddr)
	}
}

// return download Speed
func downloadHandler(ip *net.IPAddr) float64 {
	client := &http.Client{
		Transport: &http.Transport{DialContext: getDialContext(ip)},
		Timeout:   Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 10 { // Limit up to 10 redirects
				return http.ErrUseLastResponse
			}
			if req.Header.Get("Referer") == defaultURL { // When using the default download speed test URL, the redirect does not carry the Referer
				req.Header.Del("Referer")
			}
			return nil
		},
	}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return 0.0
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")

	response, err := client.Do(req)
	if err != nil {
		return 0.0
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return 0.0
	}
	timeStart := time.Now()           // start time (current)
	timeEnd := timeStart.Add(Timeout) // Add the end time obtained by downloading the speed test time

	contentLength := response.ContentLength // File size
	buffer := make([]byte, bufferSize)

	var (
		contentRead     int64 = 0
		timeSlice             = Timeout / 100
		timeCounter           = 1
		lastContentRead int64 = 0
	)

	var nextTime = timeStart.Add(timeSlice * time.Duration(timeCounter))
	e := ewma.NewMovingAverage()

	// Loop calculation, if the file is downloaded (the two are equal), exit the loop (terminate the speed test)
	for contentLength != contentRead {
		currentTime := time.Now()
		if currentTime.After(nextTime) {
			timeCounter++
			nextTime = timeStart.Add(timeSlice * time.Duration(timeCounter))
			e.Add(float64(contentRead - lastContentRead))
			lastContentRead = contentRead
		}
		// If the download speed test time is exceeded, exit the loop (terminate the speed test)
		if currentTime.After(timeEnd) {
			break
		}
		bufferRead, err := response.Body.Read(buffer)
		if err != nil {
			if err != io.EOF { // If an error (such as Timeout) is encountered during the file download process, and it is not because the file download is complete, exit the loop (terminate speed measurement)
				break
			} else if contentLength == -1 { // When the file download is complete and the file size is unknown, exit the loop (terminate the speed test)，For example：https://speed.cloudflare.com/__down?bytes=200000000 In this way, if the download is completed within 10 seconds, the speed test result will be significantly lower or even displayed as 0.00 (when the download speed is too fast)
				break
			}
			// Get last time slice
			last_time_slice := timeStart.Add(timeSlice * time.Duration(timeCounter-1))
			// download data volume / (use current time - last time slice/ time slice)
			e.Add(float64(contentRead-lastContentRead) / (float64(currentTime.Sub(last_time_slice)) / float64(timeSlice)))
		}
		contentRead += int64(bufferRead)
	}
	return e.Value() / (Timeout.Seconds() / 120)
}
