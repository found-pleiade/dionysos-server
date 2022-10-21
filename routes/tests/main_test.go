package routes_test

import (
	"os"
	"testing"

	utils "github.com/Brawdunoir/dionysos-server/utils/tests"
)

func TestMain(m *testing.M) {
	utils.SetupTestEnvironment()

	exitVal := m.Run()

	os.Exit(exitVal)
}
