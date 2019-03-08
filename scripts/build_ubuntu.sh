     
     #WASMCEPTION_URL="https://github.com/ooozws/clang-heroku-slug/raw/master/precomp/wasmception-linux-bin.tar.gz"
     WASMCEPTION_URL="http://192.168.9.251:9000/temp/wasmception-linux-bin.tar.gz"
     OS_VERSION=$( cat /etc/os-release | grep ^VERSION_ID | cut -d'=' -f2 | sed 's/\"//gI' )
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


     printf "\\tInstall libclang.so in /usr/local/lib.\\n"
	if [ ! -d $ROOT/lib/clang ]
	then
		mkdir -p $ROOT/lib/clang
		wget  -O $ROOT/lib/clang/clang.tar.xz $CLANG_URL
		cd  $ROOT/lib/clang
		mkdir -p clang
		tar -xvf clang.tar.xz --strip-components 1 -C ./clang
		if ! sudo ln -s  $ROOT/lib/clang/clang/lib/libclang.so /usr/local/lib
		then
		     printf "\\tlibclang.so has installed.\\n"
		fi
		echo export VNT_INCLUDE="$ROOT/lib/clang/clang/lib/clang/5.0.0/include" >> ~/.bash_profile 
	fi

	cd $ROOT
	go install -ldflags -s -v ./...