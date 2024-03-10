@echo off
setlocal
cd /d %~dp0

:: Set Make version
set MAKE_VER=4.4.1

:: Get the current timestamp
:: See: https://stackoverflow.com/a/19163883
for /f %%a in ('powershell -Command "Get-Date -format yyyyMMdd_HHmmss"') do set TSTAMP=%%a

:: Download the Make archive from Chocolatey
:: See: https://superuser.com/a/373068
set WORKDIR=%temp%
set MAKE_ZIP=%WORKDIR%\make_%MAKE_VER%_%TSTAMP%.zip
set MAKE_FOLDER=%WORKDIR%\make_%MAKE_VER%_%TSTAMP%
certutil.exe -urlcache -split -f "https://community.chocolatey.org/api/v2/package/make/%MAKE_VER%" %MAKE_ZIP%

:: Drop a vbs file to extract the zip file
Call :UnZipFile "%MAKE_FOLDER%" "%MAKE_ZIP%"

:: Drop `make.exe`
copy /y "%MAKE_FOLDER%\tools\install\bin\make.exe" "%cd%\make.exe"

:: Cleanup
del %MAKE_ZIP%
rd /s /q %MAKE_FOLDER%
exit /b

:: ===============================================

:: Utils
:: See: https://stackoverflow.com/a/21709923
:UnZipFile <ExtractTo> <newzipfile>
set vbs="%temp%\%TSTAMP%_.vbs"
if exist %vbs% del /f /q %vbs%
>%vbs%  echo Set fso = CreateObject("Scripting.FileSystemObject")
>>%vbs% echo If NOT fso.FolderExists(%1) Then
>>%vbs% echo fso.CreateFolder(%1)
>>%vbs% echo End If
>>%vbs% echo set objShell = CreateObject("Shell.Application")
>>%vbs% echo set FilesInZip=objShell.NameSpace(%2).items
>>%vbs% echo objShell.NameSpace(%1).CopyHere(FilesInZip)
>>%vbs% echo Set fso = Nothing
>>%vbs% echo Set objShell = Nothing
cscript //nologo %vbs%
if exist %vbs% del /f /q %vbs%
if exist %vbs% del /f /q %vbs%
