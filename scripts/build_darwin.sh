     CLANG_URL="http://releases.llvm.org/5.0.0/clang+llvm-5.0.0-x86_64-apple-darwin.tar.xz"
	WASMCEPTION_URL="https://github.com/ooozws/clang-heroku-slug/raw/master/precomp/wasmception-darwin-bin.tar.gz"
	if ! XCODESELECT=$( command -v xcode-select)
	then
		printf "\\n\\tXCode must be installed in order to proceed.\\n\\n"
		printf "\\tExiting now.\\n"
		exit 1
	fi

	printf "\\tInstall libclang.dylib in /usr/local/lib.\\n"
	if [ ! -f $ROOT/lib/clang/clang.tar.xz ]
	then
		mkdir -p $ROOT/lib/clang
		wget  -O $ROOT/lib/clang/clang.tar.xz $CLANG_URL
		cd  $ROOT/lib/clang
		mkdir -p clang
		tar -xvf clang.tar.xz --strip-components 1 -C ./clang
		if ! sudo ln -s  $ROOT/lib/clang/clang/lib/libclang.dylib /usr/local/lib
		then
		     printf "\\tlibclang.dylib has installed.\\n"
		fi
	fi

	cd $ROOT
	if ! go install -ldflags -s -v ./...
     then 
          printf "\\tError compiling bottle.\\n"
          printf "\\tExiting now.\\n\\n"
          exit 1;
     fi 

	function print_instructions()
	{	 
		printf "\\tAdd\\n"
		printf "\\texport VNT_INCLUDE=\"$ROOT/lib/clang/clang/lib/clang/5.0.0/include\"\\n" 
		printf "\\tto .bash_profile or another initialization script for your terminal and restart your terminal\\n"
		printf "\\tTo verify your installation run the following commands:\\n"
		printf "\\tcd %s; ./bottle --help\\n\\n" "build/bin/"
		return 0
	}

	

