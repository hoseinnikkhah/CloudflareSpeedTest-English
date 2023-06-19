package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/XIU2/CloudflareSpeedTest/task"
	"github.com/XIU2/CloudflareSpeedTest/utils"
)

var (
	version, versionNew string
)

func init() {
	var printVersion bool
	var help = `
CloudflareSpeedTest ` + version + `
Test the latency and speed of all IPs of Cloudflare CDN, get the fastest IP (IPv4+IPv6)!
https://github.com/XIU2/CloudflareSpeedTest

parameter:
    -n 200
		Delay speed measurement thread; the more delays, the faster the speed measurement, and devices with weak performance (such as routers) should not be too high; (default 200, maximum 1000)
    -t 4
		Delay speed test times; single IP delay speed test times; (default 4 times)
    -dn 10
		Number of download speed tests; after delay and sorting, the number of download speed tests from the lowest delay; (default 10)
    -dt 10
		Download speed test time; the maximum time for a single IP download speed test, not too short; (default 10 seconds)
    -tp 443
		Specify the speed test port; the port used for delay speed test/download speed test; (default port 443)
    -url https://cf.xiu2.xyz/url
		Specify the speed test address; the address used for delayed speed test (HTTPing)/download speed test. The default address does not guarantee availability, and it is recommended to build it yourself;

    -httping
		Switch the speed measurement mode; change the delayed speed measurement mode to the HTTP protocol, and the test address used is the [-url] parameter; (default TCPing)
    -httping-code 200
		Valid status code; only one valid HTTP status code returned by the web page during HTTPing delay speed test; (default 200 301 302)
    -cfcolo HKG,KHH,NRT,LAX,SEA,SJC,FRA,MAD
		Match the specified region; the region name is the three-character code of the local airport, separated by English commas, only available in HTTPing mode; (default all regions)

    -tl 200
		The upper limit of the average delay; only output IPs lower than the specified average delay, and the upper and lower limit conditions can be used together; (default 9999 ms)
    -tll 40
		Average delay lower limit; only output IPs higher than the specified average delay; (default 0 ms)
    -tlr 0.2
		The upper limit of the packet loss probability; only output the IP that is lower than/equal to the specified packet loss rate, the range is 0.00~1.00, 0 filters out any IP with packet loss; (default 1.00)
    -sl 5
		Download speed lower limit; only output IPs higher than the specified download speed, and the speed measurement will stop when the specified number [-dn] is reached; (default 0.00 MB/s)

    -p 10
		Display the number of results; directly display the specified number of results after the speed measurement, and exit without displaying the results when it is 0; (default 10)
    -f ip.txt
		IP segment data file; if the path contains spaces, please add quotation marks; support other CDN IP segments; (default ip.txt)
    -ip 1.1.1.1,2.2.2.2/24,2606:4700::/32
		Specify the IP segment data; directly specify the IP segment data to be tested by parameters, separated by English commas; (default empty)
    -o result.csv
		Write the result file; if the path contains spaces, please add quotation marks; if the value is empty, do not write to the file [-o ""]; (default result.csv)

    -dd
		Disable download speed measurement; after disabling, the speed test results will be sorted by delay (default sorted by download speed); (default enabled)
    -allip
		Speed test all IPs; test the speed of each IP in the IP segment (only supports IPv4); (by default, each /24 segment will randomly test one IP)

    -v
		Print program version + check for version updates
    -h
	print help instructions
`
	var minDelay, maxDelay, downloadTime int
	var maxLossRate float64
	flag.IntVar(&task.Routines, "n", 200, "delay thread")
	flag.IntVar(&task.PingTimes, "t", 4, "Delay speed test times")
	flag.IntVar(&task.TestCount, "dn", 10, "Number of download speed tests")
	flag.IntVar(&downloadTime, "dt", 10, "Download speed test time")
	flag.IntVar(&task.TCPPort, "tp", 443, "Designated speed test port")
	flag.StringVar(&task.URL, "url", "https://cf.xiu2.xyz/url", "Specify speed test address")

	flag.BoolVar(&task.Httping, "httping", false, "Switch speed measurement mode")
	flag.IntVar(&task.HttpingStatusCode, "httping-code", 0, "valid status code")
	flag.StringVar(&task.HttpingCFColo, "cfcolo", "", "Match the specified region")

	flag.IntVar(&maxDelay, "tl", 9999, "Average Latency Cap")
	flag.IntVar(&minDelay, "tll", 0, "Average Latency Lower Limit")
	flag.Float64Var(&maxLossRate, "tlr", 1, "Maximum chance of packet loss")
	flag.Float64Var(&task.MinSpeed, "sl", 0, "download speed limit")

	flag.IntVar(&utils.PrintNum, "p", 10, "Show number of results")
	flag.StringVar(&task.IPFile, "f", "ip.txt", "IP segment data file")
	flag.StringVar(&task.IPText, "ip", "", "Specify IP segment data")
	flag.StringVar(&utils.Output, "o", "result.csv", "output result file")

	flag.BoolVar(&task.Disable, "dd", false, "Disable download speed test")
	flag.BoolVar(&task.TestAll, "allip", false, "Speed test all IP")

	flag.BoolVar(&printVersion, "v", false, "print program version")
	flag.Usage = func() { fmt.Print(help) }
	flag.Parse()

	if task.MinSpeed > 0 && time.Duration(maxDelay)*time.Millisecond == utils.InputMaxDelay {
		fmt.Println("[Tips] In use [-sl] parameters, it is recommended to match [-tl] parameters to avoid insufficient [-dn] Quantity and speed measurement...")
	}
	utils.InputMaxDelay = time.Duration(maxDelay) * time.Millisecond
	utils.InputMinDelay = time.Duration(minDelay) * time.Millisecond
	utils.InputMaxLossRate = float32(maxLossRate)
	task.Timeout = time.Duration(downloadTime) * time.Second
	task.HttpingCFColomap = task.MapColoMap()

	if printVersion {
		println(version)
		fmt.Println("Checking for version updates...")
		checkUpdate()
		if versionNew != "" {
			fmt.Printf("*** new version found [%s]！please go to [https://github.com/hoseinnikkhah/CloudflareSpeedTest-English] renew! ***", versionNew)
		} else {
			fmt.Println("Currently the latest version [" + version + "]！")
		}
		os.Exit(0)
	}
}

func main() {
	task.InitRandSeed() // set random number seed

	fmt.Printf("# hoseinnikkhah/CloudflareSpeedTest-English %s \n\n", version)

	// Start Latency Test + Filter Latency/Packet Loss
	pingData := task.NewPing().Run().FilterDelay().FilterLossRate()
	// Start download speed test
	speedData := task.TestDownloadSpeed(pingData)
	utils.ExportCsv(speedData) // output file
	speedData.Print()          // print result

	if versionNew != "" {
		fmt.Printf("\n*** new version found [%s]！please go to [https://github.com/hoseinnikkhah/CloudflareSpeedTest-English] renew! ***\n", versionNew)
	}
	endPrint()
}

func endPrint() {
	if utils.NoPrintResult() {
		return
	}
	if runtime.GOOS == "windows" { // If it is a Windows system, you need to press the Enter key or Ctrl+C to exit (avoid double-clicking to run, and close it directly after the speed measurement is completed)
		fmt.Printf("Press Enter or Ctrl+C to exit.")
		fmt.Scanln()
	}
}

// Check for updates
func checkUpdate() {
	timeout := 10 * time.Second
	client := http.Client{Timeout: timeout}
	res, err := client.Get("https://api.xiu2.xyz/ver/cloudflarespeedtest.txt")
	if err != nil {
		return
	}
	// read resource data body: []byte
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	// 关闭资源流
	defer res.Body.Close()
	if string(body) != version {
		versionNew = string(body)
	}
}
