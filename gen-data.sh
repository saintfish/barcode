go-bindata -pkg=barcode -o=./data.go -prefix=data data \
&& gofmt -r 'Asset -> _asset' -w data.go