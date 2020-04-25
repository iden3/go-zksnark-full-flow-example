package zkflowexample

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlowLocal(t *testing.T) {
	// circuit 1
	idStateInputs, err := IdStateInputs()
	require.Nil(t, err)
	_, err = ExecuteFlow("testdata/circuit1", idStateInputs)
	require.Nil(t, err)

	// circuit 2
	c2inputs, err := GenInputs1()
	require.Nil(t, err)
	_, err = ExecuteFlow("testdata/circuit2", c2inputs)
	require.Nil(t, err)
}

func TestFlowDownloadingFiles(t *testing.T) {
	m := &MobileWrapper{}

	// circuit 1
	filesServer := "http://161.35.72.58:9000/circuit1"
	dir, err := ioutil.TempDir("", "TestMobileDownloadedDataCircuit1")
	require.Nil(t, err)
	idStateInputs, err := IdStateInputs()
	require.Nil(t, err)
	_, err = m.ExecuteFlowDownloading(dir, filesServer, idStateInputs)
	os.RemoveAll(dir)
	require.Nil(t, err)

	// circuit2
	filesServer = "http://161.35.72.58:9000/circuit2"
	dir, err = ioutil.TempDir("", "TestMobileDownloadedDataCircuit2")
	require.Nil(t, err)
	c2inputs, err := GenInputs1()
	require.Nil(t, err)
	_, err = m.ExecuteFlowDownloading(dir, filesServer, c2inputs)
	os.RemoveAll(dir)
	require.Nil(t, err)
}
