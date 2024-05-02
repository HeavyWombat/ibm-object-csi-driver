package mounter

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

type mockCommand struct {
	command string
	args    []string
	start   func() error
	wait    func() error
}

func (m *mockCommand) Start() error {
	return m.start()
}

func (m *mockCommand) Wait() error {
	return m.wait()
}

func TestFuseMount(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		comm          string
		args          []string
		startFunc     func() error
		waitFunc      func() error
		expectedError error
	}{
		{
			name: "SuccessfulMount",
			path: "/mount/path",
			comm: "mount_command",
			args: []string{"arg1", "arg2"},
			startFunc: func() error {
				return nil
			},
			waitFunc: func() error {
				return nil
			},
			expectedError: nil,
		},
		{
			name: "StartCommandError",
			path: "/mount/path",
			comm: "mount_command",
			args: []string{"arg1", "arg2"},
			startFunc: func() error {
				return errors.New("start command error")
			},
			waitFunc: func() error {
				return nil
			},
			expectedError: errors.New("fuseMount: cmd start failed: <mount_command>\nargs: <[arg1 arg2]>\nerror: <start command error>"),
		},
		{
			name: "WaitCommandError",
			path: "/mount/path",
			comm: "mount_command",
			args: []string{"arg1", "arg2"},
			startFunc: func() error {
				return nil
			},
			waitFunc: func() error {
				return errors.New("wait command error")
			},
			expectedError: errors.New("fuseMount: cmd wait failed: <mount_command>\nargs: <[arg1 arg2]>\nerror: <wait command error>"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := fuseMount(tt.path, tt.comm, tt.args)
			if err != nil && tt.expectedError == nil {
				t.Errorf("TestFuseMount(): %s: unexpected error: %v", tt.name, err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("TestFuseMount(): %s: expected error: %v, got nil", tt.name, tt.expectedError)
			}
			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("TestFuseMount(): %s: expected error: %v, got: %v", tt.name, tt.expectedError, err)
			}
		})
	}
}

func TestCheckPath(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		statFunc      func(path string) (os.FileInfo, error)
		expectedExist bool
		expectedError error
	}{
		{
			name: "PathExists",
			path: "/existing/path",
			statFunc: func(path string) (os.FileInfo, error) {
				return nil, nil
			},
			expectedExist: true,
			expectedError: nil,
		},
		{
			name: "PathNotExist",
			path: "/non-existing/path",
			statFunc: func(path string) (os.FileInfo, error) {
				return nil, os.ErrNotExist
			},
			expectedExist: false,
			expectedError: nil,
		},
		{
			name: "PathStatError",
			path: "/error/path",
			statFunc: func(path string) (os.FileInfo, error) {
				return nil, errors.New("stat error")
			},
			expectedExist: false,
			expectedError: errors.New("stat error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//osStat := tt.statFunc
			exist, err := checkPath(tt.path)
			if exist != tt.expectedExist {
				t.Errorf("TestCheckPath(): %s: expected existence: %t, got: %t", tt.name, tt.expectedExist, exist)
			}
			if err != nil && tt.expectedError == nil {
				t.Errorf("TestCheckPath(): %s: unexpected error: %v", tt.name, err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("TestCheckPath(): %s: expected error: %v, got nil", tt.name, tt.expectedError)
			}
			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("TestCheckPath(): %s: expected error: %v, got: %v", tt.name, tt.expectedError, err)
			}
		})
	}
}

func TestIsCorruptedMnt(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		expectedValue bool
	}{
		{
			name:          "NilError",
			err:           nil,
			expectedValue: false,
		},
		{
			name: "PathErrorENOTCONN",
			err: &os.PathError{
				Err: syscall.ENOTCONN,
			},
			expectedValue: true,
		},
		{
			name: "LinkErrorESTALE",
			err: &os.LinkError{
				Err: syscall.ESTALE,
			},
			expectedValue: true,
		},
		{
			name: "SyscallErrorENOTCONN",
			err: &os.SyscallError{
				Err: syscall.ENOTCONN,
			},
			expectedValue: true,
		},
		{
			name:          "OtherError",
			err:           errors.New("other error"),
			expectedValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := isCorruptedMnt(tt.err)
			if value != tt.expectedValue {
				t.Errorf("TestIsCorruptedMnt(): %s: expected value: %t, got: %t", tt.name, tt.expectedValue, value)
			}
		})
	}
}

func TestWritePass(t *testing.T) {
	tests := []struct {
		name          string
		pwFileName    string
		pwFileContent string
		openFileFunc  func(name string, flag int, perm os.FileMode) (*os.File, error)
		expectedError error
	}{
		{
			name:          "WritePassSuccess",
			pwFileName:    "/path/to/passwd",
			pwFileContent: "password",
			openFileFunc: func(name string, flag int, perm os.FileMode) (*os.File, error) {
				return &os.File{}, nil
			},
			expectedError: nil,
		},
		{
			name:          "OpenFileError",
			pwFileName:    "/path/to/passwd",
			pwFileContent: "password",
			openFileFunc: func(name string, flag int, perm os.FileMode) (*os.File, error) {
				return nil, errors.New("open file error")
			},
			expectedError: errors.New("open file error"),
		},
		{
			name:          "WriteError",
			pwFileName:    "/path/to/passwd",
			pwFileContent: "password",
			openFileFunc: func(name string, flag int, perm os.FileMode) (*os.File, error) {
				return &os.File{}, nil
			},
			expectedError: errors.New("write error"),
		},
		{
			name:          "CloseError",
			pwFileName:    "/path/to/passwd",
			pwFileContent: "password",
			openFileFunc: func(name string, flag int, perm os.FileMode) (*os.File, error) {
				return &os.File{}, nil
			},
			expectedError: errors.New("close error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//osOpenFile = tt.openFileFunc
			err := writePass(tt.pwFileName, tt.pwFileContent)
			if err != nil && tt.expectedError == nil {
				t.Errorf("TestWritePass(): %s: unexpected error: %v", tt.name, err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("TestWritePass(): %s: expected error: %v, got nil", tt.name, tt.expectedError)
			}
			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("TestWritePass(): %s: expected error: %v, got: %v", tt.name, tt.expectedError, err)
			}
		})
	}
}

func TestWaitForMount(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		execFunc      func(command string, args ...string) *exec.Cmd
		expectedError error
	}{
		{
			name: "MountPointExists",
			path: "/mount/point",
			execFunc: func(command string, args ...string) *exec.Cmd {
				return &exec.Cmd{}
			},
			expectedError: nil,
		},
		{
			name: "MountPointNotExists",
			path: "/mount/point",
			execFunc: func(command string, args ...string) *exec.Cmd {
				return &exec.Cmd{
					CombinedOutput: func() ([]byte, error) {
						return []byte{}, errors.New("mountpoint not exists")
					},
				}
			},
			expectedError: errors.New("timeout waiting for mount"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execCommand = tt.execFunc
			err := waitForMount(tt.path, 10*time.Second)
			if err != nil && tt.expectedError == nil {
				t.Errorf("TestWaitForMount(): %s: unexpected error: %v", tt.name, err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("TestWaitForMount(): %s: expected error: %v, got nil", tt.name, tt.expectedError)
			}
			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("TestWaitForMount(): %s: expected error: %v, got: %v", tt.name, tt.expectedError, err)
			}
		})
	}
}

func TestNewS3fsMounterFactory(t *testing.T) {
	factory := NewS3fsMounterFactory()
	if factory == nil {
		t.Error("TestNewS3fsMounterFactory(): factory is nil")
	}
}

func TestS3fsMounterFactory_NewMounter(t *testing.T) {
	factory := &S3fsMounterFactory{}
	attrib := map[string]string{"mounter": "s3fs"}
	secretMap := map[string]string{"bucketName": "bucket", "objPath": "path", "cosEndpoint": "endpoint"}
	mounter, err := factory.NewMounter(attrib, secretMap, []string{"flag1", "flag2"})
	if err != nil {
		t.Errorf("TestS3fsMounterFactory_NewMounter(): unexpected error: %v", err)
	}
	if mounter == nil {
		t.Error("TestS3fsMounterFactory_NewMounter(): mounter is nil")
	}
}
