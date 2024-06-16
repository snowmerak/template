//go:build linux || darwin

package bufbuild

// when you are working on linux or macos
//go:generate docker run -rm --volume "$(pwd):/workspace" --workdir /workspace bufbuild/buf generate
