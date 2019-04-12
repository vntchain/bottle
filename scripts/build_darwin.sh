    # CLANG_URL="http://releases.llvm.org/5.0.0/clang+llvm-5.0.0-x86_64-apple-darwin.tar.xz"
	CLANG_URL="http://192.168.9.251:9000/temp/clang+llvm-5.0.0-x86_64-apple-darwin.tar.xz"
	# WASMCEPTION_URL="https://github.com/ooozws/clang-heroku-slug/raw/master/precomp/wasmception-darwin-bin.tar.gz"
	WASMCEPTION_URL="http://192.168.9.251:9000/temp/wasmception-darwin-bin.tar.gz"
	if ! XCODESELECT=$( command -v xcode-select)
	then
		printf "\\n\\tXCode must be installed in order to proceed.\\n\\n"
		printf "\\tExiting now.\\n"
		exit 1
	fi

	printf "\\tInstall libclang.\\n"
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


	

