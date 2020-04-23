package zkflowexample

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMobileFullFlow(t *testing.T) {
	m := MobileZKFlow{}
	dir, err := ioutil.TempDir("", "TestMobileFullFlow")
	require.Nil(t, err)
	err = m.Run(dir)
	os.RemoveAll(dir)
	require.Nil(t, err)
}
