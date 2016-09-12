package mo_conf

import (
    "os"
    "etcd_dashboard/env"
    "github.com/BurntSushi/toml"
)

type EnvConf struct {
    Mode string
}

type ServerConf struct {
    Host string
    Port int
}

type LogConf struct {
    File string
    Max_line int
    Backups int
}

type EtcdConf struct {
    Name string
    Addr string
}

type MainConf struct {
    Env EnvConf
    Log LogConf
    Server ServerConf
    Etcd []EtcdConf
}

var (
    Conf = &MainConf{}
)

func init() {
    confPath := os.Getenv(env.ENV_CONF_FILE_VARIABLE_NAME)
    if confPath == "" {
        panic("conf file not found, please set environment variable '" + env.ENV_CONF_FILE_VARIABLE_NAME + "'")
    }
    toml.DecodeFile(confPath, Conf)
}
