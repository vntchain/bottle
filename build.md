CGO_CFLAGS='-I/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/lib/clang/10.0.0/include/'  CGO_LDFLAGS='-L/Users/weisaizhang/Documents/go/src/github.com/vntchain/bottle/lib -lclang' go build *.go

CGO_CFLAGS='-I/usr/include'  CGO_LDFLAGS='-L/home/ubuntu/llvm/clang+llvm-5.0.1-x86_64-linux-gnu-ubuntu-16.04/lib' go install ./...

LD_LIBRARY_PATH=/home/ubuntu/llvm/clang+llvm-5.0.1-x86_64-linux-gnu-ubuntu-16.04/lib ./analyse

LD_LIBRARY_PATH=/home/ubuntu/llvm/clang+llvm-5.0.1-x86_64-linux-gnu-ubuntu-16.04/lib CGO_CFLAGS='-I/usr/include'  CGO_LDFLAGS='-L/home/ubuntu/llvm/clang+llvm-5.0.1-x86_64-linux-gnu-ubuntu-16.04/lib' go run *.go

${SRCDIR}