module github.com/ftl/patrix

go 1.22

toolchain go1.23.6

//replace (
//	github.com/ftl/digimodes => ../digimodes
//  github.com/jfreymuth/pulse => ../pulse
//)

require (
	github.com/ftl/digimodes v0.1.2
	github.com/jfreymuth/pulse v0.1.0
	github.com/spf13/cobra v1.8.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)
