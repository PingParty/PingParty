export GOOS=linux
export GOARCH=amd64
go build -o main lambda.go
build-lambda-zip -o main.zip main