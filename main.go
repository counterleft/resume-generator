package main

import "github.com/charmbracelet/huh"
import "log"
import "os"

var (
	daily []string
	ck    []string
)

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Options(
					huh.NewOption("Spearheaded the development of an AI-powered, HIPAA-compliant clinical notes API.", "1"),
					huh.NewOption("Managed cross-functional teams to enhance platform scalability and incident response processes.", "2"),
					huh.NewOption("Implemented strategic cost reductions and improved server-side software delivery.", "3"),
					huh.NewOption("Enhanced production incident-response process and acted as an Incident Commander", "4"),
					huh.NewOption("Shepherded a group of Staff Engineers and Support leads tasked to scale the platform for new customer traffic.", "5"),
				).
				Title("Daily").
				Value(&daily),
			huh.NewMultiSelect[string]().
				Options(
					huh.NewOption("Orchestrated an engineering team reorganization to optimize performance and team cohesion.", "6"),
					huh.NewOption("Led initiatives that maintained an email-sending error rate at approximately 0.0001%.", "7"),
					huh.NewOption("Played a key role in hiring and developing senior engineering staff: managers and ICs.", "8"),
					huh.NewOption("Helped define the initial job level/ladder for software engineers from Entry- to Staff Engineer.", "9"),
					huh.NewOption("Led the promotion plan for the first Staff Software Engineer at the company.", "10"),
					huh.NewOption("Negotiated 20%% cost-reduction for the company’s most expensive third-party vendor.", "11"),
				).
				Title("ConvertKit").
				Value(&ck),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
