package cripto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCriptografia(t *testing.T) {

	crip := New()

	msg := "1792"
	msgCripto := crip.Encode(msg)
	msgDecode := crip.Decode(msgCripto)
	fmt.Printf("Original: %s\nCriptografado: %s\nDecode: %s\n", msg, msgCripto, msgDecode)
	assert.Equal(t, msg, msgDecode)
}
