package main

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/TykTechnologies/tyk-protobuf/bindings/go"
)

type Dispatcher struct{}

func (d *Dispatcher) Dispatch(ctx context.Context, object *coprocess.Object) (*coprocess.Object, error) {
	logrus.Infof("receiving object: %v", object)

	switch object.HookName {
	case "Sha256SumHook":
		logrus.Info("Sha256SumHook is called!")
		return Sha256SumHook(object)
	}

	logrus.Warnf("unknown hook: %v", object.HookName)

	return object, nil
}

func (d *Dispatcher) DispatchEvent(ctx context.Context, event *coprocess.Event) (*coprocess.EventReply, error) {
	return &coprocess.EventReply{}, nil
}
