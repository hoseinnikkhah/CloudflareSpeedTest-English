package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	defaultOutput = "result.csv"
	maxDelay      = 9999 * time.Millisecond
	minDelay      = 0 * time.Millisecond
)

var (
	InputMaxDelay = maxDelay
	InputMinDelay = minDelay
	Output        = defaultOutput
	PrintNum      = 10
)

// Whether to print test results
func NoPrintResult() bool {
	return PrintNum == 0
}

// Whether to output to a file
func noOutput() bool {
	return Output == "" || Output == " "
}

type PingData struct {
	IP       *net.IPAddr
	Sended   int
	Received int
	Delay    time.Duration
}

type CloudflareIPData struct {
	*PingData
	recvRate      float32
	DownloadSpeed float64
}

func (cf *CloudflareIPData) getRecvRate() float32 {
	if cf.recvRate == 0 {
		pingLost := cf.Sended - cf.Received
		cf.recvRate = float32(pingLost) / float32(cf.Sended)
	}
	return cf.recvRate
}

func (cf *CloudflareIPData) toString() []string {
	result := make([]string, 6)
	result[0] = cf.IP.String()
	result[1] = strconv.Itoa(cf.Sended)
	result[2] = strconv.Itoa(cf.Received)
	result[3] = strconv.FormatFloat(float64(cf.getRecvRate()), 'f', 2, 32)
	result[4] = strconv.FormatFloat(cf.Delay.Seconds()*1000, 'f', 2, 32)
	result[5] = strconv.FormatFloat(cf.DownloadSpeed/1024/1024, 'f', 2, 32)
	return result
}

func ExportCsv(data []CloudflareIPData) {
	if noOutput() || len(data) == 0 {
		return
	}
	fp, err := os.Create(Output)
	if err != nil {
		log.Fatalf("file creation[%s]faild：%v", Output, err)
		return
	}
	defer fp.Close()
	w := csv.NewWriter(fp) //创建一个新的写入文件流
	_ = w.Write([]string{"IP address", "Sent", "Received", "Packet loss", "avg latency", "Speed (MB/s)"})
	_ = w.WriteAll(convertToString(data))
	w.Flush()
}

func convertToString(data []CloudflareIPData) [][]string {
	result := make([][]string, 0)
	for _, v := range data {
		result = append(result, v.toString())
	}
	return result
}

type PingDelaySet []CloudflareIPData

func (s PingDelaySet) FilterDelay() (data PingDelaySet) {
	if InputMaxDelay > maxDelay || InputMinDelay < minDelay {
		return s
	}
	for _, v := range s {
		if v.Delay > InputMaxDelay { // Average Latency Cap
			break
		}
		if v.Delay < InputMinDelay { // Average Latency Lower Limit
			continue
		}
		data = append(data, v) // Add to a new array when the condition is met lazily
	}
	return
}

func (s PingDelaySet) Len() int {
	return len(s)
}

func (s PingDelaySet) Less(i, j int) bool {
	iRate, jRate := s[i].getRecvRate(), s[j].getRecvRate()
	if iRate != jRate {
		return iRate < jRate
	}
	return s[i].Delay < s[j].Delay
}

func (s PingDelaySet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Sort by download speed
type DownloadSpeedSet []CloudflareIPData

func (s DownloadSpeedSet) Len() int {
	return len(s)
}

func (s DownloadSpeedSet) Less(i, j int) bool {
	return s[i].DownloadSpeed > s[j].DownloadSpeed
}

func (s DownloadSpeedSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s DownloadSpeedSet) Print() {
	if NoPrintResult() {
		return
	}
	if len(s) <= 0 { // Continue when the IP array length (number of IPs) is greater than 0
		fmt.Println("\n[Information] The IP quantity of the complete speed test result is 0, and the output result is skipped.")
		return
	}
	dateString := convertToString(s) // Convert to multidimensional array [][]String
	if len(dateString) < PrintNum {  // If IP array length (number of IPs) If it is less than the number of prints, the number of times is changed to the number of IPs
		PrintNum = len(dateString)
	}
	headFormat := "%-16s%-5s%-5s%-5s%-6s%-11s\n"
	dataFormat := "%-18s%-8s%-8s%-8s%-10s%-15s\n"
	for i := 0; i < PrintNum; i++ { // If the IP to be output contains IPv6, then you need to adjust the interval
		if len(dateString[i][0]) > 15 {
			headFormat = "%-40s%-5s%-5s%-5s%-6s%-11s\n"
			dataFormat = "%-42s%-8s%-8s%-8s%-10s%-15s\n"
			break
		}
	}
	fmt.Printf(headFormat, "IP address", "Sent", "Received", "Packet loss", "avg latency", "Speed (MB/s)")
	for i := 0; i < PrintNum; i++ {
		fmt.Printf(dataFormat, dateString[i][0], dateString[i][1], dateString[i][2], dateString[i][3], dateString[i][4], dateString[i][5])
	}
	if !noOutput() {
		fmt.Printf("\nComplete speed test results have been written %v Documents can be viewed using Notepad/Spreads software.\n", Output)
	}
}
