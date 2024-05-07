// Package mounter
package mounter

import (
	"os"
	"testing"

	"github.com/IBM/ibm-object-csi-driver/pkg/utils"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewS3fsMounter_Success(t *testing.T) {
	// Mock the secretMap and mountOptions
	secretMap := map[string]string{
		"cosEndpoint":   "test-endpoint",
		"locConstraint": "test-loc-constraint",
		"bucketName":    "test-bucket-name",
		"objPath":       "test-obj-path",
		"accessKey":     "test-access-key",
		"secretKey":     "test-secret-key",
		"apiKey":        "test-api-key",
		"kpRootKeyCRN":  "test-kp-root-key-crn",
	}

	mountOptions := []string{"opt1=val1", "opt2=val2"}

	mounter, err := NewS3fsMounter(secretMap, mountOptions, utils.NewMockStatsUtilsImpl(utils.MockStatsUtilsFuncStruct{}))
	if err != nil {
		t.Errorf("NewS3fsMounter failed: %v", err)
	}

	s3fsMounter, ok := mounter.(*S3fsMounter)
	if !ok {
		t.Errorf("NewS3fsMounter() failed to return an instance of s3fsMounter")
	}

	if s3fsMounter.BucketName != secretMap["bucketName"] {
		t.Errorf("Expected bucketName: %s, got: %s", secretMap["bucketName"], s3fsMounter.BucketName)
	}
	if s3fsMounter.ObjPath != secretMap["objPath"] {
		t.Errorf("Expected objPath:   %s, got %s ", secretMap["objPath"], s3fsMounter.ObjPath)
	}
	if s3fsMounter.EndPoint != secretMap["endPoint"] {
		t.Errorf("Expected endPoint: %s, got %s ", secretMap["cosEndpoint"], s3fsMounter.EndPoint)
	}
	if s3fsMounter.LocConstraint != secretMap["locConstraint"] {
		t.Errorf("Expected locationConstraint: %s, got %s ", secretMap["locConstraint"], s3fsMounter.LocConstraint)
	}
}

// Fake VolumeStatsUtils
type FakeVolumeStatsUtils struct {
}

func (su *FakeVolumeStatsUtils) FSInfo(path string) (int64, int64, int64, int64, int64, int64, error) {
	if path == "some/path" {
		return 0, 0, 0, 0, 0, 0, status.Error(codes.NotFound, "volume not found on some/path")
	}
	return 1, 1, 1, 1, 1, 1, nil
}

func (su *FakeVolumeStatsUtils) CheckMount(targetPath string) (bool, error) {
	return true, nil
}

func (su *FakeVolumeStatsUtils) FuseUnmount(path string) error {
	return nil
}

func (su *FakeVolumeStatsUtils) FuseMount(path string, comm string, args []string) error {
	return nil
}

func Test_Mount_Positive(t *testing.T) {
	secretMap := map[string]string{
		"cosEndpoint":        "test-endpoint",
		"locationConstraint": "test-loc-constraint",
		"bucketName":         "test-bucket-name",
		"objPath":            "test-obj-path",
		"accessKey":          "test-access-key",
		"secretKey":          "test-secret-key",
		"apiKey":             "test-api-key",
		"kpRootKeyCRN":       "test-kp-root-key-crn",
	}
	mounter, err := NewS3fsMounter(secretMap, []string{"mountOption1", "mountOption2"}, utils.NewMockStatsUtilsImpl(utils.MockStatsUtilsFuncStruct{}))
	if err != nil {
		t.Fatalf("NewS3fsMounter() returned an unexpected error: %v", err)
	}
	s3fsMounter, ok := mounter.(*S3fsMounter)
	if !ok {
		t.Fatal("NewS3fsMounter() did not return a s3fsMounter")
	}

	//s3fsMounter.StatsUtils = utils.NewMockStatsUtilsImpl(utils.MockStatsUtilsFuncStruct{})

	target := "/tmp/test-mount"

	err = s3fsMounter.Mount("source", target)
	if err != nil {
		t.Errorf("S3fsMounter_Mount() returned an unexpected error: %v", err)
	}
}

func Test_Unmount_Positive(t *testing.T) {
	secretMap := map[string]string{
		"cosEndpoint":        "test-endpoint",
		"locationConstraint": "test-loc-constraint",
		"bucketName":         "test-bucket-name",
		"objPath":            "test-obj-path",
		"accessKey":          "test-access-key",
		"secretKey":          "test-secret-key",
		"apiKey":             "test-api-key",
		"kpRootKeyCRN":       "test-kp-root-key-crn",
	}
	mounter, _ := NewS3fsMounter(secretMap, []string{"mountOption1", "mountOption2"}, utils.NewMockStatsUtilsImpl(utils.MockStatsUtilsFuncStruct{}))
	s3fsMounter := mounter.(*S3fsMounter)

	target := "/tmp/test-unmount"

	// Creating a directory to simulate a mounted path
	err := os.MkdirAll(target, os.ModePerm)
	if err != nil {
		t.Fatalf("TestS3fsMounter_Unmount() failed to create directory: %v", err)
	}

	err = s3fsMounter.Unmount(target)
	if err != nil {
		t.Errorf("TestS3fsMounter_Unmount() failed to unmount: %v", err)
	}

	// Assert the mount point non-existence
	_, err = os.Stat(target)
	if !os.IsNotExist(err) {
		t.Errorf("TestS3fsMounter_Unmount() failed to remove mount point: %v", err)
	}
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
