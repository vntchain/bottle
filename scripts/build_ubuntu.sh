     WASMCEPTION_URL="https://github.com/ooozws/clang-heroku-slug/raw/master/precomp/wasmception-linux-bin.tar.gz"
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


     printf "\\tInstall libclang.so in /usr/lib.\\n"
	if [ ! -f $ROOT/build/lib/clang/clang.tar.xz ]
	then
		mkdir -p $ROOT/build/lib/clang
		wget  -O $ROOT/build/lib/clang/clang.tar.xz $CLANG_URL
		cd  $ROOT/build/lib/clang
		mkdir -p clang
		tar -xvf clang.tar.xz --strip-components 1 -C ./clang
		# if ! sudo ln -s  $ROOT/build/lib/clang/clang/lib/libclang.so.5.0 /usr/lib/libclang.so
		# then
		#      printf "\\tlibclang.so has installed.\\n"
		# fi
	fi
	printf "\\tInstall libclang.dylib successfully.\\n"


     function print_instructions()
	{	 
		printf "\\tDONOT REMOVE BUILD DIRECTORY UNLESS YOU WANT TO REMOVE BOTTLE\\n"
		printf "\\tTo verify your installation run the following commands:\\n"
		printf "\\tcd %s; ./bottle --help\\n\\n" "build/bin/"
		return 0
	}