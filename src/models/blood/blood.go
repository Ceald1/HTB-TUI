package blood

import (
	"fmt"
	"time"

	"context"

	"github.com/Ceald1/HTB-TUI/src/format"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	HTB "github.com/gubarz/gohtb"
	"github.com/gubarz/gohtb/services/seasons"
)

type errMsg error

type bloodMsg struct {
	User string
	Root string
	Err  error
}

type model struct {
	spinner      spinner.Model
	quit         bool
	err          error
	MachineBlood Blood
	HTBClient    HTB.Client
	Machine      seasons.ActiveMachineResponse
}

type Blood struct {
	User string
	Root string
}

func InitialModel(HTBClient *HTB.Client) model {
	s := spinner.New()
	custom_spinner := spinner.Spinner{
		Frames: []string{
`▁`, `▂`, `▃`, `▄`, `▅`, `▆`, `▇`, `█`, `▇`, `▆`, `▅`, `▄`, `▃`, `▁`,
},
FPS: time.Second / 10,
	}
	s.Spinner = custom_spinner
	return model{spinner: s, HTBClient: *HTBClient}
}

func (m model) Init() tea.Cmd {
	machine, err := SeasonalMachine(&m.HTBClient)
	if err != nil {
		m.err = err
		return m.spinner.Tick
	}
	m.Machine = machine
	return tea.Batch(
		m.spinner.Tick,
		bloodTaskCmd(&m.HTBClient, m.Machine), // Start the periodic check
	)
}

// Returns (machine, error)
func SeasonalMachine(client *HTB.Client) (seasons.ActiveMachineResponse, error) {
	ctx := context.Background()
	machine, err := client.Seasons.ActiveMachine(ctx)
	if err != nil {
		return machine, fmt.Errorf("error getting current machine! %w", err)
	}
	// Defensive: check machine.Data.Id (assuming Id is int, zero means invalid)
	if machine.Data.Id == 0 {
		return machine, fmt.Errorf("active machine Data.Id is zero")
	}
	return machine, nil
}

func bloodTaskCmd(client *HTB.Client, machine seasons.ActiveMachineResponse) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		machineID := machine.Data.Id

		activeMachineInfo := client.Machines.Machine(machineID)

		defer func() {
			if r := recover(); r != nil {
				// fmt.Printf("Recovered from panic in activeMachineInfo.Info: %v\n", r)
				
			}
		}()

		machineInfo, err := activeMachineInfo.Info(ctx)
		if err != nil {
			return bloodMsg{Err: fmt.Errorf("error getting Blood INFO: %w", err)}
		}

		var user, root string

		// Defensive: UserBlood and RootBlood are structs, not pointers
		// So check their fields, not nil
		if machineInfo.Data.UserBlood.User.Name != "" {
			user = machineInfo.Data.UserBlood.User.Name
		}
		if machineInfo.Data.RootBlood.User.Name != "" {
			root = machineInfo.Data.RootBlood.User.Name
		}

		return bloodMsg{User: user, Root: root}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quit = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	case bloodMsg:
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.MachineBlood.User = msg.User
		m.MachineBlood.Root = msg.Root
		if len(msg.User) != 0 && len(msg.Root) != 0 {
			m.quit = true
			return m, tea.Quit
		}
		// Schedule next task in 30 seconds
		return m, tea.Batch(
			m.spinner.Tick,
			func() tea.Msg {
				time.Sleep(30 * time.Second)
				return bloodTaskCmd(&m.HTBClient, m.Machine)()
			},
		)

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	var str string
	redStyle := lipgloss.NewStyle().Foreground(format.Red)
	str = redStyle.Render(fmt.Sprintf("\n\n   %s Awaiting Bloods...press q to quit\n\n user: %s\n\n root: %s", m.spinner.View(), m.MachineBlood.User, m.MachineBlood.Root))

	if m.quit {
		return str + "\n"
	}
	return str
}

func Run(HTBClient *HTB.Client) (err error) {
	p := tea.NewProgram(InitialModel(HTBClient))
	if _, err = p.Run(); err != nil {
		return err
	}
	return
}