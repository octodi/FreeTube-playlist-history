package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
    "os"
    "path/filepath"
    "regexp"
    "time"

    "github.com/PuerkitoBio/goquery"
    "github.com/schollz/progressbar/v3"
)

type Video struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Duration    string `json:"duration"`
}

func main() {
    fmt.Print("Enter the Invidious API url: ")
    var invidious_api string
    fmt.Scanln(&invidious_api)
    invidious_api += "api/v1/videos/"

    var input_file string
    fmt.Print("Enter the input file path: ")
    fmt.Scanln(&input_file)
    var video_ids []string

    if filepath.Ext(input_file) == ".json" {
        jsonFile, err := os.Open(input_file)
    if err != nil {
        panic(err)
    }
    defer jsonFile.Close()

    var items []map[string]interface{}
    jsonParser := json.NewDecoder(jsonFile)
    if err = jsonParser.Decode(&items); err != nil {
        panic(err)
    }

    re := regexp.MustCompile(`v=(\S{11})`)
    for _, item := range items {
        titleUrl, ok := item["titleUrl"].(string)
        if !ok {
            continue // skip JSON object where titleUrl key is not found
        }
        match := re.FindStringSubmatch(titleUrl)
        if len(match) < 2 {
            continue // skip JSON object where regular expression v=(\S{11}) is not matched
        }
        video_ids = append(video_ids, match[1])
    }

    } else {
        // read contents of HTML file
        htmlContent, err := ioutil.ReadFile(input_file)
        if err != nil {
            panic(err)
        }

        doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(htmlContent)))
        if err != nil {
            panic(err)
        }

        doc.Find("a").Each(func(i int, s *goquery.Selection) {
            href, ok := s.Attr("href")
            if ok {
                re := regexp.MustCompile(`watch\?v=([a-zA-Z0-9_-]{11})`)
                match := re.FindStringSubmatch(href)
                if len(match) >= 2 {
                    video_ids = append(video_ids, match[1])
                }
            }
        })
    }

    bar := progressbar.New64(int64(len(video_ids)))

    file, err := os.OpenFile("watch-history.db", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()

    for i := 0; i < len(video_ids); i++ {
        response, err := http.Get(invidious_api + video_ids[i])
        if err != nil {
            fmt.Println(err)
            continue
        }
        defer response.Body.Close()
        if response.StatusCode != 200 {
            fmt.Printf("Failed to fetch video data for video ID %s. Status code: %d\n", video_ids[i], response.StatusCode)
            fmt.Println(response.Body)
            continue
        }

        jsonData, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Println(err)
            continue
        }

        m := make(map[string]interface{})
        err = json.Unmarshal([]byte(jsonData), &m)
        if err != nil {
            fmt.Println(err)
            continue
        }

        video := make(map[string]interface{})
        video["videoId"] = m["videoId"]
        video["title"] = m["title"]
        video["author"] = m["author"]
        video["authorId"] = m["authorId"]
        video["published"] = m["published"]
        video["description"] = m["description"]
        video["viewCount"] = m["viewCount"]
        video["lengthSeconds"] = m["lengthSeconds"]
        video["watchProgress"] = float64(0)
        video["timeWatched"] = int32(time.Now().UTC().UnixNano() / 1000000)
        video["isLive"] = false
        video["paid"] = false
        video["type"] = "video"

        dataBytes, err := json.Marshal(video)
        if err != nil {
            fmt.Println(err)
            continue
        }

        if _, err := file.Write(append(dataBytes, '\n')); err != nil {
            fmt.Println(err)
            continue
        }
        bar.Add(1)
    }

    fmt.Println("\nDone!")
}
