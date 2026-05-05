#!/bin/bash
set -euo pipefail

cd "${HOME}"

git clone --bare 'https://github.com/nanvenomous/unix.git' "${HOME}/.unx"

git --git-dir="${HOME}/.unx/" --work-tree="${HOME}" checkout

git --git-dir="${HOME}/.unx/" --work-tree="${HOME}" config --local status.showUntrackedFiles no

echo "dotfiles ready"
