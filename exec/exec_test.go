package exec

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("Running tests...")
	result := m.Run()
	fmt.Println("Finished running tests.")
	os.Exit(result)
}

type testReader struct{}

func (r *testReader) Read(_p []byte) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func TestRsync(t *testing.T) {
	var in = testReader{}

	t.Run("Rsync local to local fail", func(t *testing.T) {
		srcM := &testExecutionContext{
			host: "localhost",
		}
		dstM := &testExecutionContext{
			host: "localhost",
		}

		var buffer bytes.Buffer
		var log bytes.Buffer
		io := NewCommandInOut(&buffer, &buffer, &log, &in)

		err := Rsync(io,
			srcM, "/workspace/repo1", "build/chart",
			dstM, "/root/chart",
			[]string{})

		if err != nil {
			t.Error(err)
		}

		result := buffer.String()

		if result != "Skipping, source and destination are the same: localhost\n" {
			t.Fatalf("not expected: [%s]", result)
		}
	})

	t.Run("Rsync remote to local fail", func(t *testing.T) {
		sm := &testExecutionContext{
			host: "remote",
		}
		dm := &testExecutionContext{
			host: "localhost",
		}

		var buffer bytes.Buffer
		var log bytes.Buffer
		io := NewCommandInOut(&buffer, &buffer, &log, &in)

		err := Rsync(io,
			sm, "/workspace/repo1", "build/chart",
			dm, "/root/chart", []string{})

		if err == nil {
			t.Fatalf("expected error")
		}

		if err.Error() != "remote machine cannot be localhost" {
			t.Fatalf("not expected: [%s]", err.Error())
		}
	})

	t.Run("Rsync local to remote success", func(t *testing.T) {
		sm := &testExecutionContext{
			user: "localuser",
			host: "localhost",
		}
		dm := &testExecutionContext{
			user: "remoteuser",
			host: "remotehost",
		}

		var buffer bytes.Buffer
		var log bytes.Buffer
		io := NewCommandInOut(&buffer, &buffer, &log, &in)

		err := Rsync(io,
			sm, "/workspace/repo1", "build/chart",
			dm, "/root/chart",
			[]string{})

		if err != nil {
			t.Error(err)
		}

		result := buffer.String()
		expected := "rsync --archive --relative --delete --verbose -o /workspace/repo1/./build/chart remoteuser@remotehost:/root/chart"

		if result != expected {
			t.Fatalf("\nexpected:\n[%s]\ngot:\n[%s]\n", expected, result)
		}
	})

	t.Run("Rsync remote to remote success", func(t *testing.T) {
		sm := &testExecutionContext{
			user: "srcuser",
			host: "srchost",
		}
		dm := &testExecutionContext{
			user: "dstuser",
			host: "dsthost",
		}

		var buffer bytes.Buffer
		var log bytes.Buffer
		io := NewCommandInOut(&buffer, &buffer, &log, &in)

		err := Rsync(io,
			sm, "/workspace/repo1", "build/chart",
			dm, "/root/chart",
			[]string{})

		if err != nil {
			t.Error(err)
		}

		result := buffer.String()
		expected := "ssh srcuser@srchost -- rsync --archive --relative --delete --verbose -o /workspace/repo1/./build/chart dstuser@dsthost:/root/chart"

		if result != expected {
			t.Fatalf("\nexpected:\n[%s]\ngot:\n[%s]\n", expected, result)
		}
	})

	t.Run("Rsync with excluded", func(t *testing.T) {
		sm := &testExecutionContext{
			user: "localuser",
			host: "localhost",
		}
		dm := &testExecutionContext{
			user: "remoteuser",
			host: "remotehost",
		}

		var buffer bytes.Buffer
		var log bytes.Buffer
		io := NewCommandInOut(&buffer, &buffer, &log, &in)

		err := Rsync(io,
			sm, "/workspace/repo1", "build/chart",
			dm, "/root/chart",
			[]string{".git", ".gradle"})

		if err != nil {
			t.Error(err)
		}

		result := buffer.String()
		expected := "rsync --archive --relative --delete --verbose -o --exclude=.git --exclude=.gradle /workspace/repo1/./build/chart remoteuser@remotehost:/root/chart"

		if result != expected {
			t.Fatalf("\nexpected:\n[%s]\ngot:\n[%s]\n", expected, result)
		}
	})
}

func TestScp(t *testing.T) {
	var in = testReader{}

	t.Run("Scp local to local fail", func(t *testing.T) {
		srcM := &testExecutionContext{
			host: "localhost",
		}
		dstM := &testExecutionContext{
			host: "localhost",
		}

		var buffer bytes.Buffer
		var log bytes.Buffer
		io := NewCommandInOut(&buffer, &buffer, &log, &in)

		err := Scp(io,
			srcM, "/workspace/repo1/build/chart",
			dstM, "/root/chart")

		if err != nil {
			t.Error(err)
		}

		result := buffer.String()

		if result != "Skipping, source and destination are the same: localhost\n" {
			t.Fatalf("not expected: [%s]", result)
		}
	})

	t.Run("Scp local to remote success", func(t *testing.T) {
		sm := &testExecutionContext{
			user: "localuser",
			host: "localhost",
		}
		dm := &testExecutionContext{
			user: "remoteuser",
			host: "remotehost",
		}

		var buffer bytes.Buffer
		var log bytes.Buffer
		io := NewCommandInOut(&buffer, &buffer, &log, &in)

		err := Scp(io,
			sm, "/workspace/repo1/build/chart",
			dm, "/root/chart")

		if err != nil {
			t.Error(err)
		}

		result := buffer.String()
		expected := "scp /workspace/repo1/build/chart remoteuser@remotehost:/root/chart"

		if result != expected {
			t.Fatalf("\nexpected:\n[%s]\ngot:\n[%s]\n", expected, result)
		}
	})

	t.Run("Scp remote to local success", func(t *testing.T) {
		sm := &testExecutionContext{
			user: "remoteuser",
			host: "remotehost",
		}
		dm := &testExecutionContext{
			user: "localuser",
			host: "localhost",
		}

		var buffer bytes.Buffer
		var log bytes.Buffer
		io := NewCommandInOut(&buffer, &buffer, &log, &in)

		err := Scp(io,
			sm, "/workspace/repo1/build/chart",
			dm, "/root/chart")

		if err != nil {
			t.Error(err)
		}

		result := buffer.String()
		expected := "scp remoteuser@remotehost:/workspace/repo1/build/chart /root/chart"

		if result != expected {
			t.Fatalf("\nexpected:\n[%s]\ngot:\n[%s]\n", expected, result)
		}
	})

	t.Run("Scp remote to remote success", func(t *testing.T) {
		sm := &testExecutionContext{
			user: "srcuser",
			host: "srchost",
		}
		dm := &testExecutionContext{
			user: "dstuser",
			host: "dsthost",
		}

		var buffer bytes.Buffer
		var log bytes.Buffer
		io := NewCommandInOut(&buffer, &buffer, &log, &in)

		err := Scp(io,
			sm, "/workspace/repo1/build/chart",
			dm, "/root/chart")

		if err != nil {
			t.Error(err)
		}

		result := buffer.String()
		expected := "ssh srcuser@srchost -- scp /workspace/repo1/build/chart dstuser@dsthost:/root/chart"

		if result != expected {
			t.Fatalf("\nexpected:\n[%s]\ngot:\n[%s]\n", expected, result)
		}
	})
}
