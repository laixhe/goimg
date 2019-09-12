module github.com/laixhe/goimg

go 1.13

require (
	github.com/go-ini/ini v1.46.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337 // indirect
	gopkg.in/ini.v1 v1.46.0 // indirect
)

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190911031432-227b76d455e7
	golang.org/x/net => github.com/golang/net v0.0.0-20190909003024-a7b16738d86b
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190911201528-7ad0cfa0b7b5
	golang.org/x/text => github.com/golang/text v0.3.3-0.20190829152558-3d0f7978add9
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190911230505-6bfd74cf029c
	golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20190717185122-a985d3407aa7
)
