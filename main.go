package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

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
		return container, fmt.Errorf(fmt.Sprintf("unable to read `%s`; make sure the file exists in the same directory as this program", fileName))
	}

	if err := json.Unmarshal(jsonData, &container); err != nil {
		return container, fmt.Errorf("unable to parse json in data.json; are the file contents malformed")
	}

	return container, nil
}

func returnErrorf(message string, err error) error {
	if err != nil {
		return fmt.Errorf(message+", file a bug at https://github.com/counterleft/resume-generator", err)
	}

	return nil
}

type ResumeOptions struct {
	TemplateFilename string
}

func generateResume(resumeData ResumeData, outputFilename string, options ResumeOptions) error {
	templateFile := options.TemplateFilename
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
		Path: playwright.String(outputFilename),
	})
	returnErrorf("could not create PDF: %w", err)
	returnErrorf("could not close browser: %w", browser.Close())
	returnErrorf("could not stop Playwright: %w", pw.Stop())

	return nil
}

var foregroundBaseStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#353535", Dark: "#FFF7D8"})

var errorMessageStyle = lipgloss.NewStyle().Inherit(foregroundBaseStyle)

var errorIconStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#8F2BF5", Dark: "#FF4D94"}).MarginRight(2)

func printErrorAndExit(err error) {
	if err != nil {
		message := fmt.Sprintf("We ran into some trouble: %s", err.Error())
		styledMessage := errorIconStyle.Render("!!") + errorMessageStyle.Render(message)
		fmt.Println(styledMessage)

		os.Exit(1)
	}
}

func parseCliArguments() (string, ResumeOptions) {
	templateFilename := flag.String("template", "resume.tmpl", "the html template to apply for this resume")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nArguments:\n")
		fmt.Fprintf(os.Stderr, "  <data_file>  path to the resume datafile\n")
	}

	flag.Parse()

	if len(flag.Args()) < 1 {
		printErrorAndExit(fmt.Errorf("no <resume.json> file specified"))
	}

	dataFilename := flag.Args()[0]

	options := ResumeOptions{
		TemplateFilename: *templateFilename,
	}

	return dataFilename, options
}

func main() {
	dataFilename, options := parseCliArguments()

	resumeData, err := readJsonFileInto(dataFilename, ResumeData{})
	printErrorAndExit(err)

	outputFilename := strings.Split(dataFilename, ".")[0] + ".pdf"

	startGeneration := func() {
		err = generateResume(resumeData, outputFilename, options)
		printErrorAndExit(err)
	}

	_ = spinner.New().Title("Preparing your resume...").Action(startGeneration).Run()

	successStyle := lipgloss.NewStyle().
		Inherit(foregroundBaseStyle).
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#8F2BF5", Dark: "#FF4D94"}).
		Padding(1)

	fmt.Println(successStyle.Render(fmt.Sprintf("All done! Open up `%s`.", outputFilename)))
}
