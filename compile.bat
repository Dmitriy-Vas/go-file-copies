set GOARCH=386
go build -o WindowsCopies.exe main.go
set GOOS=darwin
go build -o DarwinCopies main.go
set GOOS=linux
go build -o LinuxCopies main.go
