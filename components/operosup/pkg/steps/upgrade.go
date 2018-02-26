package steps

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	losetup "github.com/freddierice/go-losetup"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/paxautoma/operos/components/common/gatekeeper"
	"github.com/paxautoma/operos/components/operosup/pkg/client"
)

type UpgradeFlags struct {
	DownloadOnly bool
	UpgradeOnly  bool
	UpgradeFile  string

	RootPath string
}

func DoUpgrade(dialer *client.Dialer, flags *UpgradeFlags) error {
	var (
		err         error
		upgradeFile string
	)

	if (flags.DownloadOnly || flags.UpgradeOnly) && flags.UpgradeFile != "" {
		upgradeFile = flags.UpgradeFile
	} else {
		tmpDir, err := ioutil.TempDir("", "operos-upgrade")
		if err != nil {
			return errors.Wrap(err, "failed to create temporary directory")
		}
		upgradeFile = path.Join(tmpDir, "upgrade.iso")
	}

	if !flags.UpgradeOnly {
		res, err := getUpgradeInfo(dialer)
		if err != nil {
			return errors.Wrap(err, "failed to query upgrade data from teamster")
		}

		log.Infof("upgrading to v%s", res.GetAvailableVersion())

		url := res.GetUrl()
		if _, err = downloadFile(upgradeFile, url); err != nil {
			return errors.Wrapf(err, "failed to download file %s", url)
		}
	}

	if flags.DownloadOnly {
		return nil
	}

	rootPath := flags.RootPath
	if rootPath == "" {
		if rootPath, err = os.Getwd(); err != nil {
			return errors.Wrap(err, "failed to get current working directory")
		}
	}

	upgradeMount := path.Join(rootPath, "/run/archiso/upgrademnt")

	unmount, err := mountISO(upgradeMount, upgradeFile)
	if err != nil {
		return errors.Wrapf(err, "failed to mount file %s", upgradeFile)
	}
	defer unmount()

	log.Infof("mounted %s at %s", upgradeFile, upgradeMount)

	if err := runUpgradeScript(upgradeMount, rootPath); err != nil {
		return errors.Wrap(err, "failed while running upgrade script")
	}

	return nil
}

func getUpgradeInfo(dialer *client.Dialer) (*gatekeeper.UpgradeCheckResp, error) {
	return &gatekeeper.UpgradeCheckResp{
		AvailableVersion: "0.3.x",
		Url:              "http://192.168.33.1:8000/operos-installer-dev-0.3.x.iso",
		// Url: "http://arch-pkgs.paxautoma.com/calicoctl-1.3.0-1-x86_64.pkg.tar.xz",
	}, nil
}

func mountISO(targetPath, fileName string) (func(), error) {
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return nil, errors.Wrap(err, "failed to create upgrade mount directory")
	}

	dev, err := losetup.Attach(fileName, 0, true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to determine free loop device")
	}

	if err := syscall.Mount(dev.Path(), targetPath, "iso9660", syscall.MS_RDONLY, ""); err != nil {
		return nil, errors.Wrap(err, "failed to mount loop device")
	}

	return func() {
		syscall.Unmount(targetPath, 0)
		dev.Detach()
	}, nil
}

func runUpgradeScript(sourcePath, rootPath string) error {
	scriptPath := path.Join(sourcePath, "/operos/installfiles/upgrade.sh")

	cmd := exec.Command(
		scriptPath,
		"-s", strings.TrimRight(sourcePath, "/"),
		"-r", strings.TrimRight(rootPath, "/"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}
