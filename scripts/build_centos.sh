     WASMCEPTION_URL="https://github.com/ooozws/clang-heroku-slug/raw/master/precomp/wasmception-linux-bin.tar.gz"
     LLVM_URL="http://releases.llvm.org/5.0.0/llvm-5.0.0.src.tar.xz"
     CLANG_URL="http://releases.llvm.org/5.0.0/cfe-5.0.0.src.tar.xz" 
     CPU_CORE=$( lscpu -pCPU | grep -v "#" | wc -l )
     printf "\\tChecking cmake installation\\n"
	if ! CMAKE=$( command -v cmake)
	then
		printf "\\n\\tCmake must be installed in order to proceed.\\n\\n"
		printf "\\tExiting now.\\n"
		exit 1
	fi

     printf "\\tChecking cmake version\\n"
     function version_gt() { test "$(printf '%s\n' "$@" | sort -V | head -n 1)" != "$1"; }
     CMAKE_VERSION=$( cmake --version |grep  "cmake versio" | cut -d' ' -f3 | cut -d'-' -f1 )
     REQUIRE_VERSION="3.4.3"
     if version_gt $REQUIRE_VERSION $CMAKE_VERSION ; then
          printf "\\n\\tCmake $REQUIRE_VERSION is the minimum required..\\n\\n"
          printf "\\n\\tPlease update cmake..\\n\\n"
		printf "\\tExiting now.\\n"
		exit 1
     fi

     
     printf "\\tBuild llvm + clang\\n"

	if [ ! -d $ROOT/build/lib/llvm/clang ]
	then
		mkdir -p $ROOT/build/lib/llvm
		wget  -O $ROOT/build/lib/llvm/clang.tar.xz $CLANG_URL
		cd  $ROOT/build/lib/llvm
		mkdir -p clang
		tar -xf clang.tar.xz --strip-components 1 -C ./clang
          rm clang.tar.xz
	fi

     if [ ! -f $ROOT/build/lib/llvm/llvm.tar.xz ]
	then
		mkdir -p $ROOT/build/lib/llvm
		wget  -O $ROOT/build/lib/llvm/llvm.tar.xz $LLVM_URL
		cd  $ROOT/build/lib/llvm
		mkdir -p llvm
		tar -xf llvm.tar.xz --strip-components 1 -C ./llvm
          rm llvm.tar.xz
          mv  $ROOT/build/lib/llvm/clang $ROOT/build/lib/llvm/llvm/tools
	fi

     mkdir -p  $ROOT/build/lib/clang/clang
     cd $ROOT/build/lib/clang/clang
     if ! cmake -DCMAKE_BUILD_TYPE=Release -G "Unix Makefiles" ../../llvm/llvm
     then
          printf "\\n\\tUnable to cmake llvm..\\n\\n"
		printf "\\tExiting now.\\n"
		exit 1
     fi


     if ! make -j"${CPU_CORE}"
     then
          printf "\\tError compiling llvm.\\n"
          printf "\\tExiting now.\\n\\n"
          exit 1;
     fi
     printf "\\tBuild llvm + clang successfully.\\n"

 
