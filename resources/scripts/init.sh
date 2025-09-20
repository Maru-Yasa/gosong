#!/bin/bash
set -e

REPO="Maru-Yasa/gosong"
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
  x86_64) ARCH=amd64 ;;
  aarch64) ARCH=arm64 ;;
esac

sudo mkdir -p /opt/gosong

ASSET_URL=$(curl -s "https://api.github.com/repos/$REPO/releases" \
  | grep "browser_download_url" \
  | grep "$OS-$ARCH.zip" \
  | head -n 1 \
  | cut -d '"' -f 4)

echo "Downloading: $ASSET_URL"
curl -L "$ASSET_URL" -o /tmp/gosong-$OS-$ARCH.zip
sudo unzip -o /tmp/gosong-$OS-$ARCH.zip -d /opt/gosong

if [ -f /opt/gosong/gosong-$OS-$ARCH ]; then
  sudo mv /opt/gosong/gosong-$OS-$ARCH /opt/gosong/gosong
fi

sudo chmod +x /opt/gosong/gosong

mkdir -p /var/lib/gosong/

if [ -f /tmp/gosong.service ]; then
  sudo mv /tmp/gosong.service /etc/systemd/system/gosong.service
  sudo systemctl daemon-reload
  sudo systemctl enable gosong
fi

if systemctl list-unit-files | grep -q gosong.service; then
  sudo systemctl restart gosong || sudo systemctl start gosong
else
  echo "gosong.service not installed, please check!"
fi
