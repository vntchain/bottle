     #CLANG_URL="http://releases.llvm.org/5.0.0/clang+llvm-5.0.0-x86_64-apple-darwin.tar.xz"
	CLANG_URL="http://192.168.9.251:9000/temp/clang.tar.xz"
	#WASMCEPTION_URL="https://github.com/ooozws/clang-heroku-slug/raw/master/precomp/wasmception-darwin-bin.tar.gz"
	WASMCEPTION_URL="http://192.168.9.251:9000/temp/wasmception-darwin-bin.tar.gz"
	
	printf "\\tChecking xcode-select installation.\\n"
	if ! XCODESELECT=$( command -v xcode-select)
	then
		printf "\\n\\tXCode must be installed in order to proceed.\\n\\n"
		printf "\\tExiting now.\\n"
		exit 1
	fi

	printf "\\tInstall libclang.dylib in /usr/local/lib.\\n"
	if [ ! -d $ROOT/lib/clang ]
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
		echo export VNT_INCLUDE="$ROOT/lib/clang/clang/lib/clang/5.0.0/include" >> ~/.bash_profile 
	fi
	# if ! sudo ln -s  $PWD/lib/clang/clang/lib/libclang.dylib /usr/local/lib
	# then
	# 	printf "\\tlibclang.dylib has installed.\\n"
	# fi
	#  export CGO_LDFLAGS="-L/usr/lib/ -L/usr/local/lib/ -L$PWD/lib/clang/clang/lib/ -lclang -lstdc++"
     #  export CGO_CFLAGS="-I$PWD/lib/clang/clang/lib/clang/5.0.0/include/ -I$PWD/lib/clang/clang/include/c++/v1"
	cd $ROOT
	go install -ldflags -s -v ./...

	

