# Self-Hosted Runner Setup Guide

Tất cả CI/CD workflows của Host Anything chạy trên **self-hosted runner** với label `debian`.
Trang này hướng dẫn cách setup runner trên máy chủ Debian của bạn.

## Yêu cầu hệ thống

- Debian 11 (Bullseye) hoặc mới hơn
- 2 CPU cores, 4GB RAM tối thiểu
- 20GB disk trống
- Quyền `sudo`

## 1. Cài đặt dependencies

```bash
# Go (phiên bản phải khớp với go.mod)
wget https://go.dev/dl/go1.23.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.23.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc

# Node.js 20 LTS
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs

# Build tools (cần cho -race flag)
sudo apt-get install -y gcc libc6-dev

# Docker (cho release pipeline)
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker $USER

# dpkg-deb (thường đã có sẵn trên Debian)
sudo apt-get install -y dpkg

# golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
  | sh -s -- -b /usr/local/bin v1.59.0

# markdownlint-cli2
sudo npm install -g markdownlint-cli2

# trivy
wget -qO - https://aquasecurity.github.io/trivy-repo/deb/public.key | \
  gpg --dearmor | sudo tee /usr/share/keyrings/trivy.gpg > /dev/null
echo "deb [signed-by=/usr/share/keyrings/trivy.gpg] https://aquasecurity.github.io/trivy-repo/deb generic main" | \
  sudo tee /etc/apt/sources.list.d/trivy.list
sudo apt-get update && sudo apt-get install -y trivy

# nancy (Go CVE checker)
go install github.com/sonatype-nexus-community/nancy@latest

# govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest
```

## 2. Dăng ký GitHub Actions Runner

Truy cập GitHub repo → **Settings** → **Actions** → **Runners** → **New self-hosted runner**.

Chọn **Linux** và làm theo hướng dẫn. Ví dụ:

```bash
mkdir actions-runner && cd actions-runner

# Tải runner (thay URL bằng URL từ GitHub)
curl -o actions-runner-linux-x64.tar.gz -L \
  https://github.com/actions/runner/releases/download/v2.317.0/actions-runner-linux-x64-2.317.0.tar.gz
tar xzf ./actions-runner-linux-x64.tar.gz

# Cấu hình (thay TOKEN và REPO_URL)
./config.sh \
  --url https://github.com/YOUR_ORG/host-anything \
  --token YOUR_TOKEN \
  --labels "debian" \
  --name "ha-runner-01" \
  --work _work
```

> [!IMPORTANT]
> Label **`debian`** là bắt buộc — tất cả workflows đều target `runs-on: [self-hosted, debian]`.

## 3. Chạy runner như systemd service

```bash
# Cài đặt service
sudo ./svc.sh install

# Khởi động
sudo ./svc.sh start

# Kiểm tra status
sudo ./svc.sh status
```

Service sẽ tự khởi động lại sau reboot.

## 4. Bảo mật runner

```bash
# Tạo user riêng cho runner (không dùng root)
sudo useradd -m -s /bin/bash github-runner
sudo usermod -aG docker github-runner

# Giới hạn quyền truy cập repo
# Runner chỉ nên có quyền đọc/ghi vào thư mục _work
```

> [!CAUTION]
> **Không** đăng ký self-hosted runner trên public repo mà không có biện pháp bảo vệ.
> Self-hosted runners có thể bị exploit qua malicious PRs.
> Host Anything là private project — đây là cấu hình an toàn.

## 5. Scaling (tùy chọn)

Nếu cần nhiều runner song song:

```bash
# Đăng ký thêm runner với tên khác (cùng label)
./config.sh \
  --url https://github.com/YOUR_ORG/host-anything \
  --token YOUR_TOKEN \
  --labels "debian" \
  --name "ha-runner-02" \
  --work _work_02
```

GitHub Actions sẽ tự cân bằng tải giữa các runners có cùng label.
