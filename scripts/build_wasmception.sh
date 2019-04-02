     printf "\\tInstall wasmception.\\n"
	if [ ! -f $ROOT/build/lib/wasmception/wasmception.tar.xz ]
	then
		mkdir -p $ROOT/build/lib/wasmception
		wget  -O $ROOT/build/lib/wasmception/wasmception.tar.xz $WASMCEPTION_URL
		cd  $ROOT/build/lib/wasmception
		mkdir -p wasmception
		tar -xf wasmception.tar.xz  -C ./wasmception
	fi
	printf "\\tInstall wasmception successfully.\\n"