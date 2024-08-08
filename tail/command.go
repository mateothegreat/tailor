package tail

import (
	"bufio"
	"io"

	"github.com/fatih/color"
	"github.com/mateothegreat/tailer/ssh"
)

type TailConfig struct {
	Command string
	Color   color.Attribute
}

func Run(session *ssh.Session, config TailConfig) error {
	sess, err := session.Client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	stdout, err := sess.StdoutPipe()
	if err != nil {
		return err
	}

	err = sess.Start(config.Command)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		clr := color.New(config.Color)
		clr.Printf("%s: %s", session.Config.Name, line)
	}

	return nil
}
