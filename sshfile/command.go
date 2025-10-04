package sshfile

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/timeout"
	"github.com/charmbracelet/gum/sshfilepicker"
	"github.com/charmbracelet/gum/style"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

func VerifyHost(host string, remote net.Addr, key ssh.PublicKey) error {
	//
	// If you want to connect to new hosts.
	// here your should check new connections public keys
	// if the key not trusted you shuld return an error
	//

	// hostFound: is host in known hosts file.
	// err: error if key not in known hosts file OR host in known hosts file but key changed!
	hostFound, err := goph.CheckKnownHost(host, remote, key, "")

	// Host in known hosts but key mismatch!
	// Maybe because of MAN IN THE MIDDLE ATTACK!
	if hostFound && err != nil {
		return err
	}

	// handshake because public key already exists.
	if hostFound && err == nil {
		return nil
	}

	// Add the new host to known hosts file.
	return goph.AddKnownHost(host, remote, key, "")
}

// Run is the interface to picking a file.
func (o Options) Run() error {
	if !o.File && !o.Directory {
		return errors.New("at least one between --file and --directory must be set")
	}
	var (
		sshUser       string
		sshHost       string
		sshKey        string
		sshDefaultKey string
		auth          goph.Auth
		err           error
		sshClient     *goph.Client
		path          string
	)
	if !strings.Contains(o.SSHHost, "@") {
		sshUser = os.Getenv("USER")
		sshHost = o.SSHHost
	} else {
		hostInfo := strings.Split(o.SSHHost, "@")
		sshUser = hostInfo[0]
		sshHost = hostInfo[1]
	}

	keys := []string{"id_rsa", "id_dsa", "id_ecdsa", "id_ed25519"}
	keysRootFolder := os.ExpandEnv("HOME") + "/.ssh/"

	for _, k := range keys {
		keyTry := keysRootFolder + k
		if _, err := os.Stat(keyTry); err == nil {
			sshDefaultKey = keyTry
			break
		}
	}

	if len(o.SSHKeyfile) > 0 {
		if _, err := os.Stat(o.SSHKeyfile); err == nil {
			sshKey = o.SSHKeyfile
		}
	}
	if o.SSHKeyfile == "" && o.SSHPasswd == "" {
		if len(sshDefaultKey) > 0 {
			sshKey = sshDefaultKey
		}
	}

	if o.SSHAgent || goph.HasAgent() {
		auth, err = goph.UseAgent()
	} else if len(o.SSHPasswd) > 0 {
		auth = goph.Password(o.SSHPasswd)
	} else if len(sshKey) > 0 {
		auth, err = goph.Key(sshKey, "")
	} else {
		return errors.New("no authentication method given")
	}
	if err != nil {
		panic(err)
	}

	sshClient, err = goph.NewConn(&goph.Config{
		User:     sshUser,
		Addr:     sshHost,
		Port:     o.SSHPort,
		Auth:     auth,
		Callback: VerifyHost,
	})
	if err != nil {
		panic(err)
	}
	// Close client net connection
	defer sshClient.Close()

	if o.Path == "" {
		// Get current working dir
		command, err := sshClient.Command("pwd")
		if err != nil {
			panic(err)
		}
		out, err := command.CombinedOutput()
		if err != nil {
			panic(err)
		}
		path = strings.TrimSuffix(string(out), "\n")
	} else {
		path = o.Path
	}

	// if o.Path == "" {
	// 	o.Path = "."
	// }
	//
	// path, err := filepath.Abs(o.Path)
	// if err != nil {
	// 	return fmt.Errorf("file not found: %w", err)
	// }

	fp := sshfilepicker.New()
	fp.CurrentDirectory = path
	fp.Path = path
	fp.SetHeight(o.Height)
	fp.AutoHeight = o.Height == 0
	fp.Cursor = o.Cursor
	fp.DirAllowed = o.Directory
	fp.FileAllowed = o.File
	fp.ShowHidden = o.All
	fp.Styles = sshfilepicker.DefaultStyles()
	fp.Styles.Cursor = o.CursorStyle.ToLipgloss()
	fp.Styles.Symlink = o.SymlinkStyle.ToLipgloss()
	fp.Styles.Directory = o.DirectoryStyle.ToLipgloss()
	fp.Styles.File = o.FileStyle.ToLipgloss()
	fp.Styles.Permission = o.PermissionsStyle.ToLipgloss()
	fp.Styles.Selected = o.SelectedStyle.ToLipgloss()
	fp.Styles.FileSize = o.FileSizeStyle.ToLipgloss()
	top, right, bottom, left := style.ParsePadding(o.Padding)
	fp.SSH = sshClient
	m := model{
		sshfilepicker: fp,
		padding:       []int{top, right, bottom, left},
		showHelp:      o.ShowHelp,
		help:          help.New(),
		keymap:        defaultKeymap(),
		headerStyle:   o.HeaderStyle.ToLipgloss(),
		header:        o.Header,
	}

	ctx, cancel := timeout.Context(o.Timeout)
	defer cancel()

	tm, err := tea.NewProgram(
		&m,
		tea.WithOutput(os.Stderr),
		tea.WithContext(ctx),
	).Run()
	if err != nil {
		return fmt.Errorf("unable to pick selection: %w", err)
	}
	m = tm.(model)
	if m.selectedPath == "" {
		return errors.New("no file selected")
	}

	fmt.Println(m.selectedPath)
	return nil
}
