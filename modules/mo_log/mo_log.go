package mo_log

import (
    "github.com/jcelliott/lumber"
    "etcd_dashboard/modules/mo_conf"
    "etcd_dashboard/env"
)


type LoggerInterface interface{
    AddLoggers(...lumber.Logger)
    Trace(string, ...interface{})
    Debug(string, ...interface{})
    Info(string, ...interface{})
    Warn(string, ...interface{})
    Error(string, ...interface{})
    Fatal(string, ...interface{})
}

var Logger LoggerInterface


func NewLogger() LoggerInterface {
    logFile, err := lumber.NewFileLogger(
        mo_conf.Conf.Log.File,
        lumber.INFO,
        lumber.ROTATE,
        mo_conf.Conf.Log.Max_line,
        mo_conf.Conf.Log.Backups,
        256,
    )

    if err != nil {
        panic(err)
    }

    var log LoggerInterface
    if mo_conf.Conf.Env.Mode == env.ENV_MODE_PROD {
        log = lumber.NewMultiLogger()
        log.AddLoggers(logFile)
    } else {
        log = &DebugLogger{
            multiLogger: lumber.NewMultiLogger(),
        }
        logConsole := lumber.NewConsoleLogger(lumber.DEBUG)
        log.AddLoggers(logFile, logConsole)
    }
    return log
}



func init() {
    Logger = NewLogger()
}

