package blood

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/Ceald1/HTB-TUI/src/format"
	HTB "github.com/gubarz/gohtb"
	"context"
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
}

type Blood struct {
	User string
	Root string
}

func InitialModel(HTBClient *HTB.Client) model {
	s := spinner.New()
	s.Spinner = spinner.Pulse
	return model{spinner: s, HTBClient: *HTBClient}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		bloodTaskCmd(m.HTBClient), // Start the periodic check
	)
}

func bloodTaskCmd(client HTB.Client) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		machine, err := client.Seasons.ActiveMachine(ctx)
		if err != nil {
			fmt.Println("error getting current machine!")
			fmt.Println(machine.Data)
			panic(err)
			return bloodMsg{Err: err}
		}
		machineID := machine.Data.Id
		activeMachineInfo := client.Machines.Machine(machineID)
		machineInfo, err := activeMachineInfo.Info(ctx)
		if err != nil {
			fmt.Println("error getting Blood INFO")
			panic(err)
			return bloodMsg{Err: err}
		}
		root := machineInfo.Data.RootBlood.User.Name
		user := machineInfo.Data.UserBlood.User.Name
		return bloodMsg{User: user, Root: root}
	}
}

// func every3SecondsCmd() tea.Cmd {
// 	return func() tea.Msg {
// 		time.Sleep(3 * time.Second)
// 		return nil // triggers Update so we can re-run BloodTask
// 	}
// }

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
		// Schedule next task in 3 seconds
		return m, tea.Batch(
			m.spinner.Tick,
			func() tea.Msg {
				time.Sleep(30 * time.Second)
				return bloodTaskCmd(m.HTBClient)()
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
	if len(m.MachineBlood.User) < 1 && len(m.MachineBlood.Root) < 1 {
		str = redStyle.Render(fmt.Sprintf("\n\n   %s Awaiting Bloods...press q to quit\n\n", m.spinner.View()))
	} else {
		str = redStyle.Render(fmt.Sprintf("\n\n   %s Awaiting Bloods...press q to quit\n\n user: %s\n root: %s", m.spinner.View(), m.MachineBlood.User, m.MachineBlood.Root))
	}
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