./build.sh

cd testdata/unstable/
# https://github.com/golang/go/issues/37054#issuecomment-1648091135
# this is a workaround to resolve dependencies of the test code
#go mod tidy && go mod vendor && go mod verify
go build
GOSTABLE_OUT=`../../gostable .`

echo "before"
echo $GOSTABLE_OUT
echo "after"

cd ../stable
#go mod tidy && go mod vendor && go mod verify
go build
../../gostable .


cd ..
