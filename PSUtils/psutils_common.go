package PSUtils

import (
	"bytes"
	"fmt"
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

	seq := "\r\n"
	if !strings.Contains(seq, "\r\n") {
		seq = "\n"
	}
	lines := strings.Split(c, seq)
	return lines, nil
}

func (ps *PSUtils) ReadLines(filename string) ([]string, error) {
	str, err := ps.FileContent(filename)
	if err != nil {
		return nil, err
	}

	seq := "\r\n"
	if !strings.Contains(seq, "\r\n") {
		seq = "\n"
	}
	contents := strings.Split(str, seq)
	return contents, nil
}

func (ps *PSUtils) FileExists(filename string) bool {
	_, err := ps.Exec(fmt.Sprintf("stat %s", filename))
	if err != nil {
		return false
	}

	return true
}

func (ps *PSUtils) GetOSRelease() []string {
	var platform, version string
	contents, err := ps.ReadLines("/etc/os-release")
	if err != nil {
		return []string{"", ""}
	}

	for _, line := range contents {
		field := strings.Split(line, "=")
		if len(field) < 2 {
			continue
		}
		switch field[0] {
		case "ID": // use ID for lowercase
			platform = trimQuotes(field[1])
		case "VERSION":
			version = trimQuotes(field[1])
		}
	}
	return []string{platform, version}
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}
