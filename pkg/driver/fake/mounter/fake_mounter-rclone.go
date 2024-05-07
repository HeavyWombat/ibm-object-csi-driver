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
	"fmt"
	"github.com/IBM/ibm-object-csi-driver/pkg/mounter"
	"k8s.io/klog/v2"
)

type FakeRcloneMounter struct {
	MountFunc   func(source string, target string) error
	UnmountFunc func(target string) error
}

func FakeNewRcloneMounter(secretMap map[string]string, mountOptions []string) (mounter.Mounter, error) {
	klog.Infof("-newS3fsMounter-")

	var (
		val         string
		check       bool
		accessKey   string
		secretKey   string
		apiKey      string
		fakemounter *mounter.RcloneMounter
	)

	fakemounter = &(mounter.RcloneMounter{})

	if val, check = secretMap["cosEndpoint"]; check {
		fakemounter.EndPoint = val
	}
	if val, check = secretMap["locationConstraint"]; check {
		fakemounter.LocConstraint = val
	}
	if val, check = secretMap["bucketName"]; check {
		fakemounter.BucketName = val
	}
	if val, check = secretMap["objPath"]; check {
		fakemounter.ObjPath = val
	}
	if val, check = secretMap["accessKey"]; check {
		accessKey = val
	}
	if val, check = secretMap["secretKey"]; check {
		secretKey = val
	}
	if val, check = secretMap["kpRootKeyCRN"]; check {
		fakemounter.KpRootKeyCrn = val
	}
	if val, check = secretMap["apiKey"]; check {
		apiKey = val
	}
	if apiKey != "" {
		fakemounter.AccessKeys = fmt.Sprintf(":%s", apiKey)
		fakemounter.AuthType = "iam"
	} else {
		fakemounter.AccessKeys = fmt.Sprintf("%s:%s", accessKey, secretKey)
		fakemounter.AuthType = "hmac"
	}

	if val, check = secretMap["gid"]; check {
		fakemounter.Gid = val
	}
	if secretMap["gid"] != "" && secretMap["uid"] == "" {
		fakemounter.Uid = secretMap["gid"]
	} else if secretMap["uid"] != "" {
		fakemounter.Uid = secretMap["uid"]
	}

	klog.Infof("newRcloneMounter args:\n\tbucketName: [%s]\n\tobjPath: [%s]\n\tendPoint: [%s]\n\tlocationConstraint: [%s]\n\tauthType: [%s]",
		fakemounter.BucketName, fakemounter.ObjPath, fakemounter.EndPoint, fakemounter.LocConstraint, fakemounter.AuthType)

	updatedOptions, err := mounter.UpdateMountOptions(mountOptions, secretMap)

	if err != nil {
		klog.Infof("Problems with retrieving secret map dynamically %v", err)
	}
	fakemounter.MountOptions = updatedOptions

	return fakemounter, nil
}

func (f *FakeRcloneMounter) Mount(source string, target string) error {
	if f.MountFunc != nil {
		return f.MountFunc(source, target)
	}
	return nil
}

func (f *FakeRcloneMounter) Unmount(target string) error {
	if f.UnmountFunc != nil {
		return f.UnmountFunc(target)
	}
	return nil
}
