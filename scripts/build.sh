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

function Build_Linux {
     echo "linux"
}

if [ "$MARCH" = "darwin" ];then
   FILE="${SOURCE_DIR}/build_darwin.sh"
elif [ "$MARCH" = "linux" ];then
   FILE="${SOURCE_DIR}/build_darwin.sh"
fi

. "$FILE"


TIME_END=$(( $(date -u +%s) - ${TIME_BEGIN} ))
printf "\n\n${bldred}\t\n"
printf '\t ____  _____  ____  ____  __    ____ \n'
printf "\t(  _ \(  _  )(_  _)(_  _)(  )  ( ___)\n"
printf "\t ) _ < )(_)(   )(    )(   )(__  )__) \n"
printf "\t(____/(_____) (__)  (__) (____)(____)\n${txtrst}"

printf "\\n\\tBottle has been successfully built. %02d:%02d:%02d\\n\\n" $(($TIME_END/3600)) $(($TIME_END%3600/60)) $(($TIME_END%60))
printf "\\tTo verify your installation run the following commands:\\n"
print_instructions



