package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/log"
	"github.com/playwright-community/playwright-go"
)

type Job struct {
	Name            string   `json:"name"`
	StartDate       string   `json:"startDate"`
	EndDate         string   `json:"endDate"`
	Title           string   `json:"title"`
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
	Header    Header    `json:"header"`
	Education Education `json:"education"`
	Jobs      []Job     `json:"jobs"`
	Skills    []string  `json:"skills"`
}

func readJsonFileIntoContainer[T any](fileName string, container T) T {
	jsonData, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(jsonData, &container); err != nil {
		panic(err)
	}

	return container
}

func makeMultiSelects(jobs []Job, choices [][]string) []huh.Field {
	selects := make([]huh.Field, len(jobs))

	for i, job := range jobs {
		selects[i] = huh.NewMultiSelect[string]().
			Title(job.Title).
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

func makeForm(jobs []Job, choices [][]string) *huh.Form {
	selects := makeMultiSelects(jobs, choices)

	return huh.NewForm(
		huh.NewGroup(
			selects...,
		),
	)
}

func runForm(form *huh.Form) {
	err := form.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// XXX Do I need this?
func assertErrorToNilf(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

func newJobView(j Job, accomplishments []string) Job {
	return Job{
		Name:            j.Name,
		StartDate:       j.StartDate,
		EndDate:         j.EndDate,
		Title:           j.Title,
		Accomplishments: accomplishments,
	}
}

func generateResume(resumeData ResumeData) {
	templateFile := "resume.tmpl"
	template, err := template.New(templateFile).ParseFiles(templateFile)
	if err != nil {
		panic(err)
	}

	resumeText := &bytes.Buffer{}
	err = template.Execute(resumeText, resumeData)
	if err != nil {
		panic(err)
	}

	tempFile, err := os.CreateTemp("", "resume-generator-*.html")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resumeText)
	if err != nil {
		panic(err)
	}

	// Open with playwright and have it make a PDF
	pw, err := playwright.Run()
	assertErrorToNilf("could not launch playwright: %w", err)
	browser, err := pw.Chromium.Launch()
	assertErrorToNilf("could not launch Chromium: %w", err)
	context, err := browser.NewContext()
	assertErrorToNilf("could not create context: %w", err)
	page, err := context.NewPage()
	assertErrorToNilf("could not create page: %w", err)
	_, err = page.Goto(fmt.Sprintf("file://%s", tempFile.Name()))
	assertErrorToNilf("could not goto: %w", err)
	_, err = page.PDF(playwright.PagePdfOptions{
		Path: playwright.String("resume.pdf"),
	})
	assertErrorToNilf("could not create PDF: %w", err)
	assertErrorToNilf("could not close browser: %w", browser.Close())
	assertErrorToNilf("could not stop Playwright: %w", pw.Stop())
}

func main() {
	resumeData := readJsonFileIntoContainer("data.json", ResumeData{})

	choices := make([][]string, len(resumeData.Jobs))
	form := makeForm(resumeData.Jobs, choices)
	runForm(form)

	handleSubmission := func() {
		views := make([]Job, len(resumeData.Jobs))

		for i, job := range resumeData.Jobs {
			views[i] = newJobView(job, choices[i])
		}

		resumeData.Jobs = views

		generateResume(resumeData)
	}

	_ = spinner.New().Title("Preparing your resume...").Action(handleSubmission).Run()
}
