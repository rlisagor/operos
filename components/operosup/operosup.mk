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

OPEROSUP_FILES=$(shell find components/operosup/ -name "*.go")

.PHONY: operosup-novm
operosup-novm: iso/controller/airootfs/usr/bin/operosup

iso/controller/airootfs/usr/bin/operosup: components/teamster/pkg/teamster/teamster.pb.go $(OPEROSUP_FILES) vendor
	mkdir -p $(dir $@)
	go build -v -o $@ ./components/operosup/cmd/main.go

clean: clean-operosup

clean-operosup:
	rm -f iso/controller/airootfs/usr/bin/operosup