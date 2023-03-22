package demo

import (
	"testing"

	"mackey/pkg/tests"

	. "github.com/smartystreets/goconvey/convey"
)

func TestManager_Hash(t *testing.T) {
	tests.Dep("Should test hash", t, func(config TestConfig, m *Manager) {
		hash := m.Hash(config.TestHash)

		So(hash, ShouldEqual, "d41d8cd98f00b204e9800998ecf8427e")
	})
}
