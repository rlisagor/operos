/*
Copyright 2018 Pax Automa Systems, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package steps

import (
	"os"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

func downloadFile(targetPath, url string) (string, error) {
	client := grab.NewClient()
	req, _ := grab.NewRequest(targetPath, url)

	log.Infof("downloading file: %s", url)
	resp := client.Do(req)

	var t *time.Ticker
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		t = time.NewTicker(500 * time.Millisecond)
	} else {
		t = time.NewTicker(5 * time.Second)
	}
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			log.Debugf("transferred %v / %v bytes (%.2f%%)",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress())
		case <-resp.Done:
			break Loop
		}
	}

	if err := resp.Err(); err != nil {
		return "", errors.Wrap(err, "download failed")
	}

	log.Infof("download saved to %v", resp.Filename)
	return resp.Filename, nil
}
