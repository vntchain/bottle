#!/bin/bash
set -e
TIME_BEGIN=$( date -u +%s )
ROOT="$bottlePath"
MARCH=$(go env GOOS)
GOVERSION=$(go version)
clang_dir="clang"
SOURCE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
WASMFILE="${SOURCE_DIR}/build_wasmception.sh"
BOTTLEFILE="${SOURCE_DIR}/build_bottle.sh"
COMMITID=$(git rev-parse HEAD)

txtbld=$(tput bold)
bldred=${txtbld}$(tput setaf 1)
txtrst=$(tput sgr0)

if ! GO=$( command -v go)
then
     printf "\\n\\tGo must be installed in order to proceed.\\n\\n"
     printf "\\tExiting now.\\n"
     exit 1
fi


if [ "$MARCH" = "darwin" ];then
   echo "build bottle in darwin"
   FILE="${SOURCE_DIR}/build_darwin.sh"
elif [ "$MARCH" = "linux" ];then
   if [ ! -e /etc/os-release ]; then
      printf "\\n\\bottle currently supports Centos & Ubuntu Linux only.\\n"
      printf "\\tPlease install on the latest version of one of these Linux distributions.\\n"
      printf "\\thttps://www.centos.org/\\n"
      printf "\\thttps://www.ubuntu.com/\\n"
      printf "\\tExiting now.\\n"
      exit 1
   fi
   OS_NAME=$( cat /etc/os-release | grep ^NAME | cut -d'=' -f2 | sed 's/\"//gI' )
   case "$OS_NAME" in
      "CentOS Linux")
         echo "build bottle in centos"
         FILE="${SOURCE_DIR}/build_centos.sh"
      ;;
      "Ubuntu")
         echo "build bottle in ubuntu"
         FILE="${SOURCE_DIR}/build_ubuntu.sh"
      ;;
      *)
         printf "\\n\\tUnsupported Linux Distribution. Exiting now.\\n\\n"
         exit 1
   esac
fi

. "$FILE"
. "$WASMFILE"
. "$BOTTLEFILE"


TIME_END=$(( $(date -u +%s) - ${TIME_BEGIN} ))
printf "\n\n${bldred}\t\n"
printf '\t ____  _____  ____  ____  __    ____ \n'
printf "\t(  _ \(  _  )(_  _)(_  _)(  )  ( ___)\n"
printf "\t ) _ < )(_)(   )(    )(   )(__  )__) \n"
printf "\t(____/(_____) (__)  (__) (____)(____)\n${txtrst}"

printf "\\n\\tBottle has been successfully built. %02d:%02d:%02d\\n\\n" $(($TIME_END/3600)) $(($TIME_END%3600/60)) $(($TIME_END%60))

function print_instructions()
{	 
   printf "\\t************************************************************************\\n" 
   printf "\\t*****DONOT REMOVE **BUILD** FOLDER UNLESS YOU WANT TO DELETE BOTTLE*****\\n" 
   printf "\\t************************************************************************\\n"
   printf "\\n"  
   printf "\\tTo verify your installation run the following commands:\\n"
   printf "\\tcd %s; ./bottle --help\\n\\n" "build/bin/"
   return 0
}

print_instructions




