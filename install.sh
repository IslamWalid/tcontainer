set -e
wget https://busybox.net/downloads/binaries/1.35.0-x86_64-linux-musl/busybox
chmod a+x busybox
mkdir -p rootfs/bin rootfs/proc rootfs/dev
./busybox --install rootfs/bin
mv rootfs /tmp
rm busybox
