# FreeTube-playlist-history
<p>This program can help you import playlist and history from your YouTube's takeout. It grab's videoIds from your csv, html or JSON file and fetch other data required from Invidious API and outputs them as JSON data and saves in database file, so it could work with FreeTube.<p1>

<h3>Prerequisite for Python<h3>

```bash
pip install httpx tqdm
```

<h3>Prerequisite for Go<h3>
<h4> Install Go on your system if running from source</h4>
<p> Running from go source will install dependencies itself <br> Through go binary just give it executable permission<p1>

```bash
chmod +x bin_name
```
<h4>Note :- I will prefer you to use Go [binary or source] as it will function comparatively faster <br> Grab your Invidious API from here : https://api.invidious.io/ <br> You could get more 503 with go-history/faster if the file is large </h4>

<h3> Remove Starting lines from CSV </h3>
<h4> For smooth functioning remove starting lines from the csv file </h4>
<p> Like this :- <p1>

![alt text](https://github.com/octodi/FreeTube-playlist-history/blob/main/img/Untitled.png)

<h3> Some benchmarks <h3>

![alt text](https://github.com/octodi/FreeTube-playlist-history/blob/main/img/1.png)
![alt text](https://github.com/octodi/FreeTube-playlist-history/blob/main/img/2.png)

# Ditch Google
# Suggest me a better repo name !please




