#!/bin/bash

# Verify that the correct license block is present in all Go source
# files.
read -r -d '' EXPECTED <<EndOfLicense
// Copyright 2016 The go-libvirt Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
EndOfLicense

# Scan each Go source file for license.
EXIT=0
GOFILES=$(find . -name "*.go")

for FILE in $GOFILES; do
	BLOCK=$(head -n 14 $FILE)

	if [ "$BLOCK" != "$EXPECTED" ]; then
		echo "file missing license: $FILE"
		EXIT=1
	fi
done

exit $EXIT
