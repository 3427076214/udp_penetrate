set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0 
go build -o udp_server

if %ERRORLEVEL% gtr 0 (
    pause
)