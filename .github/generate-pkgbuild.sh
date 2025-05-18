#!/bin/bash

PKGBUILD="doot-PKGBUILD"
PKGBUILD_BIN="doot-bin-PKGBUILD"
PKGBUILD_GIT="doot-git-PKGBUILD"
HOMEBREW_FORMULA="doot.homebrew-formula.rb"

cp .github/doot.PKGBUILD.template $PKGBUILD
cp .github/doot-bin.PKGBUILD.template $PKGBUILD_BIN
cp .github/doot-git.PKGBUILD.template $PKGBUILD_GIT
cp .github/doot.homebrew-formula.template.rb $HOMEBREW_FORMULA

VERSION=$1
TARBALL_CHECKSUM=$2
LINUX_X86_CHECKSUM=$(sha256sum dist/doot-linux-x86_64 | cut -d ' ' -f 1)
LINUX_ARM64_CHECKSUM=$(sha256sum dist/doot-linux-arm64 | cut -d ' ' -f 1)
DARWIN_X86_CHECKSUM=$(sha256sum dist/doot-darwin-x86_64 | cut -d ' ' -f 1)
DARWIN_ARM64_CHECKSUM=$(sha256sum dist/doot-darwin-arm64 | cut -d ' ' -f 1)

replace_var() {
    local var_name="$1"
    local var_value="${!var_name}"
    if [ -z "$var_value" ]; then
        echo "Error: Variable $var_name is not set."
        exit 1
    fi
    sed -i "s/{{$var_name}}/$var_value/g" "$2"
}
available_vars=(
    "VERSION"
    "TARBALL_CHECKSUM"
    "LINUX_X86_CHECKSUM"
    "LINUX_ARM64_CHECKSUM"
    "DARWIN_X86_CHECKSUM"
    "DARWIN_ARM64_CHECKSUM"
)

for var in "${available_vars[@]}"; do
    replace_var "$var" $PKGBUILD
    replace_var "$var" $PKGBUILD_BIN
    replace_var "$var" $HOMEBREW_FORMULA
done
