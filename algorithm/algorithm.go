package algorithm

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/soloviev1d/cliautotestsuite/config"
	"gopkg.in/yaml.v3"
)

func readTestSuiteDir() ([]string, error) {
	entries, err := os.ReadDir(config.TestSuiteDir)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать содержимое директории: %s", err)
	}

	targets := []string{}
	for _, e := range entries {
		targets = append(targets,
			path.Join(config.TestSuiteDir, e.Name()))
	}

	return targets, nil
}

type TestSuite struct {
	Name        string     `yaml:"name"`
	Author      string     `yaml:"author"`
	Description string     `yaml:"description"`
	TestCases   []TestCase `yaml:"test_cases"`
}

type TestCase struct {
	Type           string   `yaml:"type"`
	Name           string   `yaml:"name"`
	Author         string   `yaml:"author"`
	Description    string   `yaml:"description"`
	CMD            string   `yaml:"cmd"`
	CMDFlags       []string `yaml:"cmd_flags"`
	REPLSteps      []string `yaml:"repl_steps"`
	ExpectedResult string   `yaml:"expected_result"`
}

func serializeTestSuites() ([]*TestSuite, error) {
	testSuites := []*TestSuite{}
	targets, err := readTestSuiteDir()
	if err != nil {
		return nil, err
	}

	for _, t := range targets {
		b, err := os.ReadFile(t)
		if err != nil {
			return nil, fmt.Errorf("не удалось считать файл: %s", err)
		}

		tmpSuite := TestSuite{}
		err = yaml.Unmarshal(b, &tmpSuite)
		if err != nil {
			return nil, fmt.Errorf("неправильный формат тест-сьюта: %s", err)
		}

		testSuites = append(testSuites, &tmpSuite)
	}

	return testSuites, nil
}

func RunTestSuites() (string, error) {
	var (
		sb bytes.Buffer
	)

	suites, err := serializeTestSuites()
	if err != nil {
		return "", err
	}

	for i, s := range suites {
		sb.WriteString(
			fmt.Sprintf("# TS%d \"%s\" \n", i+1, s.Name))
		sb.WriteString(
			fmt.Sprintf("**Автор:** %s \\\n", s.Author),
		)

		sb.WriteString(
			fmt.Sprintf("**Описание:** %s \\\n", s.Description),
		)
		var successfulTCs int = 0

		for j, tc := range s.TestCases {
			sb.WriteString(
				fmt.Sprintf("## TC%d.%d \"%s\" \n", i+1, j+1, tc.Name))
			sb.WriteString(
				fmt.Sprintf("**ID:** TC%d.%d-%s \\\n", i+1, j+1, tc.Type),
			)
			sb.WriteString(
				fmt.Sprintf("**Автор:** %s \\\n", tc.Author),
			)
			sb.WriteString(
				fmt.Sprintf("**Описание:** %s \\\n", tc.Description),
			)
			t := time.Now()
			sb.WriteString(
				fmt.Sprintf("**Дата проведения:** %d/%d/%d %d:%d \\\n",
					t.Day(),
					t.Month(),
					t.Year(),
					t.Hour(),
					t.Minute(),
				),
			)
			sb.WriteString(
				fmt.Sprintf("**Ожидаемый результат:** `%s` \\\n", tc.ExpectedResult),
			)

			if tc.CMD == "" {
				return "", fmt.Errorf("поле cmd не может быть пустым: %s", err)
			}
			cmd := exec.Command("bash", "-c", tc.CMD)
			b, _ := cmd.Output()
			rs := strings.Trim(string(b), "\n")

			sb.WriteString(
				fmt.Sprintf("**Реальный результат:** `%s` \\\n", rs),
			)
			if rs == tc.ExpectedResult {
				sb.WriteString(
					fmt.Sprint("**Итог:** Тест пройден ✅ \\\n"),
				)
				successfulTCs++
			} else if tc.ExpectedResult == "<NOT NULL>" {
				sb.WriteString(
					fmt.Sprint("**Итог:** Тест пройден ✅ \\\n"),
				)
				successfulTCs++
			} else {
				sb.WriteString(
					fmt.Sprint("**Итог:** Тест провален ❌ \\\n"),
				)
			}
		}
		sb.WriteString(
			fmt.Sprintf("**Итого:** %d/%d \n", successfulTCs, len(s.TestCases)),
		)
	}

	return sb.String(), nil
}
