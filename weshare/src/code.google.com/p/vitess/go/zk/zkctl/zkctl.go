// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Commands for controlling an external zookeeper process.
*/

package zkctl

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
	"time"

	"code.google.com/p/vitess/go/relog"
	"code.google.com/p/vitess/go/vt/env"
	"launchpad.net/gozk/zookeeper"
)

const (
	StartWaitTime    = 20 // number of seconds to wait at Start
	ShutdownWaitTime = 20 // number of seconds to wait at Shutdown
)

func init() {
	zookeeper.SetLogLevel(zookeeper.LOG_WARN)
}

type Zkd struct {
	config *ZkConfig
	//  createConnection CreateConnection
}

// createConnection: closure that returns a new connection or an error
func NewZkd(config *ZkConfig) *Zkd {
	return &Zkd{config}
}

func (zkd *Zkd) LocalClientAddr() string {
	return fmt.Sprintf("localhost:%v", zkd.config.ClientPort)
}

/*
 ZOO_LOG_DIR=""
 ZOO_CFG="/.../zoo.cfg"
 ZOOMAIN="org.apache.zookeeper.server.quorum.QuorumPeerMain"
 java -DZOO_LOG_DIR=${ZOO_LOG_DIR} -cp $CLASSPATH $ZOOMAIN $YT_ZK_CFG
*/
func (zkd *Zkd) Start() error {
	relog.Info("zkctl.Start")
	// NOTE(msolomon) use a script here so we can detach and continue to run
	// if the wrangler process dies. this pretty much the same as mysqld_safe.
	args := []string{
		zkd.config.LogDir(),
		zkd.config.ConfigFile(),
		zkd.config.PidFile(),
	}
	root, err := env.VtRoot()
	if err != nil {
		return err
	}
	dir := path.Join(root, "bin")
	cmd := exec.Command(path.Join(dir, "zksrv.sh"), args...)
	cmd.Env = os.Environ()
	cmd.Dir = dir

	if err = cmd.Start(); err != nil {
		return err
	}

	// give it some time to succeed - usually by the time the socket emerges
	// we are in good shape
	for i := 0; i < StartWaitTime; i++ {
		zkAddr := fmt.Sprintf(":%v", zkd.config.ClientPort)
		conn, connErr := net.Dial("tcp", zkAddr)
		if connErr != nil {
			err = connErr
			time.Sleep(time.Second)
			continue
		} else {
			err = nil
			conn.Write([]byte("ruok"))
			reply := make([]byte, 4)
			conn.Read(reply)
			if string(reply) != "imok" {
				err = fmt.Errorf("local zk unhealthy: %v %v", zkAddr, reply)
			}
			conn.Close()
			break
		}
	}
	// wait so we don't get a bunch of defunct processes
	go cmd.Wait()
	return err
}

func (zkd *Zkd) Shutdown() error {
	relog.Info("zkctl.Shutdown")
	pidData, err := ioutil.ReadFile(zkd.config.PidFile())
	if err != nil {
		return err
	}
	pid, err := strconv.Atoi(string(bytes.TrimSpace(pidData)))
	if err != nil {
		return err
	}
	err = syscall.Kill(pid, 9)
	if err != nil && err != syscall.ESRCH {
		return err
	}
	proc, _ := os.FindProcess(pid)
	processIsDead := false
	for i := 0; i < ShutdownWaitTime; i++ {
		if proc.Signal(syscall.Signal(0)) == syscall.ESRCH {
			processIsDead = true
			break
		}
		time.Sleep(time.Second)
	}
	if !processIsDead {
		return fmt.Errorf("Shutdown didn't kill process %v", pid)
	}
	return nil
}

func (zkd *Zkd) makeCfg() (string, error) {
	root, err := env.VtRoot()
	if err != nil {
		return "", err
	}
	cnfTemplatePaths := []string{path.Join(root, "config/zkcfg/zoo.cfg")}
	return MakeZooCfg(cnfTemplatePaths, zkd.config, "# generated by vt")
}

func (zkd *Zkd) ConfigChanged() (bool, error) {
	oldCfg, err := ioutil.ReadFile(zkd.config.ConfigFile())
	if err != nil {
		if pathErr, ok := err.(*os.PathError); ok && pathErr.Err == syscall.ENOENT {
			oldCfg = []byte{}
			err = nil
		} else {
			return false, err
		}
	}
	currentCfg, err := zkd.makeCfg()
	if err != nil {
		return false, err
	}
	return string(oldCfg) != currentCfg, nil
}

func (zkd *Zkd) Init() error {
	if zkd.Inited() {
		return fmt.Errorf("zk already inited")
	}
	return zkd.init(false)
}

func (zkd *Zkd) Reinit(preserveData bool) (err error) {
	if preserveData {
		err = zkd.Shutdown()
	} else {
		err = zkd.Teardown()
	}
	if err == nil {
		err = zkd.init(preserveData)
	}
	return
}

func (zkd *Zkd) init(preserveData bool) error {
	relog.Info("zkd.Init")
	for _, path := range zkd.config.DirectoryList() {
		if err := os.MkdirAll(path, 0775); err != nil {
			relog.Error(err.Error())
			return err
		}
		// FIXME(msolomon) validate permissions?
	}

	configData, err := zkd.makeCfg()
	if err == nil {
		err = ioutil.WriteFile(zkd.config.ConfigFile(), []byte(configData), 0664)
	}
	if err != nil {
		relog.Error("failed creating %v: %v", zkd.config.ConfigFile(), err)
		return err
	}

	err = zkd.config.WriteMyid()
	if err != nil {
		relog.Error("failed creating %v: %v", zkd.config.MyidFile(), err)
		return err
	}

	if err = zkd.Start(); err != nil {
		relog.Error("failed starting, check %v", zkd.config.LogDir())
		return err
	}

	zkAddr := fmt.Sprintf("localhost:%v", zkd.config.ClientPort)
	zk, session, err := zookeeper.Dial(zkAddr, StartWaitTime*time.Second)
	if err != nil {
		return err
	}
	event := <-session
	if event.State != zookeeper.STATE_CONNECTED {
		return err
	}
	defer zk.Close()

	if !preserveData {
		_, err = zk.Create("/zk", "", 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
		if err != nil && !zookeeper.IsError(err, zookeeper.ZNODEEXISTS) {
			return err
		}

		if zkd.config.Global {
			_, err = zk.Create("/zk/global", "", 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
			if err != nil && !zookeeper.IsError(err, zookeeper.ZNODEEXISTS) {
				return err
			}
		}
	}

	return nil
}

func (zkd *Zkd) Teardown() error {
	relog.Info("zkctl.Teardown")
	if err := zkd.Shutdown(); err != nil {
		relog.Warning("failed zookeeper shutdown: %v", err.Error())
	}
	var removalErr error
	for _, dir := range zkd.config.DirectoryList() {
		relog.Debug("remove data dir %v", dir)
		if err := os.RemoveAll(dir); err != nil {
			relog.Error("failed removing %v: %v", dir, err.Error())
			removalErr = err
		}
	}
	return removalErr
}

func (zkd *Zkd) CheckProcess() error {
	pidFile := zkd.config.PidFile()
	data, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return err
	}

	pid, err := strconv.Atoi(string(data))
	// found a pid - if the process is burned, fast-fail
	// otherwise, try to connect and fail slowly
	if err == nil {
		_, err := os.FindProcess(pid)
		if err != nil {
			return err
		}
	}

	zkAddr := fmt.Sprintf("localhost:%v", zkd.config.ClientPort)
	zk, session, err := zookeeper.Dial(zkAddr, StartWaitTime*time.Second)
	if err != nil {
		return err
	}
	defer zk.Close()
	timer := time.NewTimer(StartWaitTime * 1e9)
	defer timer.Stop()
	select {
	case event := <-session:
		if event.State != zookeeper.STATE_CONNECTED {
			return err
		}
	case <-timer.C:
		return errors.New("zk deadline exceeded connecting to " + zkAddr)
	}
	_, _, err = zk.Get("/zk")
	return err
}

func (zkd *Zkd) Inited() bool {
	myidFile := zkd.config.MyidFile()
	_, statErr := os.Stat(myidFile)
	if statErr == nil {
		return true
	} else if statErr.(*os.PathError).Err != syscall.ENOENT {
		panic("can't access file " + myidFile + ": " + statErr.Error())
	}
	return false
}