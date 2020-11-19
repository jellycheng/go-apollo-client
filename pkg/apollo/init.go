package apollo

import "runtime"

var LineBreak = "\n"

func init()  {

	if runtime.GOOS == "windows" {
		LineBreak = "\r\n"
	}

}
