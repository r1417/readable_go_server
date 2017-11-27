package utils

import (
    "fmt"
    "log"
    "os"
    "io"
    "encoding/json"
)

/*
 * Use only fatal error. write log and call os.exit(1)
 */ 
func OnErrorTerminate(err error, s string) {
    if err != nil {
        log.Print(fmt.Sprintf("fatal: [%s] %s", s,  err))
        panic(err)
    }
}

func OnErrorResume(err error, s string) bool {
    if err != nil {
        log.Print(fmt.Sprintf("error: [%s] %s", s,  err))
        return true
    }
    return false
}

func Info(message string, location string) {
    log.Print(fmt.Sprintf("info: [%s] %s", location, message))
}

func InfoJson(v interface{}, location string) {
    b, _ := json.Marshal(v)

    log.Print(fmt.Sprintf("info: [%s] %s", location, string(b)))
}

func InitLogMultiWriter(filePath string) *os.File {

    //CreateFile
    OnErrorResume(os.Remove(filePath), "InitLogMultiWriter:Remove")

    logfile, err := os.OpenFile(filePath, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)

    OnErrorTerminate(err, "InitLogMultiWriter:Open")

    //MultiWriter
    log.SetOutput(io.MultiWriter(logfile, os.Stdout))
    log.SetFlags(log.Ldate | log.Ltime)
    
    return logfile
}