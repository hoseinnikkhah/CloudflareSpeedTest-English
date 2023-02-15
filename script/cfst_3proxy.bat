:: --------------------------------------------------------------
::	project: CloudflareSpeedTest-English auto-updates 3Proxy
::	Version: 1.0.5
::	author: hoseinnikkhah (forked from XIU2)
::	project: https://github.com/hoseinnikkhah/CloudflareSpeedTest-English
:: --------------------------------------------------------------
@echo off
Setlocal Enabledelayedexpansion

::Determine whether administrator privileges have been obtained

>nul 2>&1 "%SYSTEMROOT%\system32\cacls.exe" "%SYSTEMROOT%\system32\config\system" 

if '%errorlevel%' NEQ '0' (  
    goto UACPrompt  
) else ( goto gotAdmin )  

::Write a vbs script to run this script (bat) as an administrator

:UACPrompt  
    echo Set UAC = CreateObject^("Shell.Application"^) > "%temp%\getadmin.vbs" 
    echo UAC.ShellExecute "%~s0", "", "", "runas", 1 >> "%temp%\getadmin.vbs" 
    "%temp%\getadmin.vbs" 
    exit /B  

::Delete the temporary vbs script if it exists
  
:gotAdmin  
    if exist "%temp%\getadmin.vbs" ( del "%temp%\getadmin.vbs" )  
    pushd "%CD%" 
    CD /D "%~dp0" 


::The above is to judge whether to obtain administrator privileges, if not, to obtain, the following is the main code of this script


::If the nowip_3proxy.txt file does not exist, this is the first time the script is being run
if not exist "nowip_3proxy.txt" (
    echo The function of this script is to obtain the fastest IP after CloudflareST speed measurement and replace the Cloudflare CDN IP in the 3Proxy configuration file.
    echo All Cloudflare CDN IPs can be redirected to the fastest IP to achieve once and for all acceleration of all websites using Cloudflare CDN [no need to add domain names to Hosts one by one]
    echo Please read before using：https://github.com/XIU2/CloudflareSpeedTest/discussions/71 [Need help with translation]
    echo.
    set /p nowip="Enter the Cloudflare CDN IP currently being used by 3Proxy and press Enter (this step is no longer needed later):"
    echo !nowip!>nowip_3proxy.txt
    echo.
)  

::Get the currently used Cloudflare CDN IP from the nowip_3proxy.txt file
set /p nowip=<nowip_3proxy.txt
echo start speed test...


:: This RESET is prepared for those who need the function of "repeating the speed test if no IP that satisfies the conditions is found"
:: If you need this function, just change the following 3 goto :STOP to goto :RESET
:RESET


:: Here you can add and modify the running parameters of CloudflareST yourself，echo.| The function is to automatically enter and exit the program (no need to add -p 0 parameter anymore)
echo.|CloudflareST.exe -o "result_3proxy.txt"


:: Determine whether the result file exists, if not, the result is 0
if not exist result_3proxy.txt (
    echo.
    echo The number of IPs in the CloudflareST speed measurement result is 0, skip the following steps...
    goto :STOP
)

:: Get the fastest IP on the first line
for /f "tokens=1 delims=," %%i in (result_3proxy.txt) do (
    set /a n+=1 
    If !n!==2 (
        set bestip=%%i
        goto :END
    )
)
:END

:: Determine whether the fastest IP just obtained is empty, and whether it is the same as the old IP
if "%bestip%"=="" (
    echo.
    echo The number of IPs in the CloudflareST speed measurement result is 0, skip the following steps...
    goto :STOP
)
if "%bestip%"=="%nowip%" (
    echo.
    echo The number of IPs in the CloudflareST speed measurement result is 0, skip the following steps...
    goto :STOP
)


:: The following piece of code is the code that is only needed for "if no IP that satisfies the conditions is found, the speed measurement will continue continuously"
:: Considering that when the lower limit of the download speed is specified, but an IP that meets all the conditions is not found, CloudflareST will output all IP results
:: Therefore, when you specify the -sl parameter, you need to remove the :: colon comment at the beginning of the following code to judge the number of file lines (for example, the number of download speed tests: 10, then the following value is set to 11)
::set /a v=0
::for /f %%a in ('type result_3proxy.txt') do set /a v+=1
::if %v% GTR 11 (
::    echo.
::    echo CloudflareST speed test results did not find an IP that fully meets the conditions, re-test the speed...
::    goto :RESET
::)


echo %bestip%>nowip_3proxy.txt
echo.
echo The old IP was %nowip%
echo The new IP is %bestip%



:: Please change D:\Program Files\3Proxy in quotation marks to the directory where your 3Proxy program is located
CD /d "D:\Program Files\3Proxy"
:: Please make sure that you have tested that 3Proxy can run normally and use it before running this script!



echo.
echo Start backup of 3proxy.cfg file (3proxy.cfg_backup)...
copy 3proxy.cfg 3proxy.cfg_backup
echo.
echo start to replace...
(
    for /f "tokens=*" %%i in (3proxy.cfg_backup) do (
        set s=%%i
        set s=!s:%nowip%=%bestip%!
        echo !s!
        )
)>3proxy.cfg

net stop 3proxy
net start 3proxy

echo Finish...
echo.
:STOP
pause 
