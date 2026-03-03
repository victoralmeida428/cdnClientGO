package cripto

type ICriptografia interface {
	Encode(string) (string, error)
	Decode(string) (string, error)
}
