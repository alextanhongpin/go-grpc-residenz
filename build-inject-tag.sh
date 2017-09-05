echo "Injecting tag to go struct"
protoc-go-inject-tag -input=./proto/listing/$(ls proto/* | grep pb.go)