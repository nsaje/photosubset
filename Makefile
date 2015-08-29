photosubset: photosubset.go qrc.go
	go build

qrc.go: qml/photosubset.qml
	go generate
