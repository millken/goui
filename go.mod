module github.com/millken/goui

go 1.24.8

replace github.com/millken/go/rgfw => /workspace/Codes/github.com/millken/go/rgfw

replace github.com/millken/go/thorvg => /workspace/Codes/github.com/millken/go/thorvg

require (
	github.com/dnsoa/go/assert v1.1.2
	github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20250301202403-da16c1255728
	github.com/millken/go/plutovg v1.0.0
	github.com/millken/go/rgfw v0.0.0-20251027110334-6b451232f866
	github.com/millken/go/thorvg v0.0.0-00010101000000-000000000000
	github.com/millken/yoga v0.0.0-20251101114009-381a31894288
	golang.org/x/exp v0.0.0-20251023183803-a4bb9ffd2546
	golang.org/x/image v0.32.0
)

require (
	github.com/ebitengine/purego v0.9.0 // indirect
	golang.org/x/text v0.30.0 // indirect
)
