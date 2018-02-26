#!/bin/bash -e
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

scriptpath="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd -P)"

srcdir=
rootdir=""
noreboot=

# Find the root directory of the ISO by traversing up the directory stack until
# a directory is found that contains the file "operos/version". 
find_src_root() {
    local start=$1
    local cur=$start

    while [[ ${cur} != "/" ]]; do
        if [ -e ${cur}/operos/version ]; then
            echo $cur
            return 0
        fi
        cur=$(dirname $cur)
    done

    echo "could not find ISO root" >&2
    return 1
}

copy_files() {
    echo "> Copying files"

    mkdir -p ${rootdir}/run/archiso/bootmnt/operos-${new_version}/
    cp -af ${srcdir}/operos/x86_64 ${rootdir}/run/archiso/bootmnt/operos-${new_version}/
}

update_bootloader() {
    if [ -d /sys/firmware/efi/efivars ]; then
        update_efi
    else
        update_bios
    fi
}

update_efi() {
    echo "> updating EFI configuration"

    # kernel, ucode, initcpio image
    mkdir -p ${rootdir}/efi/EFI/operos-${new_version}
    cp ${srcdir}/operos/boot/intel_ucode.img ${rootdir}/efi/EFI/operos-${new_version}/intel_ucode.img
    cp ${srcdir}/operos/boot/intel_ucode.LICENSE ${rootdir}/efi/EFI/operos-${new_version}/intel_ucode.LICENSE
    cp ${srcdir}/operos/boot/x86_64/archiso.img ${rootdir}/efi/EFI/operos-${new_version}/archiso.img
    cp ${srcdir}/operos/boot/x86_64/vmlinuz ${rootdir}/efi/EFI/operos-${new_version}/vmlinuz.efi

    # menu entries
    mkdir -p ${rootdir}/efi/loader/entries
    cat ${srcdir}/operos/installfiles/efi-entry.templ | \
        sed "s/%CONTROLLER_DISK%/${CONTROLLER_DISK//\//\\\/}/g;
             s/%OPEROS_VERSION%/${new_version}/g;" > \
        ${rootdir}/efi/loader/entries/operos-${new_version}.cfg

    cat > ${rootdir}/efi/loader/loader.conf <<EOF
default operos-${new_version}
timeout 5
EOF
}

update_bios() {
    echo "> Updating SYSLINUX configuration"

    mkdir -p ${rootdir}/boot/syslinux-controller
    mkdir -p ${rootdir}/boot/operos-${new_version}

    # upgrade memtest
    cp ${srcdir}/operos/boot/memtest ${rootdir}/boot/memtest
    cp ${srcdir}/operos/boot/memtest.COPYING ${rootdir}/boot/memtest.COPYING

    # kernel, ucode, initcpio image
    cp ${srcdir}/operos/boot/intel_ucode.img ${rootdir}/boot/operos-${new_version}/intel_ucode.img
    cp ${srcdir}/operos/boot/intel_ucode.LICENSE ${rootdir}/boot/operos-${new_version}/intel_ucode.LICENSE
    cp ${srcdir}/operos/boot/x86_64/archiso.img ${rootdir}/boot/operos-${new_version}/archiso.img
    cp ${srcdir}/operos/boot/x86_64/vmlinuz ${rootdir}/boot/operos-${new_version}/vmlinuz

    # syslinux modules and config
    rm -rf ${rootdir}/boot/syslinux-*
    cp -af ${srcdir}/operos/boot/syslinux-* ${rootdir}/boot

    # menu entries
    cat ${srcdir}/operos/installfiles/syslinux-entry.templ | \
        sed "s/%CONTROLLER_DISK%/${CONTROLLER_DISK//\//\\\/}/g;
             s/%OPEROS_VERSION%/${new_version}/g;" > \
        ${rootdir}/boot/operos-${new_version}/entry.cfg

    cat > ${rootdir}/boot/syslinux-controller/entries.cfg <<EOF
DEFAULT operos-${new_version}
TIMEOUT 50

INCLUDE operos-${new_version}/entry.cfg
INCLUDE operos-${cur_version}/entry.cfg
EOF
}

update_settings() {
    echo "> Updating settings"
    #sed -i "s/^\(OPEROS_VERSION=\).*/\1\"${new_version}\"/g;" ${rootdir}/etc/paxautoma/settings
    #operoscfg -s OPEROS_VERSION ${new_version}
}

cleanup_files() {
    echo "> Cleaning up old files"
    local release

    for releasedir in ${rootdir}/run/archiso/bootmnt/operos-* ${rootdir}/boot/operos-* ${rootdir}/efi/EFI/operos-*; do
        if [[ ! -d ${releasedir} ]]; then
            continue
        fi

        release=$(basename $releasedir)
        if [[ "${release}" != "operos-${cur_version}" ]] &&
           [[ "${release}" != "operos-${new_version}" ]]; then
           echo "  - removing ${releasedir}"
           rm -rf ${releasedir}
        fi
    done
}

while getopts ":s:r:bh" opt; do
    case $opt in
    s) srcdir=$OPTARG ;;
    r) rootdir=$OPTARG ;;
    b) noreboot=yes ;;
    h)
        echo "Usage: $0 [-s <srcdir>] [-r <rootdir>] [-b] [-h] "
        echo "This script updates the Operos version on the controller."
        echo "It is not meant to be run directly."
        exit 0
        ;;
    *)
        echo "Unknown option: -$OPTARG" >&2
        exit 1
        ;;
    esac
done

set -a
. /etc/paxautoma/settings
set +a

if [[ -z ${srcdir} ]]; then
    srcdir=$(find_src_root ${scriptpath})
fi

cur_version=${OPEROS_VERSION}
new_version=$(cat ${srcdir}/operos/version)

if [[ ${cur_version} = ${new_version} ]]; then
    echo "Version is already ${new_version}"
    exit 1
fi

echo "Upgrading Operos: v${cur_version} -> v${new_version}"
echo ""

copy_files
update_bootloader
update_settings
cleanup_files

echo "Upgrade completed"

if [[ -z "${noreboot}" ]]; then
    echo "Rebooting"
    reboot
fi
