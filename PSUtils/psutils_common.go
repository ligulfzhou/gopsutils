package PSUtils

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

// not implemented
func (ps *PSUtils) SudoExec() (string, error) {
	return "not implemented", nil
}

func (ps *PSUtils) Exec(command string) (string, error) {

	session, err := ps.client.NewSession()
	if err != nil {
		return "", err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return "", err
	}

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		return "", err
	}

	// fmt.Println("\n\n\n---------------exec cmd----------------")
	// fmt.Printf("exec cmd: %s, res: %s \n", cmd, b.String())
	// fmt.Println("---------------exec cmd---------------- \n\n\n ")
	return StripString(b.String()), nil
}

func (ps *PSUtils) FileContent(filename string) (string, error) {
	str, err := ps.Exec(fmt.Sprintf("cat %s", filename))
	if err != nil {
		return "", err
	}

	return str, nil
}

func (ps *PSUtils) Glob(fileReg string) ([]string, error) {
	c, err := ps.Exec("ls " + fileReg)
	if err != nil {
		return nil, err
	}

	lines := SplitStringToLines(c)
	return lines, nil
}

func (ps *PSUtils) ReadLines(filename string) ([]string, error) {
	str, err := ps.FileContent(filename)
	if err != nil {
		return nil, err
	}

	contents := SplitStringToLines(str)
	return contents, nil
}

func (ps *PSUtils) FileExists(filename string) bool {
	_, err := ps.Exec(fmt.Sprintf("stat %s", filename))
	// _, err := ps.Exec(fmt.Sprintf("ls %s", filename))
	if err != nil {
		return false
	}

	return true
}

func (ps *PSUtils) ListDirectorys(dir string) ([]string, error) {
	// c, err := ps.Exec("ls -d / ")
	c, err := ps.Exec(fmt.Sprintf("ls -d %s/*/", dir))
	if err != nil {
		return nil, err
	}

	names := SplitStringWithDeeperLines(c)
	return names, nil
}

func (ps *PSUtils) NumProcs() int64 {
	var cnt int64

	names, err := ps.ListDirectorys("/proc")
	if err != nil {
		return 0
	}

	for _, v := range names {
		sp := strings.Split(v, "/")
		if len(sp) < 4 {
			continue
		}
		if _, err = strconv.ParseInt(sp[2], 10, 64); err == nil {
			cnt++
		}
	}

	return cnt
}

func (ps *PSUtils) GetVirtualization() (string, string) {
	if ps.VirtualizationSystem != "" || ps.VirtualizationRole != "" {
		return ps.VirtualizationSystem, ps.VirtualizationRole
	}

	system, role := ps.Virtualization()
	ps.VirtualizationSystem = system
	ps.VirtualizationRole = role
	return system, role
}

func (ps *PSUtils) Virtualization() (string, string) {
	var system, role string

	// /proc/xen
	if ps.FileExists("/proc/xen") {
		system = "xen"
		role = "guest"
		if ps.FileExists("/proc/xen/capabilities") {
			content, err := ps.FileContent("/proc/xen/capabilities")
			if err == nil {
				if strings.Contains(content, "control_id") {
					role = "host"
				}
			}
		}
		return system, role
	}

	if ps.FileExists("/proc/modules") {
		content, err := ps.FileContent("/proc/cpuinfo")
		flag := true
		if err == nil {
			if strings.Contains(content, "kvm") {
				system = "kvm"
				role = "host"
			} else if strings.Contains(content, "vboxdrv") {
				system = "vbox"
				role = "host"
			} else if strings.Contains(content, "vboxguest") {
				system = "vbox"
				role = "guest"
			} else if strings.Contains(content, "vmware") {
				system = "vmware"
				role = "guest"
			} else {
				flag = false
			}
		}
		if flag {
			return system, role
		}
	}

	if ps.FileExists("/proc/cpuinfo") {
		contents, err := ps.FileContent("/proc/cpuinfo")
		if err == nil {
			if strings.Contains(contents, "QEMU Virtual CPU") ||
				strings.Contains(contents, "Common KVM processor") ||
				strings.Contains(contents, "Common 32-bit KVM processor") {
				system = "kvm"
				role = "guest"
				return system, role
			}
		}
	}

	if ps.FileExists("/proc/bus/pci/devices") {
		contents, err := ps.FileContent("/proc/bus/pci/devices")
		if err == nil {
			if strings.Contains(contents, "virtio-pci") {
				role = "guest"
			}
		}
	}

	if ps.FileExists("/proc/bc/0") {
		system = "openvz"
		role = "host"
		return system, role
	} else if ps.FileExists("/proc/vz") {
		system = "openvz"
		role = "guest"
		return system, role
	}

	if ps.FileExists("/proc/self/status") {
		contents, err := ps.FileContent("/proc/self/status")
		if err == nil {
			if strings.Contains(contents, "s_context:") ||
				strings.Contains(contents, "VxID:") {
				system = "linux-vserver"
				return system, role
			}
			// TODO: guest or host
		}
	}

	if ps.FileExists("/proc/1/environ") {
		contents, err := ps.FileContent("/proc/1/environ")
		if err == nil {
			if strings.Contains(contents, "container=lxc") {
				system = "lxc"
				role = "guest"
				return system, role
			}
		}
	}

	if ps.FileExists("/proc/self/cgroup") {
		contents, err := ps.FileContent("/proc/self/cgroup")
		flagCgroup := true
		if err == nil {
			if strings.Contains(contents, "lxc") {
				system = "lxc"
				role = "guest"
			} else if strings.Contains(contents, "docker") {
				system = "docker"
				role = "guest"
			} else if strings.Contains(contents, "machine-rkt") {
				system = "rkt"
				role = "guest"
			} else if ps.FileExists("/usr/bin/lxc-version") {
				system = "lxc"
				role = "host"
			} else {
				flagCgroup = false
			}
		}
		if flagCgroup {
			return system, role
		}
	}

	if ps.FileExists("/etc/os-release") {
		pv := ps.GetOSRelease()
		if pv != nil && pv[0] == "coreos" {
			system = "rkt"
			role = "host"
			return system, role
		}
	}

	return system, role
}
