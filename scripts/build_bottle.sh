     cd "$bottlePath"
     echo "111" "$bottlePath"
     if ! CGO_LDFLAGS="-L$ROOT/build/lib/clang/clang/lib -Wl,-rpath,$ROOT/build/lib/clang/clang/lib -lclang" go install -ldflags " -X main.vntIncludeFlag=$ROOT/build/lib/clang/clang/lib/clang/5.0.0/include  -X main.wasmCeptionFlag=$ROOT/build/lib/wasmception/wasmception -X main.nodeFlag=$ROOT/build/lib/node/node -X main.gitCommit=$COMMITID  -v" ../bottle
     then 
          printf "\\tError compiling bottle.\\n"
          printf "\\tExiting now.\\n\\n"
          exit 1;
     fi 

     if ! CGO_LDFLAGS="-L$ROOT/build/lib/clang/clang/lib -Wl,-rpath,$ROOT/build/lib/clang/clang/lib -lclang" go test ./... -v -count=1
     then 
          printf "\\tRun test Error.\\n"
          printf "\\tExiting now.\\n\\n"
          exit 1;
     fi 