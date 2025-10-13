package users
import (
	"context"

	"github.com/Ceald1/HTB-TUI/src/format"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	HTB "github.com/gubarz/gohtb"
	"github.com/gubarz/gohtb/services/users"
)

var (
	ctx = context.Background()
)

const (
	quit_value = 9999999999999
)

