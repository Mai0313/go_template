FROM mcr.microsoft.com/devcontainers/go:1-ubuntu

LABEL maintainer="Wei Lee <mai@mai0313.com>" \
    org.label-schema.name="go_template" \
    org.label-schema.vendor="Wei Lee" \
    org.label-schema.schema-version="1.0" \
    com.centurylinklabs.watchtower.stop-signal="SIGINT"

RUN git config --global http.sslVerify false && \
    apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    zsh vim wget curl fontconfig ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Install Meslo Nerd Fonts for powerlevel10k
RUN mkdir -p "/usr/local/share/fonts/meslo" && \
    wget -q -O "/usr/local/share/fonts/meslo/MesloLGS NF Regular.ttf" "https://github.com/romkatv/dotfiles-public/raw/master/.local/share/fonts/NerdFonts/MesloLGS%20NF%20Regular.ttf" && \
    wget -q -O "/usr/local/share/fonts/meslo/MesloLGS NF Bold.ttf" "https://github.com/romkatv/dotfiles-public/raw/master/.local/share/fonts/NerdFonts/MesloLGS%20NF%20Bold.ttf" && \
    wget -q -O "/usr/local/share/fonts/meslo/MesloLGS NF Italic.ttf" "https://github.com/romkatv/dotfiles-public/raw/master/.local/share/fonts/NerdFonts/MesloLGS%20NF%20Italic.ttf" && \
    wget -q -O "/usr/local/share/fonts/meslo/MesloLGS NF Bold Italic.ttf" "https://github.com/romkatv/dotfiles-public/raw/master/.local/share/fonts/NerdFonts/MesloLGS%20NF%20Bold%20Italic.ttf" && \
    fc-cache -fv

# Install oh-my-zsh non-interactively
ENV RUNZSH=no CHSH=no KEEP_ZSHRC=yes
RUN sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"

# Install powerlevel10k theme and useful plugins
RUN git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/themes/powerlevel10k && \
    git clone --depth=1 https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions && \
    git clone --depth=1 https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting && \
    git clone --depth=1 https://github.com/agkozak/zsh-z ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-z

# Configure zsh theme and plugins
RUN sed -i 's|^ZSH_THEME=.*$|ZSH_THEME="powerlevel10k/powerlevel10k"|' $HOME/.zshrc && \
    sed -i 's|^plugins=(.*)$|plugins=(git zsh-autosuggestions zsh-syntax-highlighting zsh-z)|' $HOME/.zshrc && \
    sed -i 's|^# ZSH_CUSTOM=.*$|ZSH_CUSTOM=$HOME/.oh-my-zsh/custom|' $HOME/.zshrc && \
    sed -i 's|^# export PATH=\$HOME/bin:\$HOME/.local/bin:/usr/local/bin:\$PATH$|export PATH=$HOME/bin:$HOME/.local/bin:/usr/local/bin:$PATH|' $HOME/.zshrc

# Install common Go tools (optional for local quality checks)
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install golang.org/x/tools/cmd/deadcode@latest || true
