fastfetch

if [[ -r "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh" ]]; then
  source "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh"
fi

# source plugins
source /usr/share/zsh-theme-powerlevel10k/powerlevel10k/powerlevel10k.zsh-theme

source /usr/share/zsh-autosuggestions/zsh-autosuggestions.zsh

source /usr/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh

[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh

# apparently should make stuff colourful (it doesnt)
# zstyle ':completion:*' list-colors "${(s.:.)LS_COLORS}"
autoload -Uz compinit && compinit
zstyle ':completion:*' matcher-list 'm:{a-z}={A-Za-z}'

# history
HISTSIZE=5000
HISTFILE=~/.zsh_history
SAVEHIST=$HISTSIZE
HISTDUP=erase
setopt appendhistory
setopt hist_ignore_all_dups
setopt hist_save_no_dups
setopt hist_ignore_dups
setopt hist_find_no_dups
setopt hist_ignore_space

# aliases
alias hypr="nvim .config/hypr/hyprland.conf"
alias ls="ls --color --group-directories-first"
alias cd="z"
alias zat="zathura --fork"

compdef -d pacman
compdef -d npm

eval "$(zoxide init zsh)"
