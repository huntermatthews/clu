package pkg

type CommandResult struct {
	Stdout string
	RC     int
}

type MockFS struct {
	Files map[string]string
}

func NewMockFS(files map[string]string) *MockFS { return &MockFS{Files: files} }

func (m *MockFS) FileReader(path string) (string, error) {
	if m == nil || m.Files == nil {
		return "", ErrEmptyFile
	}
	if s, ok := m.Files[path]; ok {
		return s, nil
	}
	return "", ErrEmptyFile
}

type MockCommands struct {
	Outputs map[string]CommandResult
}

func NewMockCommands(outputs map[string]CommandResult) *MockCommands {
	return &MockCommands{Outputs: outputs}
}

func (m *MockCommands) CommandRunner(cmdline string) (string, int) {
	if m == nil || m.Outputs == nil {
		return "", 1
	}
	if res, ok := m.Outputs[cmdline]; ok {
		return res.Stdout, res.RC
	}
	return "", 1
}

func UseMockFS(fs *MockFS) func() {
	prev := FileReader
	FileReader = fs.FileReader
	return func() { FileReader = prev }
}

func UseMockCommands(cmds *MockCommands) func() {
	prev := CommandRunner
	CommandRunner = cmds.CommandRunner
	return func() { CommandRunner = prev }
}

func UseMockInput(fs *MockFS, cmds *MockCommands) func() {
	prevCmd := CommandRunner
	prevFile := FileReader
	if fs != nil {
		FileReader = fs.FileReader
	}
	if cmds != nil {
		CommandRunner = cmds.CommandRunner
	}
	return func() {
		CommandRunner = prevCmd
		FileReader = prevFile
	}
}
