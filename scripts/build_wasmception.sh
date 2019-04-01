     printf "\\tInstall wasmception.\\n"
	if [ ! -f $ROOT/build/lib/wasmception/wasmception.tar.xz ]
	then
		mkdir -p $ROOT/build/lib/wasmception
		wget  -O $ROOT/build/lib/wasmception/wasmception.tar.xz $WASMCEPTION_URL
		cd  $ROOT/build/lib/wasmception
		mkdir -p wasmception
		tar -xvf wasmception.tar.xz  -C ./wasmception
		# if ! sudo ln -s  $ROOT/build/lib/clang/clang/lib/libclang.dylib /usr/local/lib
		# then
		#      printf "\\tlibclang.dylib has installed.\\n"
		# fi
		echo export VNT_WASMCEPTION="$ROOT/build/lib/wasmception/wasmception" >> ~/.bash_profile 
	fi