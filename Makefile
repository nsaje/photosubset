photosubset: photosubset.go qml
	go build

qml: qml/photosubset.qml
	go generate
