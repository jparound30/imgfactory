package logger

import (
	"fmt"
	"log"
	"path"
	"runtime"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
}

func callerFunctionInfo() (string, int, string) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown", -1, "unknown"
	}
	// ポインターから関数名に変換する
	return path.Base(file), line, runtime.FuncForPC(pc).Name()
}

func Printf(format string, v ...interface{}) {
	file, line, functionName := callerFunctionInfo()
	prepend := fmt.Sprintf("%s:%d: [%v()]: ", file, line, functionName)
	var vv []interface{}
	newSlice := append(vv, prepend)
	newSlice = append(newSlice, v...)
	log.Printf("%s"+format, newSlice...)
}

func Print(v ...interface{}) {
	file, line, functionName := callerFunctionInfo()
	prepend := fmt.Sprintf("%s:%d: [%v()]: ", file, line, functionName)
	var vv []interface{}
	newSlice := append(vv, prepend)
	newSlice = append(newSlice, v...)
	log.Print(newSlice...)
}

func Println(v ...interface{}) {
	file, line, functionName := callerFunctionInfo()
	prepend := fmt.Sprintf("%s:%d: [%v()]: ", file, line, functionName)
	var vv []interface{}
	newSlice := append(vv, prepend)
	newSlice = append(newSlice, v...)
	log.Println(newSlice...)
}

func Fatal(v ...interface{}) {
	file, line, functionName := callerFunctionInfo()
	prepend := fmt.Sprintf("%s:%d: [%v()]: ", file, line, functionName)
	var vv []interface{}
	newSlice := append(vv, prepend)
	newSlice = append(newSlice, v...)
	log.Fatal(newSlice...)
}

func Panicf(format string, v ...interface{}) {
	file, line, functionName := callerFunctionInfo()
	prepend := fmt.Sprintf("%s:%d: [%v()]: ", file, line, functionName)
	var vv []interface{}
	newSlice := append(vv, prepend)
	newSlice = append(newSlice, v...)
	log.Panicf("%s"+format, newSlice...)
}
