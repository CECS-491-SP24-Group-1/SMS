@echo off

::See: https://stackoverflow.com/questions/41862786/how-can-i-download-a-file-from-the-internet-via-command-prompt

::Set Make version
set MAKEVER=4.4.1

:: Get the current timestamp
:: See: https://stackoverflow.com/a/19163883
for /f %%a in ('powershell -Command "Get-Date -format yyyyMMdd_HHmmss"') do set TSTAMP=%%a

:: Download the Make archive from Chocolatey
:: See: https://superuser.com/a/373068
set WORKDIR=%temp%
set MAKE_ZIP=%WORKDIR%\make_%MAKE_VER%_%TSTAMP%.zip
set MAKE_FOLDER=%WORKDIR%\make_%MAKE_VER%_%TSTAMP%
powershell -c "Invoke-WebRequest -Uri "https://community.chocolatey.org/api/v2/package/make/%MAKEVER%" -OutFile '%MAKE_ZIP%'"

::Unzip the downloaded archive
powershell -command "Expand-Archive %MAKE_ZIP% %MAKE_FOLDER%"

:: Drop `make.exe`
copy /y "%MAKE_FOLDER%\tools\install\bin\make.exe" "%cd%\make.exe"

:: Cleanup
del %MAKE_ZIP%
rd /s /q %MAKE_FOLDER%
exit /b