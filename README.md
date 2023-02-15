Tool is not translated yet, I am working on it and recompling it may take a little longer than expected since I am  a bit busy.
I wont update this md file until finished.

# Cloudflare Speed Test

![GitHub repo file count](https://img.shields.io/badge/Status-Under%20Development-%23e86c2e)

![GitHub repo file count](https://img.shields.io/badge/main.go-translated%2Fnot%20compiled-yellow)
![GitHub repo file count](https://img.shields.io/badge/LICENSE-rework%20needed-red)
![GitHub repo file count](https://img.shields.io/badge/csv.go-not%20translated%20-red)
![GitHub repo file count](https://img.shields.io/badge/download.go-not%20translated%20-red)
![GitHub repo file count](https://img.shields.io/badge/ip.go-not%20translated%20-red)
![GitHub repo file count](https://img.shields.io/badge/httping.go-not%20translated%20-red)
![GitHub repo file count](https://img.shields.io/badge/tcping.go-not%20translated%20-red)
![GitHub repo file count](https://img.shields.io/badge/cfst__3proxy.bat-not%20translated%20-red)
![GitHub repo file count](https://img.shields.io/badge/cfst__ddns.bat.bat-not%20translated%20-red)
![GitHub repo file count](https://img.shields.io/badge/cfst__ddns.sh-not%20translated%20-red)
![GitHub repo file count](https://img.shields.io/badge/cfst__hosts.bat-not%20translated%20-red)
![GitHub repo file count](https://img.shields.io/badge/cfst_hosts.sh-not%20translated%20-red)
![GitHub repo file count](https://img.shields.io/badge/cfst_hosts_mac.sh-not%20translated%20-red)

![GitHub User's stars](https://img.shields.io/github/stars/hoseinnikkhah?style=social)
![GitHub forks](https://img.shields.io/github/forks/hoseinnikkhah/better-cloudflare-ip-english?style=social)

> * This is forked version, code is still under review and is not stable in English yet.
I try to finish the translation as soon as I can but still it needs considerable amount of time to review whole code.

Many foreign websites are using Cloudflare CDN, but the IPs assigned to visitors in countries dealing with censorship are not friendly (high latency, high packet loss, and slow speed).
Although Cloudflare has disclosed all [IP segments](https://www.cloudflare.com/ips/), but if you want to find the one that suits you among so many IPs, you may get exhausted, in order to make it faster and avoid wasting time this tool was developed.

**"Self-selected preferred IP" tests Cloudflare CDN latency and speed, get the fastest IP (IPv4+IPv6)! If it’s useful to you, please click ⭐ to encourage us :)**

> **This project also supports other CDN / website IP** Latency speed measurement (such as Cloudflare, Gcore)

> **Please do not rely too much on Cloudflare and do not make it your only solution, it's true that too many serveices are using cloudflare but also governments can always block IPs in certain times** 

# Quick Use
### download and run
1. Download the compiled executable file [Lanzouv](https://pan.lanzouv.com/b0742hkxe) / [Github](https://github.com/XIU2/CloudflareSpeedTest/releases) and unzip it.
2. Double-click to run the `CloudflareST.exe` file (Windows system), wait for the speed test to complete...

<details>
<summary><code><strong>"Click to view the usage example under Linux system"</strong></code></summary>

****

The following commands are examples only, please go to the version number you need and file name. Check:[**Releases**](https://github.com/XIU2/CloudflareSpeedTest/releases)

``` yaml
# If it is your first time using the tool, it is recommended to create a new folder (skip this step for subsequent updates)
mkdir CloudflareST

# Go to the folder (for subsequent updates, just repeat the download and decompression commands below from here)
cd CloudflareST

# Download the CloudflareST compressed package (replace [version number] and [file name] in the URL according to your needs)
# wget -N https://github.com/XIU2/CloudflareSpeedTest/releases/download/v2.2.2/CloudflareST_linux_amd64.tar.gz

# If you are downloading from a domestic server, please use the following mirrors to speed up:
# wget -N https://download.fastgit.org/XIU2/CloudflareSpeedTest/releases/download/v2.2.2/CloudflareST_linux_amd64.tar.gz
# wget -N https://ghproxy.com/https://github.com/XIU2/CloudflareSpeedTest/releases/download/v2.2.2/CloudflareST_linux_amd64.tar.gz

# Unzip (you don’t need to delete the old file, it will be overwritten directly, and you can replace the file name according to your needs)
tar -zxf CloudflareST_linux_amd64.tar.gz
# Give execute permission
chmod +x CloudflareST
# run (without arguments)

./CloudflareST
# run (example with parameters)
./CloudflareST -dd -tll 90
```

If the average **average latency is very low** (such as 0.xx), it means that CloudflareST **passed the proxy** during the speed measurement. Please close the proxy software before measuring the speed.
If running on a **router**, it is recommended to turn off the proxy inside the router (or exclude it), otherwise the speed test results may be **inaccurate/unusable**.

</details>

****

_A Simple Tutorial for Running CloudflareST Speed Test Standalone on **Mobile**：**[Android](https://github.com/XIU2/CloudflareSpeedTest/discussions/61)、[Android APP](https://github.com/xianshenglu/cloudflare-ip-tester-app)、[IOS](https://github.com/XIU2/CloudflareSpeedTest/issues/151)**_

### Example result

After the speed test is completed, the **fastest 10 IPs** will be displayed by default, for example:

``` bash
IP address        Sent    Received  Packet loss     avg latency   speed (MB/s)
104.27.200.69     4       4         0.00            146.23        28.64
172.67.60.78      4       4         0.00            139.82        15.02
104.25.140.153    4       4         0.00            146.49        14.90
104.27.192.65     4       4         0.00            140.28        14.07
172.67.62.214     4       4         0.00            139.29        12.71
104.27.207.5      4       4         0.00            145.92        11.95
172.67.54.193     4       4         0.00            146.71        11.55
104.22.66.8       4       4         0.00            147.42        11.11
104.27.197.63     4       4         0.00            131.29        10.26
172.67.58.91      4       4         0.00            140.19        9.14
...
# If the average latency is very low (such as 0.xx), it means that CloudflareST uses a proxy when measuring the speed. Please close the proxy software before measuring the speed.
# If running on a router, please turn off the proxy inside the router first (or exclude it), otherwise the speed test results may be inaccurate/unusable.
# Because each speed test uses a random IP in each IP segment, the results of each speed test may not be the same, which is normal!
# Notice! I found that the delay of the first speed test after the computer is turned on will be obviously high (the same is true for manual TCPing), and subsequent speed tests are normal
# Therefore, it is recommended that you test a few IPs at random before the first official speed test after booting (no need to wait for the delay to complete the speed test, as long as the progress bar moves, you can directly close it)
# The general steps of the whole process of the software under the default parameters:
# 1. Delay speed measurement (default TCPing mode, HTTPing mode needs to manually add parameters)
# 2. Delay sorting (the delay is sorted from low to high, different packet loss rates will be sorted separately and independently, so there may be some IPs with low delay but packet loss that are sorted to the back)
# 3. Download speed measurement (download speed measurement from the IP with the lowest delay, and the default will stop when 10 are measured)
# 4. Speed sorting (speed sorting from high to low)
# 5. Output results (you can rely on parameters to control whether to output to the command line (-p 0)/file (-o ""))
```








