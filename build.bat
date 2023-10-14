SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=386
go build -o emailcat


SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=386
go build -o emailcat.exe
