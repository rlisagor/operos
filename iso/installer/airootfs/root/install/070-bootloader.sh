#!/bin/bash -xe
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

echo \> Installing bootloader >&3

install_bios() {
    mkdir -p /mnt/boot/operos-${OPEROS_VERSION}

    # memtest (one version for machine - not versioned with Operos)
    cp /run/archiso/bootmnt/operos/boot/memtest /mnt/boot/memtest
    cp /run/archiso/bootmnt/operos/boot/memtest.COPYING /mnt/boot/memtest.COPYING

    # kernel, ucode, initcpio image
    cp /run/archiso/bootmnt/operos/boot/intel_ucode.img /mnt/boot/operos-${OPEROS_VERSION}/intel_ucode.img
    cp /run/archiso/bootmnt/operos/boot/intel_ucode.LICENSE /mnt/boot/operos-${OPEROS_VERSION}/intel_ucode.LICENSE
    cp /run/archiso/bootmnt/operos/boot/x86_64/archiso.img /mnt/boot/operos-${OPEROS_VERSION}/archiso.img
    cp /run/archiso/bootmnt/operos/boot/x86_64/vmlinuz /mnt/boot/operos-${OPEROS_VERSION}/vmlinuz

    # syslinux modules and config
    cp -af /run/archiso/bootmnt/operos/boot/syslinux-* /mnt/boot

    # install MBR
    dd bs=440 conv=notrunc count=1 if=/usr/lib/syslinux/bios/gptmbr.bin of=${CONTROLLER_DISK}

    cat > /mnt/boot/syslinux.cfg <<EOF
PATH /syslinux-controller/
DEFAULT loadconfig

LABEL loadconfig
  CONFIG /syslinux-controller/syslinux.cfg
EOF

    # set up syslinux
    arch-chroot /mnt extlinux --install /boot

    cp /run/archiso/bootmnt/operos/installfiles/syslinux-entry.templ /mnt/boot/operos-${OPEROS_VERSION}/entry.cfg

    cat > /mnt/boot/syslinux-controller/entries.cfg <<EOF
DEFAULT operos-${OPEROS_VERSION}

INCLUDE operos-${OPEROS_VERSION}/entry.cfg
EOF

    sed -i "s/%CONTROLLER_DISK%/${CONTROLLER_DISK//\//\\\/}/g;
            s/%OPEROS_VERSION%/${OPEROS_VERSION}/g;" /mnt/boot/operos-${OPEROS_VERSION}/entry.cfg
}

install_efi() {
    mkdir -p /mnt/efi/EFI/operos-${OPEROS_VERSION}
    cp /run/archiso/bootmnt/operos/boot/intel_ucode.img /mnt/efi/EFI/operos-${OPEROS_VERSION}/intel_ucode.img
    cp /run/archiso/bootmnt/operos/boot/intel_ucode.LICENSE /mnt/efi/EFI/operos-${OPEROS_VERSION}/intel_ucode.LICENSE
    cp /run/archiso/bootmnt/operos/boot/x86_64/archiso.img /mnt/efi/EFI/operos-${OPEROS_VERSION}/archiso.img
    cp /run/archiso/bootmnt/operos/boot/x86_64/vmlinuz /mnt/efi/EFI/operos-${OPEROS_VERSION}/vmlinuz.efi

    # copy loader config
    mkdir /mnt/efi/loader /mnt/efi/loader/entries
    cp /run/archiso/bootmnt/operos/installfiles/efi-entry.templ /mnt/efi/loader/entries/operos-${OPEROS_VERSION}.conf 

    cat > /mnt/efi/loader/loader.conf <<EOF
default operos-%OPEROS_VERSION%
timeout 5
EOF

    find /mnt/efi/loader -name "*.conf" -exec \
        sed -i "s/%CONTROLLER_DISK%/${CONTROLLER_DISK//\//\\\/}/g;
                s/%OPEROS_VERSION%/${OPEROS_VERSION}/g;" {} \;

    # install the systemd-boot efi binaries
    bootctl --path=/mnt/efi install
}

if [ -d /sys/firmware/efi/efivars ]; then
    install_efi
else
    install_bios
fi
