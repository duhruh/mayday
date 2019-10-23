package pkg

//go:generate protoc -I ../ -I /opt/google/googleapis/google -I /opt/proto/include  --go_out=plugins=grpc:../ ../proto/mayday.proto
