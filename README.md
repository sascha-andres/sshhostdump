# SSH menu

parses SSH config files and prints hosts to commandline

## Usage

    ssh (ag --files-with-matches --hidden --unrestricted Host ~/.ssh/ | ssh-host-dump | peco)