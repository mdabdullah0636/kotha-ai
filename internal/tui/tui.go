package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"kotha/internal/app"
	"kotha/internal/config"
	"kotha/internal/llm/agent"
	"kotha/internal/logging"
	"kotha/internal/permission"
	"kotha/internal/pubsub"
	"kotha/internal/session"
	"kotha/internal/tui/components/chat"
	"kotha/internal/tui/components/core"
	"kotha/internal/tui/components/dialog"
	"kotha/internal/tui/layout"
	"kotha/internal/tui/manager"
	"kotha/internal/tui/page"
	"kotha/internal/tui/theme"
	"kotha/internal/tui/util"
)

type keyMap struct {
	Logs          key.Binding
	Quit          key.Binding
	Help          key.Binding
	SwitchSession key.Binding
	Commands      key.Binding
	Filepicker    key.Binding
	Models        key.Binding
	SwitchTheme   key.Binding
}

type startCompactSessionMsg struct{}

const (
	quitKey = "q"
)

var keys = keyMap{
	Logs: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("ctrl+l", "logs"),
	),

	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("ctrl+_", "ctrl+h"),
		key.WithHelp("ctrl+?", "toggle help"),
	),

	SwitchSession: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "switch session"),
	),

	Commands: key.NewBinding(
		key.WithKeys("ctrl+k"),
		key.WithHelp("ctrl+k", "commands"),
	),
	Filepicker: key.NewBinding(
		key.WithKeys("ctrl+f"),
		key.WithHelp("ctrl+f", "select files to upload"),
	),
	Models: key.NewBinding(
		key.WithKeys("ctrl+o"),
		key.WithHelp("ctrl+o", "model selection"),
	),

	SwitchTheme: key.NewBinding(
		key.WithKeys("ctrl+t"),
		key.WithHelp("ctrl+t", "switch theme"),
	),
}

var helpEsc = key.NewBinding(
	key.WithKeys("?"),
	key.WithHelp("?", "toggle help"),
)

var returnKey = key.NewBinding(
	key.WithKeys("esc"),
	key.WithHelp("esc", "close"),
)

var logsKeyReturnKey = key.NewBinding(
	key.WithKeys("esc", "backspace", quitKey),
	key.WithHelp("esc/q", "go back"),
)

type appModel struct {
	width, height      int
	previousPage       page.PageID
	pageManager        *manager.PageManager
	dialogManager      *manager.DialogManager
	status             core.StatusCmp
	app                *app.App
	selectedSession    session.Session

	permissions          dialog.PermissionDialogCmp
	help                 dialog.HelpCmp
	quit                 dialog.QuitDialog
	sessionDialog        dialog.SessionDialog
	commandDialog        dialog.CommandDialog
	commands             []dialog.Command
	modelDialog          dialog.ModelDialog
	initDialog           dialog.InitDialogCmp
	filepicker           dialog.FilepickerCmp
	themeDialog          dialog.ThemeDialog
	multiArgumentsDialog dialog.MultiArgumentsDialogCmp

	isCompacting      bool
	compactingMessage string
}

func (a appModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmd := a.pageManager.LoadPage(a.pageManager.CurrentPage())
	cmds = append(cmds, cmd)
	cmd = a.status.Init()
	cmds = append(cmds, cmd)
	cmd = a.quit.Init()
	cmds = append(cmds, cmd)
	cmd = a.help.Init()
	cmds = append(cmds, cmd)
	cmd = a.sessionDialog.Init()
	cmds = append(cmds, cmd)
	cmd = a.commandDialog.Init()
	cmds = append(cmds, cmd)
	cmd = a.modelDialog.Init()
	cmds = append(cmds, cmd)
	cmd = a.initDialog.Init()
	cmds = append(cmds, cmd)
	cmd = a.filepicker.Init()
	cmds = append(cmds, cmd)
	cmd = a.themeDialog.Init()
	cmds = append(cmds, cmd)

	cmds = append(cmds, func() tea.Msg {
		shouldShow, err := config.ShouldShowInitDialog()
		if err != nil {
			return util.InfoMsg{
				Type: util.InfoTypeError,
				Msg:  "Failed to check init status: " + err.Error(),
			}
		}
		return dialog.ShowInitDialogMsg{Show: shouldShow}
	})

	return tea.Batch(cmds...)
}

func (a appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		msg.Height -= 1 // Make space for the status bar
		a.width, a.height = msg.Width, msg.Height

		s, _ := a.status.Update(msg)
		a.status = s.(core.StatusCmp)
		currPage := a.pageManager.CurrentPage()
		currModel := a.pageManager.CurrentModel()
		currModel, cmd = currModel.Update(msg)
		a.pageManager.SetPage(currPage, currModel)
		cmds = append(cmds, cmd)

		prm, permCmd := a.permissions.Update(msg)
		a.permissions = prm.(dialog.PermissionDialogCmp)
		cmds = append(cmds, permCmd)

		help, helpCmd := a.help.Update(msg)
		a.help = help.(dialog.HelpCmp)
		cmds = append(cmds, helpCmd)

		session, sessionCmd := a.sessionDialog.Update(msg)
		a.sessionDialog = session.(dialog.SessionDialog)
		cmds = append(cmds, sessionCmd)

		command, commandCmd := a.commandDialog.Update(msg)
		a.commandDialog = command.(dialog.CommandDialog)
		cmds = append(cmds, commandCmd)

		filepicker, filepickerCmd := a.filepicker.Update(msg)
		a.filepicker = filepicker.(dialog.FilepickerCmp)
		cmds = append(cmds, filepickerCmd)

		a.dialogManager.SetSizeAll(msg.Width, msg.Height)

		if a.dialogManager.IsActive("multiArgumentsDialog") {
			a.multiArgumentsDialog.SetSize(msg.Width, msg.Height)
			args, argsCmd := a.multiArgumentsDialog.Update(msg)
			a.multiArgumentsDialog = args.(dialog.MultiArgumentsDialogCmp)
			cmds = append(cmds, argsCmd, a.multiArgumentsDialog.Init())
		}

		return a, tea.Batch(cmds...)
	// Status
	case util.InfoMsg:
		s, cmd := a.status.Update(msg)
		a.status = s.(core.StatusCmp)
		cmds = append(cmds, cmd)
		return a, tea.Batch(cmds...)
	case pubsub.Event[logging.LogMessage]:
		if msg.Payload.Persist {
			switch msg.Payload.Level {
			case "error":
				s, cmd := a.status.Update(util.InfoMsg{
					Type: util.InfoTypeError,
					Msg:  msg.Payload.Message,
					TTL:  msg.Payload.PersistTime,
				})
				a.status = s.(core.StatusCmp)
				cmds = append(cmds, cmd)
			case "info":
				s, cmd := a.status.Update(util.InfoMsg{
					Type: util.InfoTypeInfo,
					Msg:  msg.Payload.Message,
					TTL:  msg.Payload.PersistTime,
				})
				a.status = s.(core.StatusCmp)
				cmds = append(cmds, cmd)

			case "warn":
				s, cmd := a.status.Update(util.InfoMsg{
					Type: util.InfoTypeWarn,
					Msg:  msg.Payload.Message,
					TTL:  msg.Payload.PersistTime,
				})

				a.status = s.(core.StatusCmp)
				cmds = append(cmds, cmd)
			default:
				s, cmd := a.status.Update(util.InfoMsg{
					Type: util.InfoTypeInfo,
					Msg:  msg.Payload.Message,
					TTL:  msg.Payload.PersistTime,
				})
				a.status = s.(core.StatusCmp)
				cmds = append(cmds, cmd)
			}
		}
	case util.ClearStatusMsg:
		s, _ := a.status.Update(msg)
		a.status = s.(core.StatusCmp)

	// Permission
	case pubsub.Event[permission.PermissionRequest]:
		if err := a.dialogManager.Show("permissions"); err != nil {
			return a, util.ReportError(err)
		}
		return a, a.permissions.SetPermissions(msg.Payload)
	case dialog.PermissionResponseMsg:
		var cmd tea.Cmd
		switch msg.Action {
		case dialog.PermissionAllow:
			a.app.Permissions.Grant(msg.Permission)
		case dialog.PermissionAllowForSession:
			a.app.Permissions.GrantPersistant(msg.Permission)
		case dialog.PermissionDeny:
			a.app.Permissions.Deny(msg.Permission)
		}
		a.dialogManager.Hide("permissions")
		return a, cmd

	case page.PageChangeMsg:
		return a, a.moveToPage(msg.ID)

	case dialog.CloseQuitMsg:
		a.dialogManager.Hide("quit")
		return a, nil

	case dialog.CloseSessionDialogMsg:
		a.dialogManager.Hide("sessionDialog")
		return a, nil

	case dialog.CloseCommandDialogMsg:
		a.dialogManager.Hide("commandDialog")
		return a, nil

	case startCompactSessionMsg:
		// Start compacting the current session
		a.isCompacting = true
		a.compactingMessage = "Starting summarization..."

		if a.selectedSession.ID == "" {
			a.isCompacting = false
			return a, util.ReportWarn("No active session to summarize")
		}

		// Start the summarization process
		return a, func() tea.Msg {
			ctx := context.Background()
			a.app.CoderAgent.Summarize(ctx, a.selectedSession.ID)
			return nil
		}

	case pubsub.Event[agent.AgentEvent]:
		payload := msg.Payload
		if payload.Error != nil {
			a.isCompacting = false
			return a, util.ReportError(payload.Error)
		}

		a.compactingMessage = payload.Progress

		if payload.Done && payload.Type == agent.AgentEventTypeSummarize {
			a.isCompacting = false
			return a, util.ReportInfo("Session summarization complete")
		} else if payload.Done && payload.Type == agent.AgentEventTypeResponse && a.selectedSession.ID != "" {
			model := a.app.CoderAgent.Model()
			contextWindow := model.ContextWindow
			tokens := a.selectedSession.CompletionTokens + a.selectedSession.PromptTokens
			if (tokens >= int64(float64(contextWindow)*0.95)) && config.Get().AutoCompact {
				return a, util.CmdHandler(startCompactSessionMsg{})
			}
		}
		// Continue listening for events
		return a, nil

	case dialog.CloseThemeDialogMsg:
		a.dialogManager.Hide("themeDialog")
		return a, nil

	case dialog.ThemeChangedMsg:
		currPage := a.pageManager.CurrentPage()
		currModel := a.pageManager.CurrentModel()
		currModel, cmd = currModel.Update(msg)
		a.pageManager.SetPage(currPage, currModel)
		a.dialogManager.Hide("themeDialog")
		return a, tea.Batch(cmd, util.ReportInfo("Theme changed to: "+msg.ThemeName))

	case dialog.CloseModelDialogMsg:
		a.dialogManager.Hide("modelDialog")
		return a, nil

	case dialog.ModelSelectedMsg:
		a.dialogManager.Hide("modelDialog")

		model, err := a.app.CoderAgent.Update(config.AgentCoder, msg.Model.ID)
		if err != nil {
			return a, util.ReportError(err)
		}

		return a, util.ReportInfo(fmt.Sprintf("Model changed to %s", model.Name))

	case dialog.ShowInitDialogMsg:
		if msg.Show {
			if err := a.dialogManager.Show("initDialog"); err != nil {
				return a, util.ReportError(err)
			}
		} else {
			a.dialogManager.Hide("initDialog")
		}
		return a, nil

	case dialog.CloseInitDialogMsg:
		a.dialogManager.Hide("initDialog")
		if msg.Initialize {
			// Run the initialization command
			for _, cmd := range a.commands {
				if cmd.ID == "init" {
					// Mark the project as initialized
					if err := config.MarkProjectInitialized(); err != nil {
						return a, util.ReportError(err)
					}
					return a, cmd.Handler(cmd)
				}
			}
		} else {
			// Mark the project as initialized without running the command
			if err := config.MarkProjectInitialized(); err != nil {
				return a, util.ReportError(err)
			}
		}
		return a, nil

	case chat.SessionSelectedMsg:
		a.selectedSession = msg
		a.sessionDialog.SetSelectedSession(msg.ID)

	case pubsub.Event[session.Session]:
		if msg.Type == pubsub.EventTypeUpdated && msg.Payload.ID == a.selectedSession.ID {
			a.selectedSession = msg.Payload
		}
	case dialog.SessionSelectedMsg:
		a.dialogManager.Hide("sessionDialog")
		if a.pageManager.CurrentPage() == page.ChatPage {
			return a, util.CmdHandler(chat.SessionSelectedMsg(msg.Session))
		}
		return a, nil

	case dialog.CommandSelectedMsg:
		a.dialogManager.Hide("commandDialog")
		// Execute the command handler if available
		if msg.Command.Handler != nil {
			return a, msg.Command.Handler(msg.Command)
		}
		return a, util.ReportInfo("Command selected: " + msg.Command.Title)

	case dialog.ShowMultiArgumentsDialogMsg:
		// Show multi-arguments dialog
		a.multiArgumentsDialog = dialog.NewMultiArgumentsDialogCmp(msg.CommandID, msg.Content, msg.ArgNames)
		if err := a.dialogManager.Show("multiArgumentsDialog"); err != nil {
			return a, util.ReportError(err)
		}
		return a, a.multiArgumentsDialog.Init()

	case dialog.CloseMultiArgumentsDialogMsg:
		// Close multi-arguments dialog
		a.dialogManager.Hide("multiArgumentsDialog")

		// If submitted, replace all named arguments and run the command
		if msg.Submit {
			content := msg.Content

			// Replace each named argument with its value
			for name, value := range msg.Args {
				placeholder := "$" + name
				content = strings.ReplaceAll(content, placeholder, value)
			}

			// Execute the command with arguments
			return a, util.CmdHandler(dialog.CommandRunCustomMsg{
				Content: content,
				Args:    msg.Args,
			})
		}
		return a, nil

	case tea.KeyMsg:
		if a.dialogManager.IsActive("multiArgumentsDialog") {
			args, cmd := a.multiArgumentsDialog.Update(msg)
			a.multiArgumentsDialog = args.(dialog.MultiArgumentsDialogCmp)
			return a, cmd
		}

		switch {

		case key.Matches(msg, keys.Quit):
			a.toggleDialog("quit")
			a.closeDialog("help")
			a.closeDialog("sessionDialog")
			a.closeDialog("commandDialog")
			a.closeDialog("filepicker")
			a.closeDialog("modelDialog")
			a.closeDialog("multiArgumentsDialog")
			a.filepicker.ToggleFilepicker(false)
			return a, nil
		case key.Matches(msg, keys.SwitchSession):
			if a.pageManager.CurrentPage() == page.ChatPage && a.dialogsInactive("quit", "permissions", "commandDialog") {
				sessions, err := a.app.Sessions.List(context.Background())
				if err != nil {
					return a, util.ReportError(err)
				}
				if len(sessions) == 0 {
					return a, util.ReportWarn("No sessions available")
				}
				a.sessionDialog.SetSessions(sessions)
				if err := a.dialogManager.Show("sessionDialog"); err != nil {
					return a, util.ReportError(err)
				}
				return a, nil
			}
			return a, nil
		case key.Matches(msg, keys.Commands):
			if a.pageManager.CurrentPage() == page.ChatPage && a.dialogsInactive("quit", "permissions", "sessionDialog", "themeDialog", "filepicker") {
				if len(a.commands) == 0 {
					return a, util.ReportWarn("No commands available")
				}
				a.commandDialog.SetCommands(a.commands)
				if err := a.dialogManager.Show("commandDialog"); err != nil {
					return a, util.ReportError(err)
				}
				return a, nil
			}
			return a, nil
		case key.Matches(msg, keys.Models):
			if a.dialogManager.IsActive("modelDialog") {
				a.dialogManager.Hide("modelDialog")
				return a, nil
			}
			if a.pageManager.CurrentPage() == page.ChatPage && a.dialogsInactive("quit", "permissions", "sessionDialog", "commandDialog") {
				if err := a.dialogManager.Show("modelDialog"); err != nil {
					return a, util.ReportError(err)
				}
				return a, nil
			}
			return a, nil
		case key.Matches(msg, keys.SwitchTheme):
			if a.dialogsInactive("quit", "permissions", "sessionDialog", "commandDialog") {
				if err := a.dialogManager.Show("themeDialog"); err != nil {
					return a, util.ReportError(err)
				}
				return a, a.themeDialog.Init()
			}
			return a, nil
		case key.Matches(msg, returnKey) || key.Matches(msg):
			if msg.String() == quitKey {
				if a.pageManager.CurrentPage() == page.LogsPage {
					return a, a.moveToPage(page.ChatPage)
				}
			} else if !a.filepicker.IsCWDFocused() {
				if a.dialogManager.IsActive("quit") {
					a.toggleDialog("quit")
					return a, nil
				}
				if a.dialogManager.IsActive("help") {
					a.toggleDialog("help")
					return a, nil
				}
				if a.dialogManager.IsActive("initDialog") {
					a.dialogManager.Hide("initDialog")
					if err := config.MarkProjectInitialized(); err != nil {
						return a, util.ReportError(err)
					}
					return a, nil
				}
				if a.dialogManager.IsActive("filepicker") {
					a.dialogManager.Hide("filepicker")
					a.filepicker.ToggleFilepicker(false)
					return a, nil
				}
				if a.pageManager.CurrentPage() == page.LogsPage {
					return a, a.moveToPage(page.ChatPage)
				}
			}
		case key.Matches(msg, keys.Logs):
			return a, a.moveToPage(page.LogsPage)
		case key.Matches(msg, keys.Help):
			if a.dialogManager.IsActive("quit") {
				return a, nil
			}
			a.toggleDialog("help")
			return a, nil
		case key.Matches(msg, helpEsc):
			if a.app.CoderAgent.IsBusy() {
				if a.dialogManager.IsActive("quit") {
					return a, nil
				}
				a.toggleDialog("help")
				return a, nil
			}
		case key.Matches(msg, keys.Filepicker):
			a.toggleDialog("filepicker")
			if a.dialogManager.IsActive("filepicker") {
				a.filepicker.ToggleFilepicker(true)
			} else {
				a.filepicker.ToggleFilepicker(false)
			}
			return a, nil
		}
	default:
		f, filepickerCmd := a.filepicker.Update(msg)
		a.filepicker = f.(dialog.FilepickerCmp)
		cmds = append(cmds, filepickerCmd)

	}

	if a.dialogManager.IsActive("filepicker") {
		f, filepickerCmd := a.filepicker.Update(msg)
		a.filepicker = f.(dialog.FilepickerCmp)
		cmds = append(cmds, filepickerCmd)
		if _, ok := msg.(tea.KeyMsg); ok {
			return a, tea.Batch(cmds...)
		}
	}

	if a.dialogManager.IsActive("quit") {
		q, quitCmd := a.quit.Update(msg)
		a.quit = q.(dialog.QuitDialog)
		cmds = append(cmds, quitCmd)
		if _, ok := msg.(tea.KeyMsg); ok {
			return a, tea.Batch(cmds...)
		}
	}
	if a.dialogManager.IsActive("permissions") {
		d, permissionsCmd := a.permissions.Update(msg)
		a.permissions = d.(dialog.PermissionDialogCmp)
		cmds = append(cmds, permissionsCmd)
		if _, ok := msg.(tea.KeyMsg); ok {
			return a, tea.Batch(cmds...)
		}
	}

	if a.dialogManager.IsActive("sessionDialog") {
		d, sessionCmd := a.sessionDialog.Update(msg)
		a.sessionDialog = d.(dialog.SessionDialog)
		cmds = append(cmds, sessionCmd)
		if _, ok := msg.(tea.KeyMsg); ok {
			return a, tea.Batch(cmds...)
		}
	}

	if a.dialogManager.IsActive("commandDialog") {
		d, commandCmd := a.commandDialog.Update(msg)
		a.commandDialog = d.(dialog.CommandDialog)
		cmds = append(cmds, commandCmd)
		if _, ok := msg.(tea.KeyMsg); ok {
			return a, tea.Batch(cmds...)
		}
	}

	if a.dialogManager.IsActive("modelDialog") {
		d, modelCmd := a.modelDialog.Update(msg)
		a.modelDialog = d.(dialog.ModelDialog)
		cmds = append(cmds, modelCmd)
		if _, ok := msg.(tea.KeyMsg); ok {
			return a, tea.Batch(cmds...)
		}
	}

	if a.dialogManager.IsActive("initDialog") {
		d, initCmd := a.initDialog.Update(msg)
		a.initDialog = d.(dialog.InitDialogCmp)
		cmds = append(cmds, initCmd)
		if _, ok := msg.(tea.KeyMsg); ok {
			return a, tea.Batch(cmds...)
		}
	}

	if a.dialogManager.IsActive("themeDialog") {
		d, themeCmd := a.themeDialog.Update(msg)
		a.themeDialog = d.(dialog.ThemeDialog)
		cmds = append(cmds, themeCmd)
		if _, ok := msg.(tea.KeyMsg); ok {
			return a, tea.Batch(cmds...)
		}
	}

	s, _ := a.status.Update(msg)
	a.status = s.(core.StatusCmp)
	currPage := a.pageManager.CurrentPage()
	currModel := a.pageManager.CurrentModel()
	currModel, cmd = currModel.Update(msg)
	a.pageManager.SetPage(currPage, currModel)
	cmds = append(cmds, cmd)
	return a, tea.Batch(cmds...)
}

// RegisterCommand adds a command to the command dialog
func (a *appModel) RegisterCommand(cmd dialog.Command) {
	a.commands = append(a.commands, cmd)
}

func (a *appModel) findCommand(id string) (dialog.Command, bool) {
	for _, cmd := range a.commands {
		if cmd.ID == id {
			return cmd, true
		}
	}
	return dialog.Command{}, false
}

func (a *appModel) toggleDialog(name string) tea.Cmd {
	if a.dialogManager.IsActive(name) {
		a.dialogManager.Hide(name)
		return nil
	}
	if err := a.dialogManager.Show(name); err != nil {
		return util.ReportError(err)
	}
	return nil
}

func (a *appModel) closeDialog(name string) {
	a.dialogManager.Hide(name)
}

func (a *appModel) dialogsInactive(names ...string) bool {
	for _, n := range names {
		if a.dialogManager.IsActive(n) {
			return false
		}
	}
	return true
}

func (a *appModel) moveToPage(pageID page.PageID) tea.Cmd {
	cmd, err := a.pageManager.MoveTo(pageID, func() bool { return a.app.CoderAgent.IsBusy() })
	if err != nil {
		return util.ReportWarn(err.Error())
	}
	a.previousPage = a.pageManager.CurrentPage()
	sizeCmd := a.pageManager.SetPageSize(pageID, a.width, a.height)
	if sizeCmd != nil {
		return tea.Batch(cmd, sizeCmd)
	}
	return cmd
}

func (a appModel) View() string {
	components := []string{
		a.pageManager.CurrentModel().View(),
	}

	components = append(components, a.status.View())

	appView := lipgloss.JoinVertical(lipgloss.Top, components...)

	if a.dialogManager.IsActive("help") {
		bindings := layout.KeyMapToSlice(keys)
		if p, ok := a.pageManager.CurrentModel().(layout.Bindings); ok {
			bindings = append(bindings, p.BindingKeys()...)
		}
		if a.dialogManager.IsActive("permissions") {
			bindings = append(bindings, a.permissions.BindingKeys()...)
		}
		if a.pageManager.CurrentPage() == page.LogsPage {
			bindings = append(bindings, logsKeyReturnKey)
		}
		if !a.app.CoderAgent.IsBusy() {
			bindings = append(bindings, helpEsc)
		}
		a.help.SetBindings(bindings)
	}

	appView = a.dialogManager.View(appView)

	if a.isCompacting {
		t := theme.CurrentTheme()
		style := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.BorderFocused()).
			BorderBackground(t.Background()).
			Padding(1, 2).
			Background(t.Background()).
			Foreground(t.Text())

		overlay := style.Render("Summarizing\n" + a.compactingMessage)
		row := lipgloss.Height(appView) / 2
		row -= lipgloss.Height(overlay) / 2
		col := lipgloss.Width(appView) / 2
		col -= lipgloss.Width(overlay) / 2
		appView = layout.PlaceOverlay(
			col,
			row,
			overlay,
			appView,
			true,
		)
	}

	if a.dialogManager.IsActive("initDialog") {
		overlay := a.initDialog.View()
		appView = layout.PlaceOverlay(
			a.width/2-lipgloss.Width(overlay)/2,
			a.height/2-lipgloss.Height(overlay)/2,
			overlay,
			appView,
			true,
		)
	}

	return appView
}

func New(app *app.App) tea.Model {
	startPage := page.ChatPage
	pageManager := manager.NewPageManager(startPage)
	pageManager.RegisterPage(page.ChatPage, page.NewChatPage(app))
	pageManager.RegisterPage(page.LogsPage, page.NewLogsPage())

	model := &appModel{
		pageManager:    pageManager,
		status:         core.NewStatusCmp(app.LSPClients),
		help:           dialog.NewHelpCmp(),
		quit:           dialog.NewQuitCmp(),
		sessionDialog:  dialog.NewSessionDialogCmp(),
		commandDialog:  dialog.NewCommandDialogCmp(),
		modelDialog:    dialog.NewModelDialogCmp(),
		permissions:    dialog.NewPermissionDialogCmp(),
		initDialog:     dialog.NewInitDialogCmp(),
		themeDialog:    dialog.NewThemeDialogCmp(),
		app:            app,
		commands:       []dialog.Command{},
		filepicker:     dialog.NewFilepickerCmp(app),
	}

	dialogManager := manager.NewDialogManager(5)
	names := []string{"help", "permissions", "quit", "sessionDialog", "commandDialog", "modelDialog", "initDialog", "themeDialog", "filepicker", "multiArgumentsDialog"}
	views := []func() string{
		model.help.View, model.permissions.View, model.quit.View,
		model.sessionDialog.View, model.commandDialog.View, model.modelDialog.View,
		model.initDialog.View, model.themeDialog.View, model.filepicker.View, model.multiArgumentsDialog.View,
	}
	for i, name := range names {
		dialogManager.Register(name, manager.DialogConfig{
			Name: name,
			View: views[i],
		})
	}
	model.dialogManager = dialogManager

	model.RegisterCommand(dialog.Command{
		ID:          "init",
		Title:       "Initialize Project",
		Description: "Create/Update the Kotha.md memory file",
		Handler: func(cmd dialog.Command) tea.Cmd {
			prompt := `Please analyze this codebase and create a Kotha.md file containing:
1. Build/lint/test commands - especially for running a single test
2. Code style guidelines including imports, formatting, types, naming conventions, error handling, etc.

The file you create will be given to agentic coding agents (such as yourself) that operate in this repository. Make it about 20 lines long.
If there's already a kotha.md, improve it.
If there are Cursor rules (in .cursor/rules/ or .cursorrules) or Copilot rules (in .github/copilot-instructions.md), make sure to include them.`
			return tea.Batch(
				util.CmdHandler(chat.SendMsg{
					Text: prompt,
				}),
			)
		},
	})

	model.RegisterCommand(dialog.Command{
		ID:          "compact",
		Title:       "Compact Session",
		Description: "Summarize the current session and create a new one with the summary",
		Handler: func(cmd dialog.Command) tea.Cmd {
			return func() tea.Msg {
				return startCompactSessionMsg{}
			}
		},
	})
	// Load custom commands
	customCommands, err := dialog.LoadCustomCommands()
	if err != nil {
		logging.Warn("Failed to load custom commands", "error", err)
	} else {
		for _, cmd := range customCommands {
			model.RegisterCommand(cmd)
		}
	}

	return model
}
