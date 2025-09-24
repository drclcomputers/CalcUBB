// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

/*
Algoritm pentru calculul notei la admiterea la Facultatea de Matematica si Informatica:
- Fiecare întrebare poate avea 1 sau mai multe răspunsuri corecte.
- Fiecare întrebare are asociat un punctaj „p” care este primit de către candidat daca bifează toate răspunsurile corecte si numai pe acestea.
- Daca o întrebare are asociat un punctaj „p” si are un număr de „t” răspunsuri corecte si un număr de „f” răspunsuri incorecte, atunci:
    > Daca unul dintre cele „t” răspunsuri corecte este bifat, atunci candidatul primește „p/t” puncte pentru acel raspuns
    > Daca unul dintre cele „f” răspunsuri incorecte este bifat, atunci candidatul primește „(-0.66)*p/t” puncte (adică este penalizat) pentru acel raspuns.
    > Punctajul pentru aceasta întrebare este minim „0” (daca rezultatul evaluării tuturor bifărilor făcute de către candidat este negativ atunci rezultatul se înlocuiește cu „0”) si maxim „p” ”( dacă candidatul bifează toate răspunsurile corecte și numai pe acestea).

In prezent, p=3.75 (24 de intrebari), dar in trecut p=3 (30 de intrebari).
Se acorda 10 puncte din oficiu.
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const version = "v25.09"

type state int

const (
	askNumProblems state = iota
	inputAnswers
	showResult
)

type answer struct {
	correct []bool
	marked  []bool
}

type checkboxItem struct {
	title    string
	selected bool
}

type model struct {
	state        state
	numProblems  int
	pPerProblem  float64
	problemIndex int
	answers      []answer
	numInput     textinput.Model
	correctItems []checkboxItem // Checkbox-uri pentru barem
	markedItems  []checkboxItem // Checkbox-uri pentru răspunsuri bifate
	focusedGroup int            // 0 = correctItems, 1 = markedItems
	focusedIndex int            // Indexul checkbox-ului selectat
	result       resultState
	resultIndex  int // Pentru navigarea în rezultate
	err          error
	quit         bool
	debugMessage string
}

func initialModel() model {
	numInput := textinput.New()
	numInput.Placeholder = "Număr probleme"
	numInput.Focus()
	numInput.CharLimit = 2
	numInput.Width = 20

	return model{
		state:        askNumProblems,
		numInput:     numInput,
		numProblems:  0,
		focusedGroup: 0,
		focusedIndex: 0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case askNumProblems:
		return m.updateNumProblems(msg)
	case inputAnswers:
		return m.updateInputAnswers(msg)
	case showResult:
		return m.updateShowResult(msg)
	}
	return m, nil
}

func (m model) updateNumProblems(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quit = true
			return m, tea.Quit
		case "enter":
			num, err := strconv.Atoi(m.numInput.Value())
			if err != nil || num > 90 || num < 1 {
				m.err = fmt.Errorf("Introdu un numar <=90 && >=10!")
				return m, nil
			}
			m.numProblems = num
			m.pPerProblem = 3.75
			if num != 24 {
				m.pPerProblem = float64(90) / float64(num)
			}
			m.answers = make([]answer, num)
			for i := range m.answers {
				m.answers[i] = answer{
					correct: make([]bool, 4),
					marked:  make([]bool, 4),
				}
			}
			m.state = inputAnswers
			m.correctItems = createCheckboxItems()
			m.markedItems = createCheckboxItems()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.numInput, cmd = m.numInput.Update(msg)
	return m, cmd
}

func (m model) updateInputAnswers(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "Q":
			m.quit = true
			return m, tea.Quit
		case "down", "m", "M":
			m.focusedGroup = (m.focusedGroup + 1) % 2
			return m, nil
		case "up", "j", "J":
			m.focusedGroup = (m.focusedGroup - 1 + 2) % 2
			return m, nil
		case "left", "<", "k", "K":
			m.focusedIndex = (m.focusedIndex - 1 + 4) % 4
			return m, nil
		case "right", ">", "l", "L":
			m.focusedIndex = (m.focusedIndex + 1) % 4
			return m, nil
		case " ", "a", "A":
			if m.focusedGroup == 0 {
				m.correctItems[m.focusedIndex].selected = !m.correctItems[m.focusedIndex].selected
			} else {
				m.markedItems[m.focusedIndex].selected = !m.markedItems[m.focusedIndex].selected
			}
			return m, nil
		case "tab", "n", "N":
			m.saveCurrentAnswers()
			if m.problemIndex < m.numProblems-1 {
				m.problemIndex++
				m.correctItems = createCheckboxItems()
				m.markedItems = createCheckboxItems()
				for i, selected := range m.answers[m.problemIndex].correct {
					m.correctItems[i].selected = selected
				}
				for i, selected := range m.answers[m.problemIndex].marked {
					m.markedItems[i].selected = selected
				}
				m.focusedIndex = 0
			}
			return m, nil
		case "shift+tab", "p", "P":
			m.saveCurrentAnswers()
			if m.problemIndex > 0 {
				m.problemIndex--
				m.correctItems = createCheckboxItems()
				m.markedItems = createCheckboxItems()
				for i, selected := range m.answers[m.problemIndex].correct {
					m.correctItems[i].selected = selected
				}
				for i, selected := range m.answers[m.problemIndex].marked {
					m.markedItems[i].selected = selected
				}
				m.focusedIndex = 0
			}
			return m, nil
		case "enter":
			m.saveCurrentAnswers()
			total, reports := m.calculateTotalScore()
			m.result = resultState{
				totalScore: total,
				reports:    reports,
			}
			m.state = showResult
			return m, nil
		}
	}
	return m, nil
}

type resultState struct {
	totalScore float64
	reports    []questionReport
}

func (m model) updateShowResult(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quit = true
			return m, tea.Quit
		case "r":
			return initialModel(), textinput.Blink
		case "left", "<":
			if m.resultIndex > 0 {
				m.resultIndex--
			}
			return m, nil
		case "right", ">":
			if m.resultIndex < len(m.result.reports)-1 {
				m.resultIndex++
			}
			return m, nil
		}
	}
	return m, nil
}

func (m *model) saveCurrentAnswers() {
	for i, item := range m.correctItems {
		m.answers[m.problemIndex].correct[i] = item.selected
	}
	for i, item := range m.markedItems {
		m.answers[m.problemIndex].marked[i] = item.selected
	}
}

type questionReport struct {
	questionNum    int
	score          float64
	correctAnswers []string
	markedAnswers  []string
	maxScore       float64
	correctCount   int
	incorrectCount int
}

func (m model) calculateQuestionScore(qNum int, ans answer) questionReport {
	report := questionReport{
		questionNum:    qNum + 1,
		maxScore:       m.pPerProblem,
		correctAnswers: make([]string, 0),
		markedAnswers:  make([]string, 0),
	}

	letters := []string{"A", "B", "C", "D"}

	t := 0.0
	for i, c := range ans.correct {
		if c {
			t++
			report.correctAnswers = append(report.correctAnswers, letters[i])
		}
	}

	if t == 0 {
		return report
	}

	correctMarked := 0.0
	incorrectMarked := 0.0
	for i, correct := range ans.correct {
		if ans.marked[i] {
			report.markedAnswers = append(report.markedAnswers, letters[i])
			if correct {
				correctMarked++
				report.correctCount++
			} else {
				incorrectMarked++
				report.incorrectCount++
			}
		}
	}

	score := (correctMarked * (m.pPerProblem / t)) + (incorrectMarked * (-0.66 * m.pPerProblem / t))
	report.score = math.Max(0, math.Min(score, m.pPerProblem))
	return report
}

func (m model) calculateTotalScore() (float64, []questionReport) {
	total := 0.0
	reports := make([]questionReport, len(m.answers))

	for i, ans := range m.answers {
		report := m.calculateQuestionScore(i, ans)
		reports[i] = report
		total += report.score
	}

	return total + 10, reports //punctele din oficiu
}

func createCheckboxItems() []checkboxItem {
	return []checkboxItem{
		{title: "A", selected: false},
		{title: "B", selected: false},
		{title: "C", selected: false},
		{title: "D", selected: false},
	}
}

var (
	checkboxStyle = lipgloss.NewStyle().Padding(0, 1)
	focusedStyle  = lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69"))
	inactiveStyle = lipgloss.NewStyle().Padding(0, 1)
)

func (m model) renderCheckbox(items []checkboxItem, focused bool, focusedIndex int) string {
	var s []string
	for i, item := range items {
		checkbox := fmt.Sprintf("[%s] %s", map[bool]string{true: "x", false: " "}[item.selected], item.title)
		if focused && i == focusedIndex {
			s = append(s, focusedStyle.Render(checkbox))
		} else {
			s = append(s, checkboxStyle.Render(checkbox))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, s...)
}

func (m model) View() string {
	var b strings.Builder

	switch m.state {
	case askNumProblems:
		b.WriteString("Calculator pentru punctajul de la examenul de admitere UBB (info si mate) - v 25.09\n\n")
		b.WriteString("Introdu numărul de probleme (24 sau 30): ")
		b.WriteString(m.numInput.View())
		if m.err != nil {
			b.WriteString(fmt.Sprintf("\nEroare: %v\n", m.err))
		}
		b.WriteString("\n\n\nApasă Enter pentru a continua, Q sau Ctrl+C pentru ieșire.")
	case inputAnswers:
		b.WriteString(fmt.Sprintf("Întrebarea %d/%d (p=%.2f puncte per pb)\n", m.problemIndex+1, m.numProblems, m.pPerProblem))
		b.WriteString("Barem (corecte):\n")
		b.WriteString(m.renderCheckbox(m.correctItems, m.focusedGroup == 0, m.focusedIndex))
		b.WriteString("\n\nRăspunsuri bifate:\n")
		b.WriteString(m.renderCheckbox(m.markedItems, m.focusedGroup == 1, m.focusedIndex))
		b.WriteString("\n\nComenzi:\n >Up/Down pentru scrolare intre raspunsuri corecte/bifate\n >Stânga(<)/Dreapta(>) pentru selecție\n >Space/A pentru bifare\n >Tab(N)/Shift+Tab(P) pentru următoarea/precedenta intrebare\n >Enter pentru calcul\n >Q pentru ieșire")
	case showResult:
		report := m.result.reports[m.resultIndex]

		b.WriteString(fmt.Sprintf("Punctaj total: %.2f (din care 10 puncte din oficiu)\n\n", m.result.totalScore))
		b.WriteString(fmt.Sprintf("Întrebarea %d din %d\n", m.resultIndex+1, len(m.result.reports)))
		b.WriteString("──────────────────────────────────────────────────\n")

		b.WriteString(fmt.Sprintf("Întrebarea %d:\n", report.questionNum))
		b.WriteString(fmt.Sprintf("  Răspunsuri corecte: %s\n", strings.Join(report.correctAnswers, ", ")))
		if len(report.markedAnswers) > 0 {
			b.WriteString(fmt.Sprintf("  Răspunsuri bifate: %s\n", strings.Join(report.markedAnswers, ", ")))
		} else {
			b.WriteString("  Răspunsuri bifate: (niciun răspuns bifat)\n")
		}
		b.WriteString(fmt.Sprintf("  Răspunsuri corecte: %d | Răspunsuri greșite: %d\n", report.correctCount, report.incorrectCount))
		b.WriteString(fmt.Sprintf("  Punctaj: %.2f / %.2f\n", report.score, report.maxScore))
		b.WriteString("\n──────────────────────────────────────────────────\n")

		if m.resultIndex > 0 {
			b.WriteString("\n← Întrebarea anterioară (stânga)")
		}
		if m.resultIndex < len(m.result.reports)-1 {
			if m.resultIndex > 0 {
				b.WriteString(" | ")
			} else {
				b.WriteString("\n")
			}
			b.WriteString("Următoarea întrebare → (dreapta)")
		}

		b.WriteString("\n\nApasă R pentru reset, Q sau Ctrl+C pentru ieșire.")
	}

	return b.String()
}

func readAnswersFromFile(filename string) ([]answer, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, fmt.Errorf("Nu pot deschide fisierul: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return nil, 0, fmt.Errorf("Fisierul este gol")
	}

	text := scanner.Text()
	numProblems, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil || (strings.Contains(text, "0") && strings.Contains(text, "1") && !strings.Contains(text, "10") && strings.Contains(text, "A") && strings.Contains(text, "F") && strings.Contains(text, "a") && strings.Contains(text, "f")) {
		return nil, 0, fmt.Errorf("Prima linie trebuie sa contina numarul de probleme!")
	}

	answers := make([]answer, numProblems)
	lineNum := 2

	for i := 0; i < numProblems; i++ {
		if !scanner.Scan() {
			return nil, 0, fmt.Errorf("Lipsesc raspunsuri la linia %d", lineNum)
		}

		line := strings.TrimSpace(scanner.Text())
		if len(line) != 4 {
			return nil, 0, fmt.Errorf("Linia %d trebuie sa contina exact 4 caractere", lineNum)
		}

		correctCount := 0
		answers[i].correct = make([]bool, 4)
		answers[i].marked = make([]bool, 4)

		for j, char := range line {
			if char != '0' && char != '1' && char != 'A' && char != 'F' && char != 'a' && char != 'f' {
				return nil, 0, fmt.Errorf("Linia %d contine caractere invalide (doar 0 sau 1 sunt permise)", lineNum)
			}
			if char == '1' || char == 'A' || char == 'a' {
				answers[i].correct[j] = true
				correctCount++
			}
		}

		if correctCount == 0 || correctCount > 3 {
			return nil, 0, fmt.Errorf("Cel putin un raspuns si maxim 3 pot fi corecte. Verifica linia %d (problema %d)", lineNum, lineNum-1)
		}

		lineNum++
	}

	return answers, numProblems, nil
}

func main() {
	versionFlag := flag.Bool("v", false, "Afiseaza versiunea")
	fileFlag := flag.String("f", "", "Citeste raspunsurile dintr-un fisier")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Calculator punctaj UBB %s\n", version)
		return
	}

	var initialState model
	if *fileFlag != "" {
		answers, numProblems, err := readAnswersFromFile(*fileFlag)
		if err != nil {
			fmt.Printf("Eroare: %v\n", err)
			os.Exit(1)
		}
		initialState = initialModel()
		initialState.numProblems = numProblems
		initialState.answers = answers
		initialState.pPerProblem = 3.75
		if numProblems != 24 {
			initialState.pPerProblem = float64(90) / float64(numProblems)
		}
		initialState.state = inputAnswers
		initialState.correctItems = createCheckboxItems()
		initialState.markedItems = createCheckboxItems()
		for i, selected := range answers[0].correct {
			initialState.correctItems[i].selected = selected
		}
	} else {
		initialState = initialModel()
	}

	p := tea.NewProgram(initialState, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Eroare: %v\n", err)
		os.Exit(1)
	}
}
