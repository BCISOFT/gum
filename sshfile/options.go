package sshfile

import (
	"time"

	"github.com/charmbracelet/gum/style"
)

// Options are the options for the file command.
type Options struct {
	SSHHost    string `arg:"" required:"" help:"ssh connection string user@host" env:"GUM_SSHFILE_HOST"`
	SSHPort    uint   `optional:"" help:"ssh port default 22" default:"22" env:"GUM_SSHFILE_PORT"`
	SSHPasswd  string `optional:"" help:"ssh password or passphrase" env:"GUM_SSHFILE_PASSWD"`
	SSHKeyfile string `optional:"" help:"private ssh key file" env:"GUM_SSHFILE_KEYFILE"`
	SSHAgent   bool   `optional:"" default:"false" help:"use ssh agent" env:"GUM_SSHFILE_AGENT"`

	// Path is the path to the folder / directory to begin traversing.
	Path string `arg:"" optional:"" name:"path" help:"The path to the folder to begin traversing" env:"GUM_SSHFILE_PATH"`
	// Cursor is the character to display in front of the current selected items.
	Cursor    string        `short:"c" help:"The cursor character" default:">" env:"GUM_SSHFILE_CURSOR"`
	All       bool          `short:"a" help:"Show hidden and 'dot' files" default:"false" env:"GUM_SSHFILE_ALL"`
	File      bool          `help:"Allow files selection" default:"false" env:"GUM_SSHFILE_FILE"`
	Directory bool          `help:"Allow directories selection" default:"false" env:"GUM_SSHFILE_DIRECTORY"`
	ShowHelp  bool          `help:"Show help key binds" negatable:"" default:"true" env:"GUM_SSHFILE_SHOW_HELP"`
	Timeout   time.Duration `help:"Timeout until command aborts without a selection" default:"0s" env:"GUM_SSHFILE_TIMEOUT"`
	Header    string        `help:"Header value" default:"" env:"GUM_SSHFILE_HEADER"`
	Height    int           `help:"Maximum number of files to display" default:"10" env:"GUM_SSHFILE_HEIGHT"`

	CursorStyle      style.Styles `embed:"" prefix:"cursor." help:"The cursor style" set:"defaultForeground=212" envprefix:"GUM_SSHFILE_CURSOR_"`
	SymlinkStyle     style.Styles `embed:"" prefix:"symlink." help:"The style to use for symlinks" set:"defaultForeground=36" envprefix:"GUM_SSHFILE_SYMLINK_"`
	DirectoryStyle   style.Styles `embed:"" prefix:"directory." help:"The style to use for directories" set:"defaultForeground=99" envprefix:"GUM_SSHFILE_DIRECTORY_"`
	FileStyle        style.Styles `embed:"" prefix:"file." help:"The style to use for files" envprefix:"GUM_SSHFILE_FILE_"`
	PermissionsStyle style.Styles `embed:"" prefix:"permissions." help:"The style to use for permissions" set:"defaultForeground=244" envprefix:"GUM_SSHFILE_PERMISSIONS_"`
	SelectedStyle    style.Styles `embed:"" prefix:"selected." help:"The style to use for the selected item" set:"defaultBold=true" set:"defaultForeground=212" envprefix:"GUM_SSHFILE_SELECTED_"`                    //nolint:staticcheck
	FileSizeStyle    style.Styles `embed:"" prefix:"file-size." help:"The style to use for file sizes" set:"defaultWidth=8" set:"defaultAlign=right" set:"defaultForeground=240"  envprefix:"GUM_SSHFILE_FILE_SIZE_"` //nolint:staticcheck
	HeaderStyle      style.Styles `embed:"" prefix:"header." set:"defaultForeground=99" envprefix:"GUM_SSHFILE_HEADER_"`
	Padding          string       `help:"Padding" default:"${defaultPadding}" group:"Style Flags" env:"GUM_SSHFILE_PADDING"`
}
