package main

type jsonable interface {
	jsonify() []byte
	writeToFile(string)
}
