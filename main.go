package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/playwright-community/playwright-go"
)

type Job struct {
	Name            string   `json:"name"`
	StartDate       string   `json:"startDate"`
	EndDate         string   `json:"endDate"`
	Title           string   `json:"title"`
	Accomplishments []string `json:"accomplishments"`
}

var (
	dailyChoices   []string
	ckChoices      []string
	sbChoices      []string
	adskEngChoices []string
	adskMgrChoices []string
	lgChoices      []string
	yahooEaChoices []string
)

func newAccomplishmentOptions(j Job) []huh.Option[string] {
	var options []huh.Option[string]

	for _, a := range j.Accomplishments {
		options = append(options, huh.NewOption(a, a).Selected(true))
	}

	return options
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

func main() {
	dailyJob := Job{
		Name:      "Daily",
		StartDate: "May 2022",
		EndDate:   "Nov 2023",
		Title:     "Sr. Engineering Manager / Director of Engineering",
		Accomplishments: []string{
			"Spearheaded the development of an AI-powered, HIPAA-compliant clinical notes API.",
			"Managed cross-functional teams to enhance platform scalability and incident response processes.",
			"Implemented strategic cost reductions and improved server-side software delivery.",
			"Enhanced production incident-response process and acted as an Incident Commander.",
			"Shepherded a group of Staff Engineers and Support leads tasked to scale the platform for new customer traffic.",
		},
	}

	ckJob := Job{
		Name:      "ConvertKit",
		StartDate: "Mar 2020",
		EndDate:   "Apr 2022",
		Title:     "Sr. Engineering Manager",
		Accomplishments: []string{
			"Orchestrated an engineering team reorganization to optimize performance and team cohesion.",
			"Led initiatives that maintained an email-sending error rate at approximately 0.0001%.",
			"Played a key role in hiring and developing senior engineering staff: managers and ICs.",
			"Negotiated 20% cost-reduction for the company’s most expensive third-party vendor.",
			"Led the promotion plan for the first Staff Software Engineer at the company.",
			"Helped define the initial job level/ladder for software engineers from Entry- to Staff Engineer.",
		},
	}

	sbJob := Job{
		Name:      "Scotch & Blend",
		StartDate: "Feb 2017",
		EndDate:   "Feb 2019",
		Title:     "Partner",
		Accomplishments: []string{
			"Conducted technical due diligence and team building for client business acquisitions.",
			"Hired the core engineering and leadership team members after a client’s IP-acquistion, as interim Dir. of Engineering.",
			"Designed HTTP APIs (GraphQL) for mobile game backends.",
		},
	}

	adskEngJob := Job{
		Name:      "Autodesk",
		StartDate: "Sept 2015",
		EndDate:   "Jul 2016",
		Title:     "Sr. Principal Software Engineer",
		Accomplishments: []string{
			"Reviewed and enhanced system designs across multiple teams, focusing on disaster recovery and CI/CD implementations.",
			"Built a background job processor to support bulk uploads of e-commerce product definitions and SKUs.",
			"Implemented Kanban for the DevOps team in order to focus the team’s efforts and ship consistently.",
		},
	}

	adskMgrJob := Job{
		Name:      "Autodesk, Live Gamer (acquired into Autodesk)",
		StartDate: "Jul 2013",
		EndDate:   "Aug 2015",
		Title:     "Development Manager",
		Accomplishments: []string{
			"Led acquisition-to-launch of subscriptions platform.",
			"Doubled engineering team using a hybrid work model.",
			"Enabled trust between product- and reliability-engineering teams that had a \"throw it over the wall\" mentality.",
			"Sponsored engineer transitions from the IC- to manager-track.",
			"Served as scrummaster for various teams; facilitating retrospectives and process improvements.",
			"Created CI/CD pipelines from scratch after our IP was acquired by Autodesk.",
		},
	}

	lgJob := Job{
		Name:      "Live Gamer, Twofish (acquired into Live Gamer)",
		StartDate: "Feb 2008",
		EndDate:   "Jul 2013",
		Title:     "Sr. Software Engineer",
		Accomplishments: []string{
			"Built rate-limiting, fraud, and payments systems and APIs.",
			"Built a two-part secret key implementation for 1-click purchasing.",
			"Created CI/CD pipelines from scratch.",
			"Designed HTTP APIs (REST) for virtual currency transactions.",
			"Created an aggregation model for per-user payment-transaction summaries to aid with rate-limiting and customer service.",
			"Built a raffling system for virtual-currencies and items.",
			"Built frontends for purchase flows.",
			"Wrote an email-sending backend for transactional emails.",
			"Wrote MySQL query optimizations (via indexing, removal of foreign keys, deadlock prevention).",
		},
	}

	yahooEaJob := Job{
		Name:      "Yahoo!, Electronic Arts",
		StartDate: "Jan 2007",
		EndDate:   "Jan 2008",
		Title:     "Software Engineer",
		Accomplishments: []string{
			"Built mobile frontends for Yahoo! News, Weather, and Sports.",
			"Optimized services for SMS and J2ME applications for mobile games.",
		},
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(dailyJob.Title).
				Options(newAccomplishmentOptions(dailyJob)...).
				Value(&dailyChoices),
			huh.NewMultiSelect[string]().
				Title(ckJob.Title).
				Options(newAccomplishmentOptions(ckJob)...).
				Value(&ckChoices),
			huh.NewMultiSelect[string]().
				Title(sbJob.Title).
				Options(newAccomplishmentOptions(sbJob)...).
				Value(&sbChoices),
			huh.NewMultiSelect[string]().
				Title(adskEngJob.Title).
				Options(newAccomplishmentOptions(adskEngJob)...).
				Value(&adskEngChoices),
			huh.NewMultiSelect[string]().
				Title(adskMgrJob.Title).
				Options(newAccomplishmentOptions(adskMgrJob)...).
				Value(&adskMgrChoices),
			huh.NewMultiSelect[string]().
				Title(lgJob.Title).
				Options(newAccomplishmentOptions(lgJob)...).
				Value(&lgChoices),
			huh.NewMultiSelect[string]().
				Title(yahooEaJob.Title).
				Options(newAccomplishmentOptions(yahooEaJob)...).
				Value(&yahooEaChoices),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	generateResume := func() {
		jobs := []Job{
			newJobView(dailyJob, dailyChoices),
			newJobView(ckJob, ckChoices),
			newJobView(sbJob, sbChoices),
			newJobView(adskEngJob, adskEngChoices),
			newJobView(adskMgrJob, adskMgrChoices),
			newJobView(lgJob, lgChoices),
			newJobView(yahooEaJob, yahooEaChoices),
		}

		templateFile := "resume.tmpl"
		t, err := template.New(templateFile).ParseFiles(templateFile)
		if err != nil {
			panic(err)
		}

		resumeText := &bytes.Buffer{}
		err = t.Execute(resumeText, jobs)
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

	_ = spinner.New().Title("Preparing your resume...").Action(generateResume).Run()
}

func assertErrorToNilf(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}
