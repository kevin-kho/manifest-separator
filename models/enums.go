package models

type Mode int

const (
	ModeDash Mode = iota
	ModeList
	ModeAppSet
)

type Cmd int

const (
	CmdGet Cmd = iota
	CmdDiff
	CmdApply
)
