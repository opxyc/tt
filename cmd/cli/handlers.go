package main

import (
	"bufio"
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

	_, err := ttService.Start(title, description, tags)
	if err != nil {
		fmt.Printf("failed to start '%s': %s\n", title, err)
		os.Exit(1)
	}

	fmt.Println("noted ..")
	return nil
}

func pause(_ *cli.Context) error {
	activitiesCurrentlyInProgress, err := ttService.List(&tt.ListFilters{Status: []tt.ActivityStatus{tt.StatusInProgress}})
	if err != nil {
		fmt.Printf("failed to pause: %s\n", err)
		os.Exit(1)
	}

	if len(activitiesCurrentlyInProgress) < 1 {
		fmt.Println("no activities in progress")
		os.Exit(1)
	}

	activityCurrentlyInProgress := activitiesCurrentlyInProgress[0]

	input := ""
	fmt.Printf("are you sure to pause '%s'? (Y/n) ", activityCurrentlyInProgress.Title)
	fmt.Scanln(&input)
	input = strings.ToLower(input)

	if input == "n" {
		os.Exit(0)
	}

	_, err = ttService.Pause(activityCurrentlyInProgress.ID)
	if err != nil {
		fmt.Printf("failed to pause '%s': %s\n", activityCurrentlyInProgress.Title, err)
		os.Exit(1)
	}

	fmt.Println("paused")
	return nil
}

func resume(_ *cli.Context) error {
	pausedActivities, err := ttService.List(&tt.ListFilters{Status: []tt.ActivityStatus{tt.StatusPaused}})
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
	filterDate := cCtx.String("date")
	filters := &tt.ListFilters{Date: time.Now().Format("2006-01-02")}
	if filterDate == "all" {
		filters.Date = ""
	} else if len(filterDate) == len("2006-10-11") {
		filters.Date = filterDate
	}

	activitesList, err := ttService.List(filters)
	if err != nil {
		return err
	}

	if len(activitesList) == 0 {
		fmt.Println("nothing to list")
		return nil
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"DT", "Title", "Description", "Tags", "Status", "Duration"})
	var rows []table.Row
	for _, activity := range activitesList {
		rows = append(rows, table.Row{activity.CreatedAt.Format("2006-01-02 03:04:05 PM"), activity.Title, activity.Desc, activity.Tags, activity.Status, activity.Duration})
	}
	t.AppendRows(rows)
	t.Render()

	return nil
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
