module demo

go 1.24.8

replace github.com/millken/goui => ../../

replace github.com/millken/go/rgfw => /workspace/Codes/github.com/millken/go/rgfw

replace github.com/millken/thorvg => /workspace/Codes/github.com/millken/thorvg

require (
	github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71
	github.com/millken/go/rgfw v0.0.0-00010101000000-000000000000
	github.com/millken/thorvg v0.0.0-00010101000000-000000000000
	github.com/millken/goui v0.0.0-00010101000000-000000000000
)

require (
	github.com/ebitengine/purego v0.8.3 // indirect
	golang.org/x/image v0.27.0 // indirect
	golang.org/x/text v0.25.0 // indirect
)
