package main

import (
	"bytes"
	"encoding/json"
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

func getJobs(fileName string) []Job {
	jsonData, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	jobs := []Job{}
	if err := json.Unmarshal(jsonData, &jobs); err != nil {
		panic(err)
	}

	return jobs
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

func generateResume(jobs []Job, choices [][]string) {
	views := make([]Job, len(jobs))

	for i, job := range jobs {
		views[i] = newJobView(job, choices[i])
	}

	templateFile := "resume.tmpl"
	t, err := template.New(templateFile).ParseFiles(templateFile)
	if err != nil {
		panic(err)
	}

	resumeText := &bytes.Buffer{}
	err = t.Execute(resumeText, views)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("./playwright.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resumeText)
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
	_, err = page.Goto("file:///Users/brian/dev/me/resume-generator/playwright.html")
	assertErrorToNilf("could not goto: %w", err)
	_, err = page.PDF(playwright.PagePdfOptions{
		Path: playwright.String("playwright.pdf"),
	})
	assertErrorToNilf("could not create PDF: %w", err)
	assertErrorToNilf("could not close browser: %w", browser.Close())
	assertErrorToNilf("could not stop Playwright: %w", pw.Stop())
}

func main() {
	jobs := getJobs("jobs.json")
	choices := make([][]string, len(jobs))

	form := makeForm(jobs, choices)
	runForm(form)

	handleSubmission := func() {
		generateResume(jobs, choices)
	}

	_ = spinner.New().Title("Preparing your resume...").Action(handleSubmission).Run()
}
