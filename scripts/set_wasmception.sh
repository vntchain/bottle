     printf "\\tInstall wasmception.\\n"
	if [ ! -f $ROOT/lib/wasmception/wasmception.tar.xz ]
	then
		mkdir -p $ROOT/lib/wasmception
		wget  -O $ROOT/lib/wasmception/wasmception.tar.xz $WASMCEPTION_URL
		cd  $ROOT/lib/wasmception
		mkdir -p wasmception
		tar -xvf wasmception.tar.xz  -C ./wasmception
		# if ! sudo ln -s  $ROOT/lib/clang/clang/lib/libclang.dylib /usr/local/lib
		# then
		#      printf "\\tlibclang.dylib has installed.\\n"
		# fi
		echo export VNT_WASMCEPTION="$ROOT/lib/wasmception/wasmception" >> ~/.bash_profile 
	fi