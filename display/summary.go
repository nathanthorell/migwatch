package display

import (
	"fmt"
	"time"

	"github.com/nathanthorell/migwatch/model"
)

func PrintSummary(result model.EnvironmentResult) {
	if result.Error != nil {
		fmt.Println(styleFailed.Render("  error: " + result.Error.Error()))
		fmt.Println()
		return
	}

	if len(result.Migrations) == 0 {
		fmt.Println(styleDim.Render("  no migrations found"))
		fmt.Println()
		return
	}

	var lastSuccess *model.Migration
	versionedCount := 0
	repeatableCount := 0
	var failures []model.Migration

	for i := range result.Migrations {
		m := &result.Migrations[i]
		if m.Version == "" {
			repeatableCount++
		} else {
			versionedCount++
		}
		if !m.Success {
			failures = append(failures, *m)
		}
		if m.Success && m.Version != "" && (lastSuccess == nil || m.InstalledRank > lastSuccess.InstalledRank) {
			lastSuccess = m
		}
	}

	if lastSuccess != nil {
		label := styleHeader.Render("V" + lastSuccess.Version)
		desc := lastSuccess.Description
		date := styleDim.Render(lastSuccess.InstalledOn.Format(time.DateTime))
		fmt.Printf("  Last migration:  %s — %s  %s\n", label, desc, date)
	}

	totalLine := fmt.Sprintf("%d versioned", versionedCount)
	if repeatableCount > 0 {
		totalLine += fmt.Sprintf(", %d repeatable", repeatableCount)
	}
	fmt.Printf("  Total:           %s\n", totalLine)

	if len(failures) == 0 {
		fmt.Printf("  Failures:        %s\n", styleSuccess.Render("none"))
	} else {
		fmt.Printf("  Failures:        %s\n", styleFailed.Render(fmt.Sprintf("%d", len(failures))))
		for _, f := range failures {
			fmt.Printf("    %s  %s — %s\n",
				styleFailed.Render("✗"),
				styleHeader.Render("V"+f.Version),
				f.Description,
			)
		}
	}

	fmt.Println()
}
