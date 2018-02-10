# Copyright 2018 Pax Automa Systems, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

UPGRADER_FILES=$(shell find components/upgrader/ -name "*.go")

.PHONY: upgrader-novm
upgrader-novm: iso/controller/airootfs/usr/bin/upgrader

iso/controller/airootfs/usr/bin/upgrader: components/teamster/pkg/teamster/teamster.pb.go $(UPGRADER_FILES) vendor
	mkdir -p $(dir $@)
	go build -v -o $@ \
		./components/upgrader/cmd/check.go

clean: clean-upgrader

clean-upgrader:
	rm -f iso/controller/airootfs/usr/bin/upgrader
