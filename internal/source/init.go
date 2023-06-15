package source

func init() {
	Register("stdin", NewStdin)
	Register("file", NewFile)
	Register("dir", NewDirectory)
}
