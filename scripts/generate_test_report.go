package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

type TestInfo struct {
	Module      string
	TestType    string // Unit / Integration
	TestName    string
	Description string
}

func main() {
	var tests []TestInfo

	// Modules to check
	modules := []string{
		"internal/modules/course",
		"internal/modules/student",
		"internal/modules/ta_duty",
		"internal/modules/announce",
	}

	for _, mod := range modules {
		baseDir := filepath.Join(mod, "service") // Most tests are in service
		if _, err := os.Stat(baseDir); os.IsNotExist(err) {
			continue
		}

		cmd := exec.Command("go", "test", "-v", "./"+baseDir)
		output, _ := cmd.CombinedOutput() // Ignore error as we still want to parse the output if tests fail

		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if !strings.HasPrefix(line, "=== RUN   ") {
				continue
			}

			testName := strings.TrimSpace(strings.TrimPrefix(line, "=== RUN   "))

			// Clean module name
			modName := strings.TrimPrefix(mod, "internal/modules/")

			// Determine type
			testType := "Unit Test"
			if strings.Contains(testName, "_Integration") {
				testType = "Integration Test"
			}

			// Generate simple description based on naming convention
			desc := "Verifies " + testName

			// Handle subtests (which contain '/')
			if strings.Contains(testName, "/") {
				nameParts := strings.Split(testName, "/")
				parentName := nameParts[0]
				subName := nameParts[len(nameParts)-1]

				// Try to strip T001_, IT001_ prefixes for cleaner descriptions
				cleanSub := subName
				subParts := strings.SplitN(subName, "_", 2)
				if len(subParts) == 2 && (strings.HasPrefix(subParts[0], "T") || strings.HasPrefix(subParts[0], "IT") || strings.HasPrefix(subParts[0], "TC")) {
					cleanSub = strings.ReplaceAll(subParts[1], "_", " ")
				} else {
					cleanSub = strings.ReplaceAll(subName, "_", " ")
				}
				desc = "Tests " + parentName + " -> " + cleanSub
			} else {
				// Top-level test
				parts := strings.Split(testName, "_")
				if len(parts) > 1 {
					desc = "Tests " + parts[0] + " (" + strings.Join(parts[1:], " ") + ")"
				} else {
					desc = "Tests " + testName
				}
			}

			tests = append(tests, TestInfo{
				Module:      modName,
				TestType:    testType,
				TestName:    testName,
				Description: desc,
			})
		}
	}

	// Create excel
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Sheet setup
	sheet := "Sheet1"
	f.SetCellValue(sheet, "A1", "Module")
	f.SetCellValue(sheet, "B1", "Test Type")
	f.SetCellValue(sheet, "C1", "Test Name")
	f.SetCellValue(sheet, "D1", "Description")

	// Make Header bold
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	f.SetCellStyle(sheet, "A1", "D1", style)

	// Add data
	for i, t := range tests {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), t.Module)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), t.TestType)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), t.TestName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), t.Description)
	}

	// Adjust column widths
	f.SetColWidth(sheet, "A", "A", 15)
	f.SetColWidth(sheet, "B", "B", 18)
	f.SetColWidth(sheet, "C", "C", 40)
	f.SetColWidth(sheet, "D", "D", 60)

	// Save
	outputPath := "Test_Report.xlsx"
	if err := f.SaveAs(outputPath); err != nil {
		fmt.Printf("Error saving excel file: %v\n", err)
		return
	}

	fmt.Printf("Successfully generated %s with %d tests.\n", outputPath, len(tests))
}
