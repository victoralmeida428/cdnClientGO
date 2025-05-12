package cripto

type ICriptografia interface {
	Encode(string) string
	Decode(string) string
}
