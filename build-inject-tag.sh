echo Injecting tag to go struct
echo ./proto/listing/$(ls proto/* | grep pb.go)
protoc-go-inject-tag -input=./proto/listing/$(ls proto/* | grep pb.go)