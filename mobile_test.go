package zkflowexample

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMobileFullFlow(t *testing.T) {
	m := MobileZKFlow{}
	dir, err := ioutil.TempDir("", "TestMobileFullFlow")
	require.Nil(t, err)
	res, err := m.Run(dir)
	fmt.Println(res)
	os.RemoveAll(dir)
	require.Nil(t, err)
}
