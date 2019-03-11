#!/bin/sh
set -e
TIME_BEGIN=$( date -u +%s )
PWD="$(pwd)"
ROOT="$PWD"
MARCH=$(go env GOOS)
GOVERSION=$(go version)
clang_mac_url="http://releases.llvm.org/5.0.0/clang+llvm-5.0.0-x86_64-apple-darwin.tar.xz"
clang_dir="clang"
SOURCE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
WASMFILE="${SOURCE_DIR}/set_wasmception.sh"

function Build_Linux {
     echo "linux"
}

if [ "$MARCH" = "darwin" ];then
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
         echo "centos"
         FILE="${SOURCE_DIR}/build_centos.sh"
      ;;
      "Ubuntu")
         echo "ubuntu"
         FILE="${SOURCE_DIR}/build_ubuntu.sh"
      ;;
      *)
         printf "\\n\\tUnsupported Linux Distribution. Exiting now.\\n\\n"
         exit 1
   esac
fi

. "$FILE"
. "$WASMFILE"


TIME_END=$(( $(date -u +%s) - ${TIME_BEGIN} ))
printf "\n\n${bldred}\t\n"
printf '\t ____  _____  ____  ____  __    ____ \n'
printf "\t(  _ \(  _  )(_  _)(_  _)(  )  ( ___)\n"
printf "\t ) _ < )(_)(   )(    )(   )(__  )__) \n"
printf "\t(____/(_____) (__)  (__) (____)(____)\n${txtrst}"

printf "\\n\\tBottle has been successfully built. %02d:%02d:%02d\\n\\n" $(($TIME_END/3600)) $(($TIME_END%3600/60)) $(($TIME_END%60))




print_instructions



