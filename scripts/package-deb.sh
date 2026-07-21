#!/bin/bash
set -euo pipefail

ARCH=${1:-amd64}
VERSION="0.1.0"
PKG_DIR="build/hostanything_${VERSION}_${ARCH}"

echo "Creating packaging directory structure..."
mkdir -p "$PKG_DIR/DEBIAN"
mkdir -p "$PKG_DIR/usr/bin"
mkdir -p "$PKG_DIR/lib/systemd/system"
mkdir -p "$PKG_DIR/etc/hostanything"
mkdir -p "$PKG_DIR/var/lib/hostanything"
mkdir -p "$PKG_DIR/var/log/hostanything"
mkdir -p "$PKG_DIR/usr/share/hostanything/web"

echo "Copying files..."
# Assuming binary is built and frontend is compiled
cp bin/hostanything-linux-$ARCH "$PKG_DIR/usr/bin/hostanything"
cp -r web/dist/* "$PKG_DIR/usr/share/hostanything/web/"

cat <<EOF > "$PKG_DIR/DEBIAN/control"
Package: hostanything
Version: $VERSION
Architecture: $ARCH
Maintainer: Host Anything Contributors <hello@host-anything.dev>
Description: Deploy and manage services with ease.
Depends: systemd, fail2ban
EOF

cat <<EOF > "$PKG_DIR/lib/systemd/system/hostanything.service"
[Unit]
Description=Host Anything Server
After=network.target

[Service]
ExecStart=/usr/bin/hostanything
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOF

cat <<EOF > "$PKG_DIR/DEBIAN/postinst"
#!/bin/sh
set -e
systemctl daemon-reload
systemctl enable hostanything
systemctl restart hostanything
EOF
chmod +x "$PKG_DIR/DEBIAN/postinst"

cat <<EOF > "$PKG_DIR/DEBIAN/prerm"
#!/bin/sh
set -e
systemctl stop hostanything
systemctl disable hostanything
EOF
chmod +x "$PKG_DIR/DEBIAN/prerm"

echo "Building .deb package..."
dpkg-deb --build "$PKG_DIR"

echo "Done! Package created."
