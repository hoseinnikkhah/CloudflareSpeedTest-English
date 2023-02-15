# XIU2/CloudflareSpeedTest - Script

Here are some scripts based on **XIU2/CloudflareSpeedTest** and **extended with more features**.
Can you tell me if you have any suggestions on the existing script functions, if you have some easy-to-use scripts for your own use, you can also pass [**Issues**](https://github.com/hoseinnikkhah/CloudflareSpeedTest-English/issues) Or send pull requests to me and add them here so that more people can use them (the author will be marked)

> Tip: Click the icon button with three horizontal bars in the upper left corner of ↖ to view the catalog~

****
## 📑 cfst_hosts.sh / cfst_hosts.bat (built in)

After running CloudflareST to get the fastest IP, the script replaces the old CDN IP in the Hosts file.

> **Instructions for use：https://github.com/XIU2/CloudflareSpeedTest/issues/42**

This thread needs translation, feel free to help.

<details>
<summary><code><strong>"Changelog"</strong></code></summary>

****

#### December 17, 2021, version v1.0.6
 - **1. optimization** [If no IP that satisfies the conditions can be found, the speed measurement will be repeated continuously] Function, there is no problem of re-testing when specifying the lower limit of download speed (default comment)   

#### December 17, 2021, version v1.0.3
 - **1. Add** If no IP that satisfies the conditions can be found, the speed measurement function will be cycled (default comment)  
 - **2. optimization** the code  

#### September 29, 2021, version v1.0.2
 - **1. repair** When the number of IPs in the speed test result is 0, the script does not exit the problem

#### April 29, 2021, version v1.0.1
 - **1. optimization** It is no longer necessary to add the -p 0 parameter to avoid the Enter key to exit (now the result can be displayed immediately, and there is no need to worry about the Enter key to exit the program)

#### January 28, 2021, version v1.0.0
 - **1. release** first version  

</details>

****

## 📑 cfst_3proxy.bat (built in)

The function of this script is to obtain the fastest IP after CloudflareST speed measurement and replace the Cloudflare CDN IP in the 3Proxy configuration file.
All Cloudflare CDN IPs can be redirected to the fastest IP to achieve once and for all acceleration of all websites using Cloudflare CDN (no need to add domain names to Hosts one by one).

> **Instructions for use：https://github.com/XIU2/CloudflareSpeedTest/discussions/71**

This thread needs translation, feel free to help.

<details>
<summary><code><strong>"Changelog"</strong></code></summary>

****

#### December 17, 2021, version v1.0.5
 - **1. optimization** [Continuous cycle speed measurement if no IP meeting the conditions can be found] function, there is no problem of re-measurement when the lower limit of download speed measurement is specified (default comment)   

#### December 17, 2021, version v1.0.4
 - **1. Add** If no IP that satisfies the conditions can be found, the speed measurement function will be cycled (default comment)  
 - **2. optimization** the code  

#### September 29, 2021, version v1.0.3
 - **1. repair** When the number of IPs in the speed test result is 0, the script does not exit the problem

#### April 29, 2021, version v1.0.2
 - **1. optimization** It is no longer necessary to add the -p 0 parameter to avoid the Enter key to exit (now the result can be displayed immediately, and there is no need to worry about the Enter key to exit the program)  

#### March 16, 2021, version v1.0.1
 - **1. optimization** Code and comment content  

#### March 13, 2021, version v1.0.0
 - **1. release** first version 

</details>

****

## 📑 cfst_ddns.sh / cfst_ddns.bat

If your domain name is hosted on Cloudflare, you can automatically update domain name resolution records through Cloudflare's official API!

> **Instructions for use：https://github.com/XIU2/CloudflareSpeedTest/issues/40**

This thread needs translation, feel free to help.

<details>
<summary><code><strong>"Changelog"</strong></code></summary>

****

#### 2December 17, 2021, version v1.0.4
 - **1. Add** If no IP that satisfies the conditions can be found, the speed measurement function will be cycled (default comment)  
 - **2. optimization** 代码  

#### September 29, 2021, version v1.0.3
 - **1. repair** When the number of IPs in the speed test result is 0, the script does not exit the problem  

#### April 29, 2021, version v1.0.2
 - **1. optimization** It is no longer necessary to add the -p 0 parameter to avoid the Enter key to exit (now the result can be displayed immediately, and there is no need to worry about the Enter key to exit the program)

#### January 27, 2021, version v1.0.1
 - **1. optimization** Configuration is read from a file

#### January 26, 2021, version v1.0.0
 - **1. release** first version

</details>

****

## Feature Suggestion/Question Feedback

If you have any problems, you can go to [**Issues**](https://github.com/hoseinnikkhah/CloudflareSpeedTest-English/issues) Check here to see if anyone else has asked (remember to check out here: [**Closed**](https://github.com/hoseinnikkhah/CloudflareSpeedTest-English/issues?q=is%3Aissue+is%3Aclosed)) 
If you don't find a similar question, please open a new one [**Issues**](https://github.com/hoseinnikkhah/CloudflareSpeedTest-English/issues/new) Come tell me!
