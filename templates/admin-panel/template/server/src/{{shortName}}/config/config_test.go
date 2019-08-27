package config

import (
	stdOS "os"
	"path"
	"strings"
	"test_helpers"
	"testing"
)

type mockedFileInfo struct {
	stdOS.FileInfo
}

type mockedOS struct {
	reportErr bool
	calls     *test_helpers.FunctionsCallsStorage
}

func (mockedOS) IsNotExist(err error) bool { return stdOS.IsNotExist(err) }

func (m mockedOS) Stat(name string) (stdOS.FileInfo, error) {
	m.calls.For("Stat").Add(name)
	if m.reportErr {
		return nil, stdOS.ErrNotExist
	}
	return mockedFileInfo{}, nil
}

type mockedIOutil struct {
	calls *test_helpers.FunctionsCallsStorage
}

func (m mockedIOutil) WriteFile(filename string, data []byte, perm stdOS.FileMode) error {
	m.calls.For("WriteFile").Add(filename, data, perm)
	return nil
}

func TestCreateConfigIfNotExists(t *testing.T) {
	test_helpers.DisableLogger()
	oldOs := os
	mos := &mockedOS{
		reportErr: true,
		calls:     test_helpers.NewMockFunctionCalls(),
	}
	os = mos
	oldIoUtil := ioutil
	moutil := &mockedIOutil{
		calls: test_helpers.NewMockFunctionCalls(),
	}
	ioutil = moutil
	configFilename := path.Join("testdata", "config.cfg")
	defer func() {
		os = oldOs
		ioutil = oldIoUtil
	}()

	LoadConfiguration(configFilename)

	test_helpers.Equals(t, 1, mos.calls.For("Stat").Count())
	test_helpers.Equals(t, 1, moutil.calls.For("WriteFile").Count())
	test_helpers.Equals(t, configFilename, moutil.calls.For("WriteFile").ArgsAt(0)[0])
	test_helpers.Equals(t, []byte(strings.TrimSpace(defaultConfig)), moutil.calls.For("WriteFile").ArgsAt(0)[1])
}

func TestNotCreateConfigIfExists(t *testing.T) {
	test_helpers.DisableLogger()
	oldOs := os
	mos := &mockedOS{
		calls: test_helpers.NewMockFunctionCalls(),
	}
	os = mos
	oldIoUtil := ioutil
	moutil := &mockedIOutil{
		calls: test_helpers.NewMockFunctionCalls(),
	}
	ioutil = moutil
	configFilename := path.Join("testdata", "config.cfg")
	defer func() {
		os = oldOs
		ioutil = oldIoUtil
	}()

	LoadConfiguration(configFilename)

	test_helpers.Equals(t, 1, mos.calls.For("Stat").Count())
	test_helpers.Equals(t, 0, moutil.calls.For("WriteFile").Count())
}

func TestReadConfig(t *testing.T) {
	test_helpers.DisableLogger()
	configFilename := path.Join("testdata", "config.cfg")
	LoadConfiguration(configFilename)

	test_helpers.Equals(t, server{Port: 3001, AllowedDomains: []string{"http://127.0.0.1:3000", "http://localhost:3000"}}, Server)
	test_helpers.Equals(t, dataBaseConfig{Host: "127.0.0.1", Port: 5432, Name: "dbname", User: "dbuser", Pass: "dbuserpassword"}, DataBase)
	test_helpers.Equals(t, session{Secret: "secret-key"}, Session)
	test_helpers.Equals(t, "postgres://dbuser:dbuserpassword@127.0.0.1:5432/dbname", DataBase.GetURL())
}
