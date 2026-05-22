package tui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
	"k8s.io/client-go/kubernetes"

	"github.com/kha7iq/pvu/internal/k8s"
	"github.com/kha7iq/pvu/internal/models"
)

type screen int

const (
	screenList screen = iota
	screenDetail
)

type Model struct {
	clientset kubernetes.Interface
	namespace string
	keys      KeyMap

	screen       screen
	pods         []models.PodListItem
	filteredPods []models.PodListItem
	selected     int
	loading      bool
	err          error

	filtering   bool
	filterInput textinput.Model

	detail models.PodVolumeView

	width    int
	height   int
	viewport viewport.Model
}

type podsLoadedMsg struct {
	pods []models.PodListItem
	err  error
}

type detailLoadedMsg struct {
	view models.PodVolumeView
	err  error
}

func NewModel(clientset kubernetes.Interface, namespace string) Model {
	vp := viewport.New(0, 0)

	input := textinput.New()
	input.Prompt = "/ "
	input.Placeholder = "Filter pods..."
	input.CharLimit = 64
	input.Blur()

	return Model{
		clientset:    clientset,
		namespace:    namespace,
		keys:         DefaultKeyMap(),
		screen:       screenList,
		pods:         []models.PodListItem{},
		filteredPods: []models.PodListItem{},
		selected:     0,
		loading:      true,
		filtering:    false,
		filterInput:  input,
		viewport:     vp,
	}
}

func (m Model) Init() tea.Cmd {
	return loadPodsCmd(m.clientset, m.namespace)
}

func loadPodsCmd(clientset kubernetes.Interface, namespace string) tea.Cmd {
	return func() tea.Msg {
		pods, err := k8s.ListPods(context.Background(), clientset, namespace)
		return podsLoadedMsg{
			pods: pods,
			err:  err,
		}
	}
}

func loadDetailCmd(clientset kubernetes.Interface, namespace, podName string) tea.Cmd {
	return func() tea.Msg {
		view, err := k8s.GetPodVolumeView(context.Background(), clientset, namespace, podName)
		return detailLoadedMsg{
			view: view,
			err:  err,
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.syncViewport()
		return m, nil

	case podsLoadedMsg:
		m.loading = false
		m.err = msg.err
		m.pods = msg.pods
		m.applyFilter()
		m.syncViewport()
		return m, nil

	case detailLoadedMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			m.detail = msg.view
			m.screen = screenDetail
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		switch m.screen {
		case screenList:
			return m.updateList(msg)
		case screenDetail:
			return m.updateDetail(msg)
		}
	}

	return m, nil
}

func (m Model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.filtering {
		switch msg.String() {
		case "esc":
			if m.filterInput.Value() == "" {
				m.filtering = false
				m.filterInput.Blur()
			} else {
				m.filterInput.SetValue("")
				m.applyFilter()
				m.syncViewport()
			}
			return m, nil

		case "enter":
			if len(m.filteredPods) == 0 {
				return m, nil
			}
			m.loading = true
			m.err = nil
			pod := m.filteredPods[m.selected]
			m.filtering = false
			m.filterInput.Blur()
			return m, loadDetailCmd(m.clientset, pod.Namespace, pod.Name)

		case "up", "k":
			if len(m.filteredPods) > 0 && m.selected > 0 {
				m.selected--
				m.syncViewport()
			}
			return m, nil

		case "down", "j":
			if len(m.filteredPods) > 0 && m.selected < len(m.filteredPods)-1 {
				m.selected++
				m.syncViewport()
			}
			return m, nil
		}

		var cmd tea.Cmd
		m.filterInput, cmd = m.filterInput.Update(msg)
		m.applyFilter()
		m.syncViewport()
		return m, cmd
	}

	switch msg.String() {
	case "/":
		m.filtering = true
		return m, m.filterInput.Focus()

	case "up", "k":
		if len(m.filteredPods) > 0 && m.selected > 0 {
			m.selected--
			m.syncViewport()
		}

	case "down", "j":
		if len(m.filteredPods) > 0 && m.selected < len(m.filteredPods)-1 {
			m.selected++
			m.syncViewport()
		}

	case "enter":
		if len(m.filteredPods) == 0 {
			return m, nil
		}
		m.loading = true
		m.err = nil
		pod := m.filteredPods[m.selected]
		return m, loadDetailCmd(m.clientset, pod.Namespace, pod.Name)
	}

	return m, nil
}

func (m Model) updateDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "b", "esc":
		m.screen = screenList
		m.err = nil
		m.syncViewport()
		return m, nil
	}

	return m, nil
}

func (m *Model) applyFilter() {
	query := m.filterInput.Value()

	if query == "" {
		m.filteredPods = append([]models.PodListItem(nil), m.pods...)
		if m.selected >= len(m.filteredPods) {
			m.selected = 0
		}
		return
	}

	type candidate struct {
		pod  models.PodListItem
		text string
	}

	candidates := make([]candidate, 0, len(m.pods))
	values := make([]string, 0, len(m.pods))

	for _, pod := range m.pods {
		text := pod.Name + " " + pod.Namespace
		candidates = append(candidates, candidate{
			pod:  pod,
			text: text,
		})
		values = append(values, text)
	}

	matches := fuzzy.Find(query, values)
	filtered := make([]models.PodListItem, 0, len(matches))

	for _, match := range matches {
		filtered = append(filtered, candidates[match.Index].pod)
	}

	m.filteredPods = filtered

	if len(m.filteredPods) == 0 {
		m.selected = 0
		m.viewport.SetYOffset(0)
		return
	}

	if m.selected >= len(m.filteredPods) {
		m.selected = len(m.filteredPods) - 1
	}
	if m.selected < 0 {
		m.selected = 0
	}
}

func (m *Model) syncViewport() {
	if m.screen != screenList {
		return
	}

	width := m.contentWidth()
	listContent := m.renderListContent(width)

	headerHeight := lipgloss.Height(m.renderListHeader())
	footerHeight := lipgloss.Height(m.renderListFooter())
	filterHeight := lipgloss.Height(m.renderFilterBar())
	verticalPadding := 6

	availableHeight := m.height - headerHeight - footerHeight - filterHeight - verticalPadding
	if availableHeight < 3 {
		availableHeight = 3
	}

	m.viewport.Width = width
	m.viewport.Height = availableHeight
	m.viewport.SetContent(listContent)

	m.ensureSelectionVisible()
}

func (m *Model) ensureSelectionVisible() {
	if len(m.filteredPods) == 0 {
		m.viewport.SetYOffset(0)
		return
	}

	if m.selected < m.viewport.YOffset {
		m.viewport.SetYOffset(m.selected)
		return
	}

	bottom := m.viewport.YOffset + m.viewport.Height - 1
	if m.selected > bottom {
		m.viewport.SetYOffset(m.selected - m.viewport.Height + 1)
	}
}

func Run(clientset kubernetes.Interface, namespace string) error {
	p := tea.NewProgram(NewModel(clientset, namespace))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("run tui: %w", err)
	}
	return nil
}
