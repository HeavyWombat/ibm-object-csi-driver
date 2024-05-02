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

// Mounter interface defined in mounter.go
// s3fsMounter Implements Mounter
type FakeS3fsMounter struct {
	MountFunc   func(source string, target string) error
	UnmountFunc func(target string) error
}

func fakenewS3fsMounter(secretMap map[string]string, mountOptions []string) (mounter.Mounter, error) {
	var (
		val         string
		check       bool
		accessKey   string
		secretKey   string
		apiKey      string
		fakemounter *mounter.s3fsMounter
	)

	fakemounter = &s3fsMounter{}

	if val, check = secretMap["cosEndpoint"]; check {
		fakemounter.endPoint = val
	}
	if val, check = secretMap["locationConstraint"]; check {
		fakemounter.locConstraint = val
	}
	if val, check = secretMap["bucketName"]; check {
		fakemounter.bucketName = val
	}
	if val, check = secretMap["objPath"]; check {
		fakemounter.objPath = val
	}
	if val, check = secretMap["accessKey"]; check {
		accessKey = val
	}
	if val, check = secretMap["secretKey"]; check {
		secretKey = val
	}
	if val, check = secretMap["apiKey"]; check {
		apiKey = val
	}
	if val, check = secretMap["kpRootKeyCRN"]; check {
		fakemounter.kpRootKeyCrn = val
	}

	if apiKey != "" {
		fakemounter.accessKeys = fmt.Sprintf(":%s", apiKey)
		fakemounter.authType = "iam"
	} else {
		fakemounter.accessKeys = fmt.Sprintf("%s:%s", accessKey, secretKey)
		fakemounter.authType = "hmac"
	}

	klog.Infof("newS3fsMounter args:\n\tbucketName: [%s]\n\tobjPath: [%s]\n\tendPoint: [%s]\n\tlocationConstraint: [%s]\n\tauthType: [%s]kpRootKeyCrn: [%s]",
		fakemounter.bucketName, fakemounter.objPath, fakemounter.endPoint, fakemounter.locConstraint, fakemounter.authType, fakemounter.kpRootKeyCrn)

	updatedOptions, err := mounter.UpdateS3FSMountOptions(mountOptions, secretMap)
	if err != nil {
		klog.Infof("Problems with retrieving secret map dynamically %v", err)
	}
	fakemounter.mountOptions = updatedOptions

	return fakemounter, nil
	//return &FakeS3fsMounter{}, nil
}

func (f *FakeS3fsMounter) Mount(source string, target string) error {
	if f.MountFunc != nil {
		return f.MountFunc(source, target)
	}
	return nil
}

func (f *FakeS3fsMounter) Unmount(target string) error {
	if f.UnmountFunc != nil {
		return f.UnmountFunc(target)
	}
	return nil
}
