/**
 * Copyright 2021 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package mounter

import (
	"github.com/IBM/ibm-object-csi-driver/pkg/mounter"
	"k8s.io/klog/v2"
)

const (
	s3fsMounterType   = "s3fs"
	rcloneMounterType = "rclone"
)

type FakeMounter struct {
	MountFunc   func(source, target string) error
	UnmountFunc func(target string) error
}

func (f *FakeMounter) Mount(source, target string) error {
	if f.MountFunc != nil {
		return f.MountFunc(source, target)
	}
	return nil
}

func (f *FakeMounter) Unmount(target string) error {
	if f.UnmountFunc != nil {
		return f.UnmountFunc(target)
	}
	return nil
}

type FakeMounterFactory struct{}

func NewFakeMounterFactory() *FakeMounterFactory {
	return &FakeMounterFactory{}
}

func (f *FakeMounterFactory) NewMounter(attrib map[string]string, secretMap map[string]string, mountFlags []string) (mounter.Mounter, error) {
	klog.Info("-NewMounter-")
	var mounter, val string
	var check bool

	// Select mounter as per storage class
	if val, check = attrib["mounter"]; check {
		mounter = val
	} else {
		// if mounter not set in storage class
		if val, check = secretMap["mounter"]; check {
			mounter = val
		}
	}
	switch mounter {
	case s3fsMounterType:
		return FakeNewS3fsMounter(secretMap, mountFlags)
	case rcloneMounterType:
		return FakeNewRcloneMounter(secretMap, mountFlags)
	default:
		return FakeNewS3fsMounter(secretMap, mountFlags)
	}
}
