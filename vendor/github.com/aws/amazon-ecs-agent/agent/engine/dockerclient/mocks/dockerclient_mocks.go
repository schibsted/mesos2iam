// Copyright 2015-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/aws/amazon-ecs-agent/agent/engine/dockerclient (interfaces: Factory)

package mock_dockerclient

import (
	dockerclient "github.com/aws/amazon-ecs-agent/agent/engine/dockerclient"
	dockeriface "github.com/aws/amazon-ecs-agent/agent/engine/dockeriface"
	gomock "github.com/golang/mock/gomock"
)

// Mock of Factory interface
type MockFactory struct {
	ctrl     *gomock.Controller
	recorder *_MockFactoryRecorder
}

// Recorder for MockFactory (not exported)
type _MockFactoryRecorder struct {
	mock *MockFactory
}

func NewMockFactory(ctrl *gomock.Controller) *MockFactory {
	mock := &MockFactory{ctrl: ctrl}
	mock.recorder = &_MockFactoryRecorder{mock}
	return mock
}

func (_m *MockFactory) EXPECT() *_MockFactoryRecorder {
	return _m.recorder
}

func (_m *MockFactory) FindAvailableVersions() []dockerclient.DockerVersion {
	ret := _m.ctrl.Call(_m, "FindAvailableVersions")
	ret0, _ := ret[0].([]dockerclient.DockerVersion)
	return ret0
}

func (_mr *_MockFactoryRecorder) FindAvailableVersions() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "FindAvailableVersions")
}

func (_m *MockFactory) GetClient(_param0 dockerclient.DockerVersion) (dockeriface.Client, error) {
	ret := _m.ctrl.Call(_m, "GetClient", _param0)
	ret0, _ := ret[0].(dockeriface.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockFactoryRecorder) GetClient(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetClient", arg0)
}

func (_m *MockFactory) GetDefaultClient() (dockeriface.Client, error) {
	ret := _m.ctrl.Call(_m, "GetDefaultClient")
	ret0, _ := ret[0].(dockeriface.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockFactoryRecorder) GetDefaultClient() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDefaultClient")
}
