package main

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"
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

    var videoData []map[string]interface{}
    for _, videoID := range videoIDs {
        u := invidiousAPI + "api/v1/videos/" + videoID
        resp, err := http.Get(u)
        if err != nil {
            log.Fatal(err)
        }
        defer resp.Body.Close()

        if resp.StatusCode == 200 {
            bodyBytes, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                log.Fatal(err)
            }

            var bodyMap map[string]interface{}
            err = json.Unmarshal(bodyBytes, &bodyMap)
            if err != nil {
                log.Fatal(err)
            }

            video := map[string]interface{}{
                "videoId":         bodyMap["videoId"],
                "title":           bodyMap["title"],
                "author":          bodyMap["author"],
                "authorId":        bodyMap["authorId"],
                "published":       "",
                "description":     "",
                "viewCount":       bodyMap["viewCount"],
                "lengthSeconds":   bodyMap["lengthSeconds"],
                "timeAdded":       strconv.FormatInt(time.Now().UnixNano(), 10),
                "isLive":          false,
                "paid":            false,
                "type":            "video",
                "channel_url":     fmt.Sprintf("%s/user/%s", invidiousAPI, bodyMap["author"]),
                "thumbnail_url":   fmt.Sprintf("%s/vi/%s/maxresdefault.jpg", invidiousAPI, bodyMap["videoId"]),
                "external_player": fmt.Sprintf("%s/watch?v=%s", invidiousAPI, bodyMap["videoId"]),
            }

            videoData = append(videoData, video)
        } else {
            fmt.Printf("Failed to fetch video data for video ID %s. Status code: %d\n", videoID, resp.StatusCode)
        }
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

    r := csv.NewReader(f)
    r.Read()

    var ids []string
    for {
        record, err := r.Read()
        if err != nil {
            break
        }
        ids = append(ids, record[0])
    }
    return ids, nil
}