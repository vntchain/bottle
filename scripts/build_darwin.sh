	CLANG_URL="http://releases.llvm.org/5.0.0/clang+llvm-5.0.0-x86_64-apple-darwin.tar.xz"
	WASMCEPTION_URL="https://github.com/ooozws/clang-heroku-slug/raw/master/precomp/wasmception-darwin-bin.tar.gz"
	NODE_URL="https://nodejs.org/download/release/v10.16.0/node-v10.16.0-darwin-x64.tar.xz"
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

	printf "\\tInstall nodejs.\\n"
	if [ ! -d $ROOT/build/lib/node/node ]
	then
		mkdir -p $ROOT/build/lib/node
		wget  -O $ROOT/build/lib/node/node.tar.xz $NODE_URL
		cd  $ROOT/build/lib/node
		mkdir -p node
		tar -xf node.tar.xz --strip-components 1 -C ./node
		rm node.tar.xz
	fi
	printf "\\tInstall nodejs successfully.\\n"



	

