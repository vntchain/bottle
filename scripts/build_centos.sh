     WASMCEPTION_URL="https://github.com/ooozws/clang-heroku-slug/raw/master/precomp/wasmception-linux-bin.tar.gz"
     LLVM_URL="http://releases.llvm.org/5.0.0/llvm-5.0.0.src.tar.xz"
     CLANG_URL="http://releases.llvm.org/5.0.0/cfe-5.0.0.src.tar.xz" 
     printf "\\tChecking cmake installation\\n"
	if ! CMAKE=$( command -v cmake)
	then
		printf "\\n\\tCmake must be installed in order to proceed.\\n\\n"
		printf "\\tExiting now.\\n"
		exit 1
	fi

     printf "\\tChecking cmake version\\n"
     function version_gt() { test "$(printf '%s\n' "$@" | sort -V | head -n 1)" != "$1"; }
     CMAKE_VERSION=$( cmake --version |grep  "cmake versio" | cut -d' ' -f3 )
     REQUIRE_VERSION="3.4.3"
     if version_gt $REQUIRE_VERSION $CMAKE_VERSION ; then
          printf "\\n\\tCmake $REQUIRE_VERSION is the minimum required..\\n\\n"
          printf "\\n\\Please update cmake..\\n\\n"
		printf "\\tExiting now.\\n"
		exit 1
     fi

     
     printf "\\tBuild llvm + clang\\n"

	if [ ! -f $ROOT/lib/llvm/clang.tar.xz ]
	then
		mkdir -p $ROOT/lib/llvm
		wget  -O $ROOT/lib/llvm/clang.tar.xz $CLANG_URL
		cd  $ROOT/lib/llvm
		mkdir -p clang
		tar -xvf clang.tar.xz --strip-components 1 -C ./clang
	fi

     if [ ! -f $ROOT/lib/llvm/llvm.tar.xz ]
	then
		mkdir -p $ROOT/lib/llvm
		wget  -O $ROOT/lib/llvm/llvm.tar.xz $LLVM_URL
		cd  $ROOT/lib/llvm
		mkdir -p llvm
		tar -xvf llvm.tar.xz --strip-components 1 -C ./llvm
          mv  $ROOT/lib/llvm/clang $ROOT/lib/llvm/llvm/tools
	fi

     cd  $ROOT/lib/llvm
     mkdir -p llvm_build
     cd llvm_build
     cmake ../llvm
