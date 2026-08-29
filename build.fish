#!/usr/bin/env fish

function cacti-build
    cd ~/Projects/cactios-iso/rootfs
    sudo sh -c 'find . -print0 | cpio --null -ov --format=newc | gzip -9 > ../iso/boot/initramfs.img'
    cd ..
    grub-mkrescue -o cactios-v0.1.iso iso/
end

function cacti-qemu
    cd ~/Projects/cactios-iso
    qemu-system-x86_64 \
        -m 2G \
        -smp 2 \
        -enable-kvm \
        -cpu host \
        -nic user,model=virtio-net-pci \
        -drive file=cacti-test.qcow2,format=qcow2,if=virtio \
        -cdrom cactios-v0.1.iso \
        -boot d
end

function cacti-test
    cacti-build; and cacti-qemu
end

switch $argv[1]
    case build
        cacti-build
    case qemu
        cacti-qemu
    case test
        cacti-test
    case '*'
        echo "usage: ./build.fish build|qemu|test"
end
