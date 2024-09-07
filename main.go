package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/playwright-community/playwright-go"
)

type Job struct {
	Company         string   `json:"company"`
	StartDate       string   `json:"startDate"`
	EndDate         string   `json:"endDate"`
	Title           string   `json:"title"`
	Summary         string   `json:"summary"`
	Accomplishments []string `json:"accomplishments"`
}

type Header struct {
	Name        string `json:"name"`
	Headline    string `json:"headline"`
	LinkedInUrl string `json:"linkedInUrl"`
	Email       string `json:"email"`
	Location    string `json:"location"`
	Phone       string `json:"phone"`
}

type Education struct {
	School    string `json:"school"`
	Degree    string `json:"degree"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type ResumeData struct {
	Header    Header        `json:"header"`
	Education Education     `json:"education"`
	Summary   template.HTML `json:"summary"`
	Jobs      []Job         `json:"jobs"`
	Skills    []string      `json:"skills"`
}

func readJsonFileInto[T any](fileName string, container T) (T, error) {
	jsonData, err := os.ReadFile(fileName)
	if err != nil {
		return container, fmt.Errorf("unable to read data.json; make sure the file exists in the same directory as this program")
	}

	if err := json.Unmarshal(jsonData, &container); err != nil {
		return container, fmt.Errorf("unable to parse json in data.json; are the file contents malformed")
	}

	return container, nil
}

func newMultiSelects(jobs []Job, choices [][]string) []huh.Field {
	selects := make([]huh.Field, len(jobs))

	for i, job := range jobs {
		selects[i] = huh.NewMultiSelect[string]().
			Title(fmt.Sprintf("%s (%s)", job.Title, job.Company)).
			Options(newAccomplishmentOptions(job)...).
			Value(&choices[i])
	}

	return selects
}

func newAccomplishmentOptions(j Job) []huh.Option[string] {
	var options []huh.Option[string]

	for _, a := range j.Accomplishments {
		options = append(options, huh.NewOption(a, a).Selected(true))
	}

	return options
}

func newForm(jobs []Job, choices [][]string) *huh.Form {
	selects := newMultiSelects(jobs, choices)

	foreground := lipgloss.CompleteAdaptiveColor{
		Light: lipgloss.CompleteColor{TrueColor: "#8F2BF5", ANSI256: "55", ANSI: "45"},
		Dark:  lipgloss.CompleteColor{TrueColor: "#F27EDE", ANSI256: "213", ANSI: "105"},
	}

	highlighted := lipgloss.NewStyle().Foreground(foreground)

	desc := fmt.Sprintf("We've loaded your resume data.\nPlease %s you want in your generated resume.\nAll have been selected by default.", highlighted.Render("pick the accomplishments"))

	fields := []huh.Field{
		huh.NewNote().
			Title("Welcome to Resume Generator").
			Description(desc),
	}

	fields = append(fields, selects...)

	form := huh.NewForm(
		huh.NewGroup(
			fields...,
		),
	).WithProgramOptions(tea.WithAltScreen())

	if os.Getenv("ACCESSIBLE") != "" {
		form.WithAccessible(true)
	}

	return form
}

func returnErrorf(message string, err error) error {
	if err != nil {
		return fmt.Errorf(message+", file a bug at https://github.com/counterleft/resume-generator", err)
	}

	return nil
}

// func newJobView(j Job, accomplishments []string) Job {
// 	return Job{
// 		Company:         j.Company,
// 		StartDate:       j.StartDate,
// 		EndDate:         j.EndDate,
// 		Title:           j.Title,
// 		Summary:         j.Summary,
// 		Accomplishments: accomplishments,
// 	}
// }

func generateResume(resumeData ResumeData) error {
	templateFile := "resume.tmpl"
	template, err := template.New(templateFile).ParseFiles(templateFile)
	returnErrorf("unable to parse the template file", err)

	resumeText := &bytes.Buffer{}
	err = template.Execute(resumeText, resumeData)
	returnErrorf("unable to execute the template file", err)

	tempFile, err := os.CreateTemp("", "resume-generator-*.html")
	returnErrorf("unable to create temporary html file of resume", err)
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resumeText)
	returnErrorf("unable to copy resume data to temp html file", err)

	// Open with playwright and have it make a PDF
	pw, err := playwright.Run()
	returnErrorf("could not launch playwright: %w", err)

	browser, err := pw.Chromium.Launch()
	returnErrorf("could not launch Chromium: %w", err)

	context, err := browser.NewContext()
	returnErrorf("could not create browser context: %w", err)

	page, err := context.NewPage()
	returnErrorf("could not create browser page: %w", err)

	_, err = page.Goto(fmt.Sprintf("file://%s", tempFile.Name()))
	returnErrorf("could not goto: %w", err)

	_, err = page.PDF(playwright.PagePdfOptions{
		Path: playwright.String("resume.pdf"),
	})
	returnErrorf("could not create PDF: %w", err)
	returnErrorf("could not close browser: %w", browser.Close())
	returnErrorf("could not stop Playwright: %w", pw.Stop())

	return nil
}

func printErrorAndExit(err error) {
	if err != nil {
		fmt.Printf("We ran into some trouble: %s\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 2 {
		printErrorAndExit(fmt.Errorf("no <resume.json> file specified"))
	}

	dataFileName := os.Args[1]

	resumeData, err := readJsonFileInto(dataFileName, ResumeData{})
	printErrorAndExit(err)

	// choices := make([][]string, len(resumeData.Jobs))
	// form := newForm(resumeData.Jobs, choices)
	// err = form.Run()
	// printErrorAndExit(err)

	handleSubmission := func() {
		// views := make([]Job, len(resumeData.Jobs))

		// for i, job := range resumeData.Jobs {
		// 	views[i] = newJobView(job, choices[i])
		// }

		// resumeData.Jobs = views

		err = generateResume(resumeData)
		printErrorAndExit(err)
	}

	_ = spinner.New().Title("Preparing your resume...").Action(handleSubmission).Run()

	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.AdaptiveColor{Light: "#353535", Dark: "#FFF7D8"}).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#8F2BF5", Dark: "#FF4D94"}).
		Padding(1)

	fmt.Println(style.Render("All done! Open up `resume.pdf`."))
}
