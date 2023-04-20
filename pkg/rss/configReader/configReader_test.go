package configReader

import (
	"testing"
)

func TestConfigReader_Read(t *testing.T) {
	config := Read("./../../../config.json")
	t.Log(config)
}
