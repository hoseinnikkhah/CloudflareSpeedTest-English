

# Cloudflare Speed Test

![GitHub repo file count](https://img.shields.io/badge/Status-%20Developed-brightgreen)
![GitHub repo file count](https://img.shields.io/badge/LICENSE-GPL%203-pink)

![GitHub repo file count](https://img.shields.io/badge/main.go-%20Translated-brightgreen)
![GitHub repo file count](https://img.shields.io/badge/csv.go-%20Translated-brightgreen)
![GitHub repo file count](https://img.shields.io/badge/download.go-%20Translated-brightgreen)
![GitHub repo file count](https://img.shields.io/badge/ip.go-%20Translated-brightgreen)
![GitHub repo file count](https://img.shields.io/badge/httping.go-%20Translated-brightgreen)
![GitHub repo file count](https://img.shields.io/badge/tcping.go-%20Translated-brightgreen)
![GitHub repo file count](https://img.shields.io/badge/cfst__3proxy.bat-Translated-brightgreen)
![GitHub repo file count](https://img.shields.io/badge/cfst__ddns.bat-Translated-brightgreen)
![GitHub repo file count](https://img.shields.io/badge/cfst__ddns.sh-Translated-brightgreen)
![GitHub repo file count](https://img.shields.io/badge/cfst__hosts.bat-Translated-brightgreen)
![GitHub repo file count](https://img.shields.io/badge/cfst_hosts.sh-Translated-brightgreen)
![GitHub repo file count](https://img.shields.io/badge/cfst_hosts_mac.sh-Translated-brightgreen)

![GitHub User's stars](https://img.shields.io/github/stars/hoseinnikkhah?style=social)
![GitHub forks](https://img.shields.io/github/forks/hoseinnikkhah/CloudflareSpeedTest-English?style=social)

> * This is forked version, code has been reviewed completely and it is stable in English. Please report any bug if you faced one


Many foreign websites are using Cloudflare CDN, but the IPs assigned to visitors in countries dealing with censorship are not friendly (high latency, high packet loss, and slow speed).
Although Cloudflare has disclosed all [IP segments](https://www.cloudflare.com/ips/), but if you want to find the one that suits you among so many IPs, you may get exhausted, in order to make it faster and avoid wasting time this tool was developed.

**"Self-selected preferred IP" tests Cloudflare CDN latency and speed, get the fastest IP (IPv4+IPv6)! If it’s useful to you, please click ⭐ to encourage us :)**

> **This project also supports other CDN / website IP** Latency speed measurement (such as Cloudflare, Gcore)

> **Please do not rely too much on Cloudflare and do not make it your only solution, it's true that too many serveices are using cloudflare but also governments can always block IPs in certain times** 

# Quick Use
### download and run
1. Download the compiled executable file [Lanzouv](https://pan.lanzouv.com/b0742hkxe) / [Github](https://github.com/hoseinnikkhah/CloudflareSpeedTest-English/releases) and unzip it.
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

# How to use

It's pretty easy to use this software, Just download right version and run the app, keep in mind app won't run if there's no *ip.txt* file so you need to create it in first place, and also keep in mind this app only can check cloudflare IP ranges but it has an option for other CDNs which I will tell in next section.
Mac/Linux users need to give permission in some cases, right click on the downloaded file and excute it as an application, otherwsie app won't work.

# Konwn Issues

These are issues I found during testing the app, basically I won't call them issues since there are very basic solutios provided by original author.

> App does not open

It's not an issue with app, you need to creat a text file and rename it to *ip* and pate all cloudflare IP ranges in the file.
you can also copy these IPs and use them as you wish:
```
173.245.48.0/20
103.21.244.0/22
103.22.200.0/22
103.31.4.0/22
141.101.64.0/18
108.162.192.0/18
190.93.240.0/20
188.114.96.0/20
197.234.240.0/22
198.41.128.0/17
162.158.0.0/15
104.16.0.0/13
104.24.0.0/14
172.64.0.0/13
131.0.72.0/22
```

> Speed download is 0.00 Mb

This is an issue cuased by cloudflare itself, main author used his own server for speed test and due to high use of single link, cloudflare denies to send or recive data from that specefic link, you can either change this link directly from source code and recompile the wole code for yourself or you can excute the link with a single command, point is that you can't just use any link and it must be a cloudflare link or if you use other CDNs, it must be a link from a server that use that specefic CDN. You can also use your own server for this but keep in mind that target file must be over 200Mb in size.

### How to solve the issue?

Depending on your OS there are different ways to excute a specfic link for speed check

To use a different download speed address, simply add a parameter when running CloudflareST, for example:
```
# Windows
CloudflareST.exe -url https://speed.cloudflare.com/__down?bytes=200000000

# Linux/Mac
./CloudflareST -url https://speed.cloudflare.com/__down?bytes=200000000
```
Here are some ready to use links:
```
-url https://cloudflaremirrors.com/archlinux/community/os/x86_64/endless-sky-high-dpi-0.9.16-1-any.pkg.tar.zst
-url https://cloudflaremirrors.com/archlinux/images/latest/Arch-Linux-x86_64-basic.qcow2
-url https://cloudflaremirrors.com/archlinux/iso/latest/archlinux-x86_64.iso
-url https://download.parallels.com/desktop/v15/15.1.5-47309/ParallelsDesktop-15.1.5-47309.dmg
-url https://download.parallels.com/desktop/v17/17.1.1-51537/ParallelsDesktop-17.1.1-51537.dmg
-url https://cloudflare.cdn.openbsd.org/pub/OpenBSD/7.1/src.tar.gz
-url https://cloudflare.cdn.openbsd.org/pub/OpenBSD/7.0/i386/base70.tgz
-url https://cloudflare.cdn.openbsd.org/pub/OpenBSD/7.1/alpha/install71.iso
-url https://speedtest.galgamer.eu.org/200m.png
```

# Licence

GPL v3.0




