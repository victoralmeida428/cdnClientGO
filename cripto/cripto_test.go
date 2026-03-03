package cripto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCriptografia(t *testing.T) {

	crip := New("./public_key.pem", "./private_key.pem")

	msg := "1792"
	msgCripto, _ := crip.Encode(msg)
	msgDecode, _ := crip.Decode(msgCripto)
	fmt.Printf("Original: %s\nCriptografado: %s\nDecode: %s\n", msg, msgCripto, msgDecode)
	assert.Equal(t, msg, msgDecode)
}
