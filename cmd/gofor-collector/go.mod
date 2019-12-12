module github.com/Go-Forensics/Windows-Collector/cmd/gofor-collector

go 1.13.5

require (
	github.com/Go-Forensics/Windows-Collector v0.3.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/sirupsen/logrus v1.4.2
)

replace github.com/Go-Forensics/Windows-Collector => ../../
