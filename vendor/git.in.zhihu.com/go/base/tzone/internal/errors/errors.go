package errors

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"git.in.zhihu.com/go/base/telemetry"
)

func WrapError(err error) telemetry.Error {
	if transErr, ok := err.(thrift.TTransportException); ok {
		var input telemetry.Error
		errTypes := map[int]string{
			1: "NotOpen",
			2: "AleradyOpen",
			3: "TimeOut",
			4: "EndOfFile",
		}
		if class, ok := errTypes[transErr.TypeId()]; ok {
			input = telemetry.WrapErr(transErr, class)
		} else {
			input = telemetry.WrapErrWithUnknownClass(transErr.Err())
		}
		return input
	}
	if protoErr, ok := err.(thrift.TProtocolException); ok {
		var input telemetry.Error
		errTypes := map[int]string{
			1: "InvalidData",
			2: "MegativeSize",
			3: "SizeLimit",
			4: "BadVersion",
			5: "NotImplemented",
			6: "DepthLimit",
		}
		if class, ok := errTypes[protoErr.TypeId()]; ok {
			input = telemetry.WrapErr(protoErr, class)
		} else {
			input = telemetry.WrapErrWithUnknownClass(protoErr)
		}
		return input
	}
	if appErr, ok := err.(thrift.TApplicationException); ok {
		var input telemetry.Error
		errTypes := map[int32]string{
			1:  "UnknownMethod",
			2:  "InvalidMessageType",
			3:  "WrongMethodName",
			4:  "BadSequenceID",
			5:  "MissingResult",
			6:  "UnknownInternalError",
			7:  "UnknownProtocolError",
			8:  "InvalidTransform",
			9:  "InvalidProtocol",
			10: "UnsupportedClientType",
		}
		if class, ok := errTypes[appErr.TypeId()]; ok {
			input = telemetry.WrapErr(appErr, class)
		} else {
			input = telemetry.WrapErrWithUnknownClass(appErr)
		}
		return input
	}
	return telemetry.WrapErrWithUnknownClass(err)
}
