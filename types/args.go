package types

import "errors"

type Args struct {
	TransportType string
	Line          string
	Stop          string
	Help          bool
}

func (args *Args) GetArgs(cmdArgs []string) (err error) {
	if len(cmdArgs) == 1 || cmdArgs[2] == "-h" || cmdArgs[2] == "--help" {
		args.Help = true
		return
	}
	if len(cmdArgs) < 4 {
		err = errors.New("not enough argument provided, please refer to the help: nextbus -h")
		return
	}
	if len(cmdArgs) > 4 {
		err = errors.New("too many argument provided, please refer to the help: nextbus -h")
		return
	}

	args.TransportType = cmdArgs[1]
	args.Line = cmdArgs[2]
	args.Stop = cmdArgs[3]
	return
}