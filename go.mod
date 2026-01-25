module github.com/miu200521358/mu_model_viewer

go 1.25.6

replace github.com/miu200521358/walk => ../walk

replace github.com/miu200521358/win => ../win

replace github.com/miu200521358/dds => ../dds

replace github.com/miu200521358/mlib_go => ../mlib_go_t4

require (
	github.com/ftrvxmtrx/tga v0.0.0-20150524081124-bd8e8d5be13a
	github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20240118000515-a250818d05e3
	github.com/miu200521358/dds v0.0.1
	github.com/miu200521358/walk v0.0.6
	github.com/miu200521358/win v0.0.2
	github.com/tiendc/go-deepcopy v1.7.2
	golang.org/x/sys v0.39.0
	golang.org/x/text v0.33.0
	gonum.org/v1/gonum v0.16.0
)

require (
	github.com/felixge/fgprof v0.9.3 // indirect
	github.com/go-gl/mathgl v1.2.0
	github.com/google/pprof v0.0.0-20211214055906-6f57359322fd // indirect
	github.com/pkg/profile v1.7.0
	golang.org/x/image v0.35.0
	golang.org/x/tools v0.40.0 // indirect
	gopkg.in/Knetic/govaluate.v3 v3.0.0 // indirect
)
