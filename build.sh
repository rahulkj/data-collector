rm -rf staging/bin
mkdir -p staging/bin

go clean
go get github.com/rahulkj/server
go get github.com/mitchellh/gox

./bin/gox -osarch="linux/amd64" github.com/rahulkj/server

if [[ $? -ne 0 ]]; then
	./bin/gox -build-toolchain
	./bin/gox -osarch="linux/amd64" github.com/rahulkj/server
fi	

mv server_linux_amd64 staging/bin/server
cp -r web staging/

cd staging
gcf push data-collector -b http://github.com/ryandotsmith/null-buildpack.git -c ./bin/server
