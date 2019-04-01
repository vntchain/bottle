     CLANG_URL="http://releases.llvm.org/5.0.0/clang+llvm-5.0.0-x86_64-apple-darwin.tar.xz"
	WASMCEPTION_URL="https://github.com/ooozws/clang-heroku-slug/raw/master/precomp/wasmception-darwin-bin.tar.gz"
	if ! XCODESELECT=$( command -v xcode-select)
	then
		printf "\\n\\tXCode must be installed in order to proceed.\\n\\n"
		printf "\\tExiting now.\\n"
		exit 1
	fi

	printf "\\tInstall libclang.dylib.\\n"
	if [ ! -f $ROOT/build/lib/clang/clang.tar.xz ]
	then
		mkdir -p $ROOT/build/lib/clang
		wget  -O $ROOT/build/lib/clang/clang.tar.xz $CLANG_URL
		cd  $ROOT/build/lib/clang
		mkdir -p clang
		tar -xvf clang.tar.xz --strip-components 1 -C ./clang
		# if ! sudo ln -s  $ROOT/build/lib/clang/clang/lib/libclang.dylib /usr/local/lib
		# then
		#      printf "\\tlibclang.dylib has installed.\\n"
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

	

