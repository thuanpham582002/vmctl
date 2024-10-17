#!/bin/zsh
# Đường dẫn đến thư mục hoàn thành của Oh My Zsh
COMPLETION_DIR="$HOME/.oh-my-zsh/completions"

# Tạo thư mục nếu nó chưa tồn tại
mkdir -p "$COMPLETION_DIR"

# Sao chép tệp _vmctl vào thư mục hoàn thành
cp ./_vmctl "$COMPLETION_DIR/_vmctl"
echo "Copied _vmctl to $COMPLETION_DIR/_vmctl"

# Tải lại cấu hình Oh My Zsh
source "$HOME/.zshrc"
exec zsh