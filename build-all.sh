echo "Compiling proto files"
./build-proto.sh && \
./build-inject-tag.sh && \
./build-reverse-proxy.sh 