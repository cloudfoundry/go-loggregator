module code.cloudfoundry.org/go-loggregator/v10/examples

go 1.21
toolchain go1.24.1

require (
	code.cloudfoundry.org/go-loggregator/v10 v10.0.0-00010101000000-000000000000
	github.com/cloudfoundry/dropsonde v1.1.0
	google.golang.org/protobuf v1.36.6
)

require (
	code.cloudfoundry.org/go-diodes v0.0.0-20180905200951-72629b5276e3 // indirect
	code.cloudfoundry.org/tlsconfig v0.24.0 // indirect
	github.com/cloudfoundry/sonde-go v0.0.0-20220627221915-ff36de9c3435 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/grpc v1.71.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)

replace code.cloudfoundry.org/go-loggregator/v10 => ../
