package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/opxyc/tt/pkg/tt"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

func start() error {
	activitiesInProgress, err := ttService.List(&tt.ListFilters{Status: []tt.ActivityStatus{tt.StatusInProgress}})
	if err != nil {
		fmt.Printf("failed to start: %s\n", err)
		os.Exit(1)
	}

	if len(activitiesInProgress) != 0 {
		err = pauseActivity(&activitiesInProgress[0])
		if err != nil {
			fmt.Printf("failed to pause '%s': %s\n", activitiesInProgress[0].Title, err)
			os.Exit(1)
		}

		fmt.Println("new activity:")
	}

	buffer := bufio.NewReader(os.Stdin)

	fmt.Print("> title ")
	title, _ := buffer.ReadString('\n')
	title = strings.TrimSpace(strings.TrimSuffix(title, "\n"))

	fmt.Print("> description ")
	description, _ := buffer.ReadString('\n')
	description = strings.TrimSpace(strings.TrimSuffix(description, "\n"))

	fmt.Print("> tags ")
	tags, _ := buffer.ReadString('\n')
	tags = strings.TrimSpace(strings.TrimSuffix(tags, "\n"))

	_, err = ttService.Start(title, description, tags)
	if err != nil {
		fmt.Printf("failed to start '%s': %s\n", title, err)
		os.Exit(1)
	}

	fmt.Println("noted ..")
	return nil
}

func pause(_ *cli.Context) error {
	// there will be only one activity currently in progress since we do not allow to start another
	// activity when one is in progress
	activitiesInProgress, err := ttService.List(&tt.ListFilters{Status: []tt.ActivityStatus{tt.StatusInProgress}})
	if err != nil {
		fmt.Printf("failed to pause: %s\n", err)
		os.Exit(1)
	}

	if len(activitiesInProgress) < 1 {
		fmt.Println("no activities in progress")
		os.Exit(1)
	}

	err = pauseActivity(&activitiesInProgress[0])
	if err != nil {
		fmt.Printf("failed to pause '%s': %s\n", activitiesInProgress[0].Title, err)
		os.Exit(1)
	}

	fmt.Println("paused")
	return nil
}

func pauseActivity(activityInProgress *tt.Activity) error {
	input := ""
	fmt.Printf("are you sure to pause '%s'? (Y/n) ", activityInProgress.Title)
	fmt.Scanln(&input)
	input = strings.ToLower(input)

	if input == "n" {
		os.Exit(0)
	}

	_, err := ttService.Pause(activityInProgress.ID)
	return err
}

func resume(_ *cli.Context) error {
	pausedActivities, err := ttService.List(&tt.ListFilters{Status: []tt.ActivityStatus{tt.StatusPaused}})
	if err != nil {
		fmt.Printf("failed to resume: %s\n", err)
		os.Exit(1)
	}

	activitiesInProgress, err := ttService.List(&tt.ListFilters{Status: []tt.ActivityStatus{tt.StatusInProgress}})
	if err != nil {
		fmt.Printf("failed to resume: %s\n", err)
		os.Exit(1)
	}

	totalPausedActivities := len(pausedActivities)
	if totalPausedActivities < 1 {
		fmt.Println("no activities to resume")
		os.Exit(1)
	}

	for i, activity := range pausedActivities {
		fmt.Printf("%d. %s\n", i+1, activity.Title)
	}

	input := ""
	fmt.Printf("> ")
	fmt.Scanln(&input)
	indexSelected, err := strconv.Atoi(input)
	indexSelected -= 1
	if err != nil || indexSelected > totalPausedActivities-1 {
		fmt.Println("invalid input")
		os.Exit(1)
	}

	if len(activitiesInProgress) != 0 {
		err = pauseActivity(&activitiesInProgress[0])
		if err != nil {
			fmt.Printf("failed to resume: %s\n", err)
			os.Exit(1)
		}
	}

	_, err = ttService.Resume(pausedActivities[indexSelected].ID)
	if err != nil {
		fmt.Printf("failed to resume '%s': %s\n", pausedActivities[indexSelected].Title, err)
		os.Exit(1)
	}

	fmt.Println("resumed ..")

	return nil
}

func stop(_ *cli.Context) error {
	stoppableActivities, err := ttService.List(&tt.ListFilters{Status: []tt.ActivityStatus{tt.StatusInProgress, tt.StatusPaused}})
	if err != nil {
		fmt.Printf("failed to stop: %s\n", err)
		os.Exit(1)
	}

	totalStoppableActivities := len(stoppableActivities)
	if totalStoppableActivities < 1 {
		fmt.Println("no activities to stop")
		os.Exit(1)
	}

	i := 0
	for _, activity := range stoppableActivities {
		fmt.Printf("%d. %s (%s)\n", i+1, activity.Title, strings.ToLower(string(activity.Status)))
		i++
	}

	input := ""
	fmt.Printf("> ")
	fmt.Scanln(&input)
	indexSelected, err := strconv.Atoi(input)
	indexSelected -= 1
	if err != nil || indexSelected > totalStoppableActivities-1 {
		fmt.Println("invalid input")
		os.Exit(1)
	}

	activityToStop := stoppableActivities[indexSelected]
	activityDetails, err := ttService.Stop(activityToStop.ID)
	if err != nil {
		fmt.Printf("failed to stop '%s': %s\n", activityToStop.Title, err)
		os.Exit(1)
	}

	fmt.Printf("stopped. total duration: %v\n", activityDetails.Duration)

	return nil
}

func list(cCtx *cli.Context) error {
	filterDateFlagValue := cCtx.String("date")
	printAsCsvFlagValue := cCtx.Bool("csv")
	filters := &tt.ListFilters{Date: time.Now().Format("2006-01-02")}
	if filterDateFlagValue == "all" {
		filters.Date = ""
	} else if len(filterDateFlagValue) == len("2006-10-11") {
		filters.Date = filterDateFlagValue
	}

	activityList, err := ttService.List(filters)
	if err != nil {
		return err
	}

	if len(activityList) == 0 {
		fmt.Println("nothing to list")
		return nil
	}

	if printAsCsvFlagValue {
		return printAsCsv(&activityList)
	}

	printAsTable(&activityList)
	return nil
}

func printAsCsv(activityList *[]tt.Activity) error {
	w := csv.NewWriter(os.Stdout)
	var rows [][]string
	for _, activity := range *activityList {
		rows = append(rows, []string{activity.CreatedAt.Format("2006-01-02 03:04:05 PM"), activity.Title, activity.Desc, activity.Tags, string(activity.Status), fmt.Sprintf("%f", activity.Duration)})
	}

	return w.WriteAll(rows)
}

func printAsTable(activityList *[]tt.Activity) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"DT", "Title", "Description", "Tags", "Status", "Duration"})
	var rows []table.Row
	for _, activity := range *activityList {
		rows = append(rows, table.Row{activity.CreatedAt.Format("2006-01-02 03:04:05 PM"), activity.Title, activity.Desc, activity.Tags, activity.Status, activity.Duration})
	}
	t.AppendRows(rows)
	t.SetStyle(table.StyleRounded)
	t.Render()
}

func delete(_ *cli.Context) error {
	deletableActivities, err := ttService.List(&tt.ListFilters{Status: []tt.ActivityStatus{tt.StatusInProgress, tt.StatusPaused}})
	if err != nil {
		fmt.Printf("failed to delete: %s\n", err)
		os.Exit(1)
	}

	totalDeletableActivities := len(deletableActivities)
	if totalDeletableActivities < 1 {
		fmt.Println("no running or paused activities to delete")
		os.Exit(1)
	}

	i := 0
	for _, activity := range deletableActivities {
		fmt.Printf("%d. %s (%s)\n", i+1, activity.Title, strings.ToLower(string(activity.Status)))
		i++
	}

	input := ""
	fmt.Printf("> ")
	fmt.Scanln(&input)
	indexSelected, err := strconv.Atoi(input)
	indexSelected -= 1
	if err != nil || indexSelected > totalDeletableActivities-1 {
		fmt.Println("failed to delete: invalid input")
		os.Exit(1)
	}

	activityToDelete := deletableActivities[indexSelected]
	_, err = ttService.Delete(activityToDelete.ID)
	if err != nil {
		fmt.Printf("failed to delete '%s': %s\n", activityToDelete.Title, err)
		os.Exit(1)
	}

	fmt.Println("deleted")

	return nil
}
