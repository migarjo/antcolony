package antcolony

type jsonable interface {
	jsonify() []byte
	writeToFile(string)
}
