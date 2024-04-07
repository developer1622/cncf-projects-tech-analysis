package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"context"
	"os"

	"github.com/google/go-github/github"
	"github.com/wcharczuk/go-chart"
	"gopkg.in/yaml.v2"

	"golang.org/x/oauth2"
)

var (
	// you need to generate personal access token at
	// https://github.com/settings/applications#personal-access-tokens
	personalAccessToken = "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"
 )
 
 
 type TokenSource struct {
	AccessToken string
 }
 
 
func writeFile(filename string, data []byte) error {
	// Create or open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
 
 
	// Write the JSON data to the file
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

 func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
 }

func main() {
	
	// TODO: make it modularize
	// should be in seperate file, let's keep it in
	// initial version
	tokenSource := &TokenSource{
		AccessToken: personalAccessToken,
	}

	ctx := context.Background()
	oauthClient := oauth2.NewClient(ctx, tokenSource)
	client := github.NewClient(oauthClient)

	// Read the YAML file taken from
	// https://github.com/cncf/devstats/blob/master/projects.yaml
	yamlFile, err := os.ReadFile("projects.yaml")
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}

	// Parse YAML into a map for processing
	var projectsMap map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &projectsMap)
	if err != nil {
		fmt.Printf("Error unmarshaling YAML data: %v\n", err)
		return
	}

	// Extract "main_repo" from each project
	// assuming that we would have field called
	// main_repo
	mainRepos := make([]string, 0)
	projects, ok := projectsMap["projects"].(map[interface{}]interface{})
	if !ok {
		fmt.Println("Error parsing projects data")
		return
	}

	for _, project := range projects {
		projectMap, ok := project.(map[interface{}]interface{})
		if !ok {
			continue
		}
		mainRepo, ok := projectMap["main_repo"].(string)
		if ok {
			mainRepos = append(mainRepos, mainRepo)
		}
	}

	// Print the main repositories
	fmt.Println("main Repositories:")
	languagesList:=make([]map[string]int,0)
	for _, repo := range mainRepos {
		ownerRepo := strings.Split(repo, "/")
       if len(ownerRepo) == 1 || len(ownerRepo) > 3 {
           fmt.Printf("***** seems like some issue with ****** : %s\n", repo)
           continue
       }

	   // not sure about this API rate limit for normal user
	   // but make sure we do not run into it
       time.Sleep(time.Second * 2)


	   // TODO: do not ignore error
       langList, _, err := client.Repositories.ListLanguages(ctx, ownerRepo[0], ownerRepo[1])
       if err != nil {
           fmt.Printf("failed to fetch languages for: %s\n", repo)
       }
	   // debugging 
	   fmt.Println(langList)
       languagesList = append(languagesList, langList)
	}

// we do not want to run all the time, let's save this to the file
	// Marshal the mapList into JSON
	jsonData, err := json.MarshalIndent(languagesList, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling data to JSON:", err)
		return
	}

	// Write the JSON data to a file
	err = os.WriteFile("projects_cncf.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON data to file:", err)
		return
	}

	fmt.Println("Data written to 'projects_cncf.json' file successfully")


	// Initialize a map to store counts for each programming language
	fmt.Println("we got list of all the project languages, let build chart")
	languageCounts := make(map[string]int)
	for _, item := range languagesList {
		for lang, count := range item {
			// Increment the count for the current language
			languageCounts[lang] += count
		}
	}

	// Create data for the pie chart
	var values []chart.Value
	for lang, count := range languageCounts {
		values = append(values, chart.Value{Label: lang, Value: float64(count)})
	}

	// Create a pie chart
	pie := chart.PieChart{
		Width:  4096,  // Increased width
		Height: 4096,  // Increased height
		Values: values,
	}

	// Save the chart as an image file
	f, err := os.Create("pie_chart.png")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	err = pie.Render(chart.PNG, f)
	if err != nil {
		fmt.Println("Error rendering pie chart:", err)
		return
	}
	fmt.Println("Pie chart saved as 'pie_chart.png'")
}


