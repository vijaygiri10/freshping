package util

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func GetExePath() string {

	var err error
	ExePath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error GetExePath : ", err)
	}
	return ExePath
}
func FuncName() string {
	defer func() {
		if errD := recover(); errD != nil {
			fmt.Println("Exception Occurred and Recovered in FuncName(), Error Info: ", errD)
		}
	}()

	pc, _, _, _ := runtime.Caller(1)
	funcName := strings.TrimSuffix(runtime.FuncForPC(pc).Name(), ".func1") // This is for defer function
	funcName = strings.TrimSuffix(funcName, ".1")                          // This is for go runtine function
	return funcName
}

func RecoverExceptionDetails(strfuncName string) string {
	defer func() {
		if errD := recover(); errD != nil {
			fmt.Println("Exception Occurred and Recovered in RecoverExceptionDetails(), Error Info: ", errD)
		}
	}()

	var output string
	flag := false
	for skip := 1; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		strfunctionName := runtime.FuncForPC(pc).Name()
		if strings.Contains(file, "/runtime/") && strings.Contains(strfunctionName, "runtime.") {
			flag = true
			continue
		}
		if flag && strings.HasSuffix(file, ".go") {
			output += strfunctionName + ":" + strconv.Itoa(line) + " << "
			if strfuncName == strfunctionName {
				output = strings.TrimSuffix(output, " << ")
				break
			}
		}
	}
	return output
}
