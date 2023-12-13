
.PHONY: gen
gen:
	go build -o protoc-gen-fetch-types main.go
	buf generate

.PHONY: protocol 
protocol:
	buf generate --template buf.gen.go.yaml --config buf.go.yaml

