//go:build windows

package bufbuild

// when you are working on windows
//go:generate docker run -rm --volume "${PWD}:/workspace" --workdir /workspace bufbuild/buf generate
