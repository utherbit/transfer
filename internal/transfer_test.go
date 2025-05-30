package internal

import (
	"bytes"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const testdataDirectory = "./testdata"

func TestTransferWithTestData(t *testing.T) {
	walkTestdata(t, testdataDirectory)
}

func walkTestdata(t *testing.T, dir string) {
	entry, err := os.ReadDir(dir)
	require.NoError(t, err)

	for _, e := range entry {
		if !e.IsDir() {
			continue
		}

		dir := filepath.Join(dir, e.Name())

		if thisDirIsTestCase(t, dir) {
			RunTestCase(t, dir)
			continue
		}

		walkTestdata(t, dir)
	}
}

func thisDirIsTestCase(t *testing.T, dir string) bool {
	entry, err := os.ReadDir(dir)
	require.NoError(t, err)

	var (
		containsEntity         = false
		containsEntityTransfer = false
	)

	for _, e := range entry {
		if e.Name() == "entity.go" {
			containsEntity = true
		}
		if e.Name() == "entity_transfer.go" {
			containsEntityTransfer = true
		}
	}

	return containsEntity && containsEntityTransfer
}

func RunTestCase(t *testing.T, dir string) {
	transferPath, _ := url.JoinPath(dir, "entity_transfer.go")
	transferFile, err := os.Open(transferPath)
	require.NoError(t, err)

	transferExpect, err := io.ReadAll(transferFile)
	require.NoError(t, err)

	t.Run(filepath.Base(dir), func(t *testing.T) {
		parseReq, err := findStructByDirAndType(dir, "Entity")
		require.NoError(t, err)

		structInfo, err := parseStruct(parseReq)
		require.NoError(t, err)

		buf := bytes.NewBuffer(make([]byte, 0))

		err = generateTransfer(structInfo, buf)
		require.NoError(t, err)

		require.Equal(t, string(transferExpect), buf.String())
	})
}
