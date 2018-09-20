package main

import (
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/TykTechnologies/tyk-protobuf/bindings/go"
	"google.golang.org/grpc"
)

const (
	listenAddress = ":9111"
)

func main() {
	lis, err := net.Listen("tcp", listenAddress)
	fatalOnError(err, "failed to start tcp listener")

	logrus.Infof("starting grpc middleware on %s", listenAddress)
	s := grpc.NewServer()
	coprocess.RegisterDispatcherServer(s, &Dispatcher{})
	s.Serve(lis)
}

func fatalOnError(err error, msg string) {
	if err != nil {
		logrus.WithError(err).Fatal(msg)
	}
}
