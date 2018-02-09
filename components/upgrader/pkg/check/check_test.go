package check_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/paxautoma/operos/components/common/gatekeeper"
	"github.com/paxautoma/operos/components/upgrader/pkg/check"
)

const mockServerAddr = "127.0.0.33:8888"

type MockClient struct {
	AvailableVersion string
	URLs             []string
}

func (mc *MockClient) AuthorizeDiagUpload(ctx context.Context, in *gatekeeper.AuthorizeDiagUploadReq, opts ...grpc.CallOption) (*gatekeeper.AuthorizeDiagUploadResp, error) {
	return &gatekeeper.AuthorizeDiagUploadResp{}, nil
}

func (mc *MockClient) UpgradeCheck(ctx context.Context, in *gatekeeper.UpgradeCheckReq, opts ...grpc.CallOption) (*gatekeeper.UpgradeCheckResp, error) {
	return &gatekeeper.UpgradeCheckResp{
		AvailableVersion: mc.AvailableVersion,
		DownloadUrls:     mc.URLs,
	}, nil
}

func (mc *MockClient) Close() error {
	return nil
}

func NewMockClient(availableVersion string, urls []string) (io.Closer, gatekeeper.GatekeeperClient, error) {
	mc := &MockClient{
		AvailableVersion: availableVersion,
		URLs:             urls,
	}
	return mc, mc, nil
}

func doTest(t *testing.T, c check.UpgradeCheck, shouldDownload bool) {
	tmpDir, err := ioutil.TempDir("", "operos-upgrade-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	c.DialFunc = func(addr string, noTLS bool) (io.Closer, gatekeeper.GatekeeperClient, error) {
		return NewMockClient("0.0.1", []string{
			fmt.Sprintf("http://%s/testfile.tar.gz", mockServerAddr),
		})
	}
	c.TargetPath = tmpDir

	err = c.DoCheck()
	require.NoError(t, err)

	tarFile := path.Join(tmpDir, "testfile.tar.gz")
	testFile := path.Join(tmpDir, "testfile")

	if shouldDownload {
		require.FileExists(t, tarFile)
		require.FileExists(t, testFile)
		data, err := ioutil.ReadFile(testFile)
		require.NoError(t, err)
		require.Equal(t, "operos upgrade contents\n", string(data))
	} else {
		requireFileNotExists(t, tarFile)
	}

}

func TestUpgradeCheck(t *testing.T) {
	srv := &http.Server{
		Addr:    mockServerAddr,
		Handler: http.FileServer(http.Dir("testdata")),
	}

	listener, err := net.Listen("tcp", mockServerAddr)
	require.NoError(t, err)

	go func() {
		srv.Serve(listener)
	}()
	defer srv.Shutdown(context.Background())

	tests := []struct {
		Name           string
		Spec           check.UpgradeCheck
		ShouldDownload bool
	}{
		{
			"SameVersion_DoesNotDownloadUpgrade",
			check.UpgradeCheck{Version: "0.0.1"},
			false,
		},
		{
			"DifferentVersion_DownloadsUpgrade",
			check.UpgradeCheck{Version: "0.0.2"},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			doTest(t, test.Spec, test.ShouldDownload)
		})
	}

}

func requireFileNotExists(t *testing.T, filename string, msgAndArgs ...interface{}) {
	_, err := os.Lstat(filename)
	require.True(t, os.IsNotExist(err), msgAndArgs...)
}
