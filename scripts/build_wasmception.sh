     printf "\\tInstall wasmception.\\n"
	if [ ! -d $ROOT/build/lib/wasmception/wasmception ]
	then
		mkdir -p $ROOT/build/lib/wasmception
		wget  -O $ROOT/build/lib/wasmception/wasmception.tar.xz $WASMCEPTION_URL
		cd  $ROOT/build/lib/wasmception
		mkdir -p wasmception
		tar -xf wasmception.tar.xz  -C ./wasmception
		rm wasmception.tar.xz
	fi
	printf "\\tInstall wasmception successfully.\\n"