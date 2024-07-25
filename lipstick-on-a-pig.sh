#! /bin/bash
##TOD Add warning of destruction

printf "\033[1;31mGrabbing latest binary\033[0m\n"
#TODO make OS dependent
wget https://github.com/Porco-Rosso/Lipstick-on-a-Pig/releases/latest/download/lipgloss-on-a-pig_darwin_amd64
mkdir /usr/local/bin/
mv lipgloss-on-a-pig_darwin_amd64 /usr/local/bin/lipgloss-on-a-pig

printf "\033[1;31mMaking executable\033[0m\n"
chmod +x /usr/local/bin/lipgloss-on-a-pig

printf "\033[1;31mRemoving old MOTD\033[0m\n"
echo -n >/etc/update-motd.d/10-uname
echo -n > /etc/motd

printf "\033[1;31mAdding lipstick-on-a-pig to MOTD\033[0m\n"
LINE='/usr/local/bin/lipgloss-on-a-pig"'
FILE='~/.bashrc'
grep -qF -- "$LINE" "$FILE" || echo "$LINE" >> "$FILE"

printf "\033[1;31mInstalled!\033[0m\n"
/usr/local/bin/lipgloss-on-a-pig