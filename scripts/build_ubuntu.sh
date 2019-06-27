     WASMCEPTION_URL="https://github.com/ooozws/clang-heroku-slug/raw/master/precomp/wasmception-linux-bin.tar.gz"
     OS_VERSION=$( cat /etc/os-release | grep ^VERSION_ID | cut -d'=' -f2 | sed 's/\"//gI' )
     NODE_URL="https://nodejs.org/download/release/v10.16.0/node-v10.16.0-linux-x64.tar.xz"
     case "$OS_VERSION" in
          "16.04")
          CLANG_URL="http://releases.llvm.org/5.0.0/clang+llvm-5.0.0-linux-x86_64-ubuntu16.04.tar.xz"
          ;;
          "14.04")
          CLANG_URL="http://releases.llvm.org/5.0.0/clang+llvm-5.0.0-linux-x86_64-ubuntu14.04.tar.xz"
          ;;
          *)
          printf "\\n\\tUnsupported Linux Distribution. Exiting now.\\n\\n"
          exit 1
     esac


     printf "\\tInstall libclang in /usr/lib.\\n"
	if [ ! -d $ROOT/build/lib/clang/clang ]
	then
		mkdir -p $ROOT/build/lib/clang
		wget  -O $ROOT/build/lib/clang/clang.tar.xz $CLANG_URL
		cd  $ROOT/build/lib/clang
		mkdir -p clang
		tar -xf clang.tar.xz --strip-components 1 -C ./clang
          rm clang.tar.xz
	fi
	printf "\\tInstall libclang successfully.\\n"
