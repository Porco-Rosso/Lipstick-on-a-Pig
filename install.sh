#! /bin/bash

# Check if interactive and warn user of destruction
case $- in
*i*)
	printf "\033[1;31mThis script will permanently delete your current MOTD. I am not responsible for any potential harm.\033[0m\n"
	printf "\033[1;31mAre you sure you want to proceed?\033[0m\n"
	read -p -r "Continue (y/n)?" CONT
	;;
*) CONT="y" ;;
esac
printf "\033[1;31mThis script will permanently delete your current MOTD. I am not responsible for any potential harm.\033[0m\n"
printf "\033[1;31mAre you sure you want to proceed?\033[0m\n"
read -p -r "Continue (y/n)?" CONT
if [ "$CONT" = "y" ]; then

	printf "\033[1;31mGrabbing latest binary\033[0m\n"
	#Get OS
	platform='unknown'
	unamestr=$(uname)
	if [[ "$unamestr" == 'Linux' ]]; then
		platform='linux'
	elif [[ "$unamestr" == 'Darwin' ]]; then
		platform='darwin'
	fi

	#Get Arch
	architecture=""
	case $(uname -m) in
	x86_64) architecture="amd64" ;;
	arm) dpkg --print-architecture | grep -q "arm64" && architecture="arm64" || architecture="arm" ;;
	esac

	if [[ $platform == 'linux' ]]; then
		if [[ $architecture == 'amd64' ]]; then
			wget https://github.com/Porco-Rosso/Lipstick-on-a-Pig/releases/latest/download/lipgloss-on-a-pig_linux_amd64
		elif [[ "$unamestr" == 'arm64' ]]; then
			wget https://github.com/Porco-Rosso/Lipstick-on-a-Pig/releases/latest/download/lipgloss-on-a-pig_linux_arm64
		else
			echo "sorry not supported arch"
		fi
	elif [[ $platform == 'darwin' ]]; then
		if [[ $architecture == 'amd64' ]]; then
			wget https://github.com/Porco-Rosso/Lipstick-on-a-Pig/releases/latest/download/lipgloss-on-a-pig_darwin_amd64
		elif [[ "$unamestr" == 'arm64' ]]; then
			wget https://github.com/Porco-Rosso/Lipstick-on-a-Pig/releases/latest/download/lipgloss-on-a-pig_darwin_arm64
		else
			echo "sorry not supported arch"
		fi
	else
		echo "sorry not supported OS"
		return
	fi

	# move and make executable
	mkdir /usr/local/bin/
	mv lipgloss-on-a-pig_* /usr/local/bin/lipgloss-on-a-pig

	printf "\033[1;31mMaking executable\033[0m\n"
	chmod +x /usr/local/bin/lipgloss-on-a-pig

	# Remove old MOTD
	printf "\033[1;31mRemoving old MOTD\033[0m\n"
	echo -n >/etc/update-motd.d/10-uname
	echo -n >/etc/motd

	# Add to .bashrc
	printf "\033[1;31mAdding lipstick-on-a-pig to MOTD\033[0m\n"
	LINE='/usr/local/bin/lipgloss-on-a-pig'
	FILE=$HOME'/.bashrc'
	grep -qF -- "$LINE" "$FILE" || echo "$LINE" >>"$FILE"

	#Run at close
	printf "\033[1;31mInstalled!\033[0m\n"
	/usr/local/bin/lipgloss-on-a-pig

else
	printf "\033[1;31mBetter safe than sorry, maybe try manually installing?\033[0m\n"
fi
