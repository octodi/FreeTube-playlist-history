package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "strings"
    "net/http"
    "os"
    "strconv"
    "sync"
    "time"

    "github.com/schollz/progressbar/v3"
)

func main() {
    fmt.Print("Enter the Invidious API url: ")
    var invidiousAPI string
    fmt.Scanln(&invidiousAPI)
    fmt.Print("Enter the CSV file path: ")
    var csvPath string
    fmt.Scanln(&csvPath)

    videoIDs, err := getVideoIDsFromCSV(csvPath)
    if err != nil {
        log.Fatal(err)
    }

    bar := progressbar.New64(int64(len(videoIDs)))

    videoDataCh := make(chan map[string]interface{})
    var wg sync.WaitGroup
    for _, videoID := range videoIDs {
        wg.Add(1)
        go func(id string) {
            defer wg.Done()
            u := invidiousAPI + "api/v1/videos/" + id
            resp, err := http.Get(u)
            if err != nil {
                log.Printf("Failed to fetch video data for video ID %s. Error: %v\n", id, err)
                return
            }
            defer resp.Body.Close()

            if resp.StatusCode == 200 {
                bodyBytes, err := ioutil.ReadAll(resp.Body)
                if err != nil {
                    log.Printf("Failed to read response for video ID %s. Error: %v\n", id, err)
                    return
                }

                var bodyMap map[string]interface{}
                err = json.Unmarshal(bodyBytes, &bodyMap)
                if err != nil {
                    log.Printf("Failed to unmarshall JSON for video ID %s. Error: %v\n", id, err)
                    return
                }

                video := make(map[string]interface{})
                video["videoId"] = bodyMap["videoId"]
                video["title"] = bodyMap["title"]
                video["author"] = bodyMap["author"]
                video["authorId"] = bodyMap["authorId"]
                video["published"] = ""
                video["description"] = ""
                video["viewCount"] = bodyMap["viewCount"]
                video["lengthSeconds"] = bodyMap["lengthSeconds"]
                video["timeAdded"] = strconv.FormatInt(time.Now().UnixNano(), 10)
                video["isLive"] = false
                video["paid"] = false
                video["type"] = "video"
                video["channel_url"] = invidiousAPI + "/user/" + bodyMap["author"].(string)
                video["thumbnail_url"] = invidiousAPI + "/vi/" + bodyMap["videoId"].(string) + "/maxresdefault.jpg"
                video["external_player"] = invidiousAPI + "/watch?v=" + bodyMap["videoId"].(string)
                videoDataCh <- video
            } else {
                log.Printf("Failed to fetch video data for video ID %s. Status code: %d\n", id, resp.StatusCode)
            }
        }(videoID)
    }
    
    go func() {
        wg.Wait()
        close(videoDataCh)
    }()

    videoData := make([]map[string]interface{}, 0, len(videoIDs))
    for video := range videoDataCh {
        videoData = append(videoData, video)
        bar.Add(1)
    }

    playlistData := map[string]interface{}{
        "playlistName": "Favorites",
        "videos":       videoData,
    }

    file, _ := json.MarshalIndent([]map[string]interface{}{playlistData}, "", " ")
    _ = ioutil.WriteFile("playlist_data.db", file, 0644)
}

func getVideoIDsFromCSV(path string) ([]string, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var ids []string
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        fields := strings.Split(line, ",")
        id := strings.TrimSpace(fields[0])
        ids = append(ids, id)
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return ids, nil
}
