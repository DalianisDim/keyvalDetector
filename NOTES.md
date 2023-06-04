```
cobra-cli --config=.cobra.yaml init
go mod init keyvalDetector
go mod tidy
go get github.com/spf13/cobra
go run main.go 
```


Release: 

```
git tag v0.1.0
git push origin v0.1.0
```