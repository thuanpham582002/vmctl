#!/bin/zsh
go build
vmctl completion zsh > _vmctl
# Đường dẫn đến thư mục hoàn thành của Oh My Zsh
COMPLETION_DIR="$HOME/.oh-my-zsh/completions"

# Tạo thư mục nếu nó chưa tồn tại
mkdir -p "$COMPLETION_DIR"

# Sao chép tệp _limactl vào thư mục hoàn thành
cp ./_vmctl "$COMPLETION_DIR/"

# Tải lại cấu hình Oh My Zsh
source "$HOME/.zshrc"

exec zsh
