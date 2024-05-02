package mounter

import (
	//"os"
	"testing"

	//"github.com/IBM/ibm-object-csi-driver/pkg/utils"
	"github.com/stretchr/testify/assert"
)

var mounter_s3fs = &s3fsMounter{
	bucketName:    "testBucket",
	objPath:       "testObjPath",
	endPoint:      "testEndPoint",
	locConstraint: "testLocConstraint",
	authType:      "testAuthType",
	accessKeys:    "testAccessKeys",
	kpRootKeyCrn:  "testKpRootKeyCrn",
	mountOptions:  []string{"testOption1", "testOption2"},
}

func Test_Mount_Positive(t *testing.T) {
	mounter := mounter_s3fs

	//expectedPath := "/var/lib/ibmc-s3fs/testTargetPath"
	/*
		expectedArgs := []string{
			"testBucket:/testObjPath",
			"testTargetPath",
			"-o", "sigv2",
			"-o", "use_path_request_style",
			"-o", "passwd_file=/var/lib/ibmc-s3fs/.passwd-s3fs",
			"-o", "url=testEndPoint",
			"-o", "endpoint=testLocConstraint",
			"-o", "allow_other",
			"-o", "mp_umask=002",
			"-o", "ibm_iam_auth",
			"-o", "ibm_iam_endpoint=https://iam.cloud.ibm.com",
		}*/

	err := mounter.Mount("testSource", "testTargetPath")

	assert.NoError(t, err)
}

// Mocks for testing
type fakeVolumeStatsUtils struct{}

func (f *fakeVolumeStatsUtils) FuseUnmount(target string) error {
	return nil
}

func Test_Unmount_Positive(t *testing.T) {
	mounter := mounter_s3fs

	//os.RemoveAll = func(path string) error {
	//	return nil
	//}

	//utils.VolumeStatsUtils = &fakeVolumeStatsUtils{}

	err := mounter.Unmount("testTargetPath")

	assert.NoError(t, err)
}

func TestUpdateS3FSMountOptions(t *testing.T) {

	defaultMountOp := []string{"option1=value1", "option2=value2"}
	secretMap := map[string]string{
		"tmpdir":       "/tmp",
		"use_cache":    "true",
		"gid":          "1001",
		"mountOptions": "additional_option=value3",
	}

	updatedOptions, err := UpdateS3FSMountOptions(defaultMountOp, secretMap)

	assert.NoError(t, err)
	assert.ElementsMatch(t, updatedOptions, []string{
		"option1=value1",
		"option2=value2",
		"tmpdir=/tmp",
		"use_cache=true",
		"gid=1001",
		"uid=1001",
		"additional_option=value3",
	})
}
