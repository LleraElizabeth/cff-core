package main

import "C"
import (
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/hub/executor"
	"github.com/Dreamacro/clash/hub/route"
	"github.com/Dreamacro/clash/log"
	"github.com/oschwald/geoip2-golang"
	"go.uber.org/automaxprocs/maxprocs"
	"os"
	"path/filepath"
)

func init() {
	_, _ = maxprocs.Set(maxprocs.Logger(func(string, ...any) {}))
}

//export SetHomeDir
func SetHomeDir(homeStr *C.char) bool {
	homeDir := C.GoString(homeStr)
	info, err := os.Stat(homeDir)
	if err != nil {
		log.Errorln("[Clash Lib] SetHomeDir: %s : %+v", homeDir, err)
		return false
	}
	if !info.IsDir() {
		log.Errorln("[Clash Lib] SetHomeDir: Path is not directory %s", homeDir)
		return false
	}
	constant.SetHomeDir(homeDir)
	return true
}

//export SetConfig
func SetConfig(configStr *C.char) bool {
	configFile := C.GoString(configStr)
	if configFile == "" {
		return false
	}
	if !filepath.IsAbs(configFile) {
		configFile = filepath.Join(constant.Path.HomeDir(), configFile)
	}
	constant.SetConfig(configFile)
	return true
}

//export VerifyMMDB
func VerifyMMDB(path *C.char) bool {
	instance, err := geoip2.Open(C.GoString(path))
	if err == nil {
		_ = instance.Close()
	}
	return err == nil
}

//export StartController
func StartController(controller *C.char) {
	go route.Start(C.GoString(controller), "")
}

//export StartService
func StartService() bool {
	cfg, err := executor.Parse()
	if err != nil {
		log.Errorln("[Clash Lib] StartService: Parse config error: %+v", err)
		return false
	}
	executor.ApplyConfig(cfg, true)
	return true
}

func main() {}
