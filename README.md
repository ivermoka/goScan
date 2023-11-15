# goScan

goScan is a program used for checking go version and docker FROM line. 

run command in directory with files and goScan.go:

```
go build -o goScan/goScan goScan/goScan.go && goScan/goScan | xsv table
```
