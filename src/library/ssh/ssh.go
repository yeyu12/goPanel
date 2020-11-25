package ssh

import (
	"bufio"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	gossh "golang.org/x/crypto/ssh"
	"net"
	"time"
	"unicode/utf8"
)

type Ssh struct {
	host    string
	user    string
	port    int
	passwd  string
	client  *gossh.Client
	session *gossh.Session
}

type TermConfig struct {
	Cols uint32
	Rows uint32
}

type ptyRequestMsg struct {
	Term     string
	Cols     uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

func NewSsh(host, user, passwd string, port int) *Ssh {
	return &Ssh{
		host:   host,
		user:   user,
		port:   port,
		passwd: passwd,
	}
}

func (s *Ssh) connect() (*Ssh, error) {
	config := &gossh.ClientConfig{}
	config.SetDefaults()
	config.User = s.user
	config.Auth = []gossh.AuthMethod{gossh.Password(s.passwd)}
	config.HostKeyCallback = func(hostname string, remote net.Addr, key gossh.PublicKey) error { return nil }
	client, err := gossh.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port), config)
	if err != nil {
		return nil, err
	}
	s.client = client

	return s, nil
}

func (s *Ssh) exec(cmd string) (string, error) {
	var buf bytes.Buffer
	session, err := s.client.NewSession()
	if nil != err {
		return "", err
	}

	session.Stdout = &buf
	session.Stderr = &buf
	err = session.Run(cmd)
	if err != nil {
		return "", err
	}
	defer session.Close()

	stdout := buf.String()

	return stdout, nil
}

func (s *Ssh) RunShell(term TermConfig) (gossh.Channel, error) {
	shConn, err := s.connect()
	if err != nil {
		return nil, err
	}

	channel, incomingRequests, err := shConn.client.Conn.OpenChannel("session", nil)
	if err != nil {
		return nil, err
	}

	go func() {
		for req := range incomingRequests {
			if req.WantReply {
				err = req.Reply(false, nil)
				if err != nil {
					return
				}
			}
		}
	}()

	err = s.Pty(channel, term)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (s *Ssh) Pty(channel gossh.Channel, term TermConfig) error {
	sshModes := gossh.TerminalModes{
		gossh.ECHO:          1,
		gossh.TTY_OP_ISPEED: 14400,
		gossh.TTY_OP_OSPEED: 14400,
	}

	var sshModeList []byte
	for k, v := range sshModes {
		kv := struct {
			Key byte
			Val uint32
		}{k, v}
		sshModeList = append(sshModeList, gossh.Marshal(&kv)...)
	}
	sshModeList = append(sshModeList, 0)

	req := ptyRequestMsg{
		Term:     "xterm",
		Cols:     term.Cols,
		Rows:     term.Rows,
		Width:    term.Cols * 8,
		Height:   term.Rows * 8,
		Modelist: string(sshModeList),
	}

	ok, err := channel.SendRequest("pty-req", true, gossh.Marshal(&req))
	if !ok || err != nil {
		return err
	}

	ok, err = channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		return err
	}

	return nil
}

func (s *Ssh) Read(channel gossh.Channel, sshRead chan []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	br := bufio.NewReader(channel)
	buf := []byte{}

	t := time.NewTimer(time.Millisecond * 100)
	defer t.Stop()
	rn := make(chan rune)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Info(err)
			}
		}()

		for {
			r, size, err := br.ReadRune()
			if err != nil {
				return
			}
			if size > 0 {
				rn <- r
			}
		}
	}()

	for {
		select {
		case <-t.C:
			if len(buf) != 0 {
				sshRead <- buf
				buf = []byte{}
			}

			t.Reset(time.Millisecond * 100)
		case d := <-rn:
			if d != utf8.RuneError {
				p := make([]byte, utf8.RuneLen(d))
				utf8.EncodeRune(p, d)
				buf = append(buf, p...)
			} else {
				buf = append(buf, []byte("@")...)
			}
		}
	}
}

func (s *Ssh) Write(channel gossh.Channel, sshWrite chan []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	for {
		select {
		case sw := <-sshWrite:
			if _, err := channel.Write(sw); nil != err {
				log.Error(err)
				return
			}
		}
	}
}
