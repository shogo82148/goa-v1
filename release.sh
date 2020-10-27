#!/bin/bash

set -eux

CURRENT=$(cd "$(dirname "$0")" && pwd)
VERSION=$1
MAJOR=$(echo "$VERSION" | cut -d. -f1)
MINOR=$(echo "$VERSION" | cut -d. -f2)
PATCH=$(echo "$VERSION" | cut -d. -f3)

cd "$CURRENT"
git switch main

cat <<EOF > version/version.go
package version

const (
	// Major version number
	Major = ${MAJOR}
	// Minor version number
	Minor = ${MINOR}
	// Build version number
	Build = ${PATCH}
)
EOF
perl -i -pe "s(Current Release: \`v[.0-9]+\`)(Current Release: \`v$VERSION\`)" README.md
git add version/version.go README.md
git commit -m "bump up v$VERSION"
git tag "v$VERSION"
git push origin main
git push --tags
