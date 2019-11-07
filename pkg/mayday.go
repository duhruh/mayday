package pkg

//go:generate protoc -I ../ -I /opt/google/googleapis/google -I /opt/proto/include  --go_out=plugins=grpc:../ ../proto/mayday.proto

var (
	// GitCommit - holds the value of the commit this code was built at
	GitCommit string
	// BuildTime - holds the time this release was built
	BuildTime string
	// Version - the release version
	Version string
)
