./build.sh

cd testdata/unstable/
go build
../../gostable .

cd ../stable
go build
../../gostable .

cd ../..