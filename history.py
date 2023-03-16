import csv
import requests
from datetime import datetime
from tqdm import tqdm
import json
import re
import os

invidious_api = input('Enter the Invidious API url: ') + 'api/v1/videos/'
input_file = input('Enter the input file path: ')

# check if input file is a JSON file
if os.path.splitext(input_file)[1] == '.json':
    with open(input_file) as f:
        video_ids = [re.search(r'v=(\S{11}?)(&|$)', item.get('titleUrl', '')).group(1) if re.search(r'v=(\S{11}?)(&|$)', item.get('titleUrl', '')) else None for item in json.load(f) if 'titleUrl' in item and item['titleUrl'] is not None]



else:
    with open(input_file, "r") as f:
        html_content = f.read()
        video_ids = re.findall(r'watch\?v=(\S{11})', html_content)

video_data = []
for video_id in tqdm(video_ids, desc='Fetching video details', unit='videos'):
    response = requests.get(invidious_api + video_id)
    if response.status_code == 200:
        video = response.json()
        video_data.append({"videoId": video["videoId"],
                           "title": video["title"],
                           "author": video["author"],
                           "authorId": video["authorId"],
                           "published": video["published"],
                           "description": video["description"],
                           "viewCount": video["viewCount"],
                           "lengthSeconds": video["lengthSeconds"],
                           "watchProgress": 0,
                           "timeWatched": int(datetime.now().replace(second=0, microsecond=0).timestamp()) * 1000,
                           "isLive": False,
                           "paid": False,
                           "type": "video"})
    else:
        print(f'Failed to fetch video data for video ID {video_id}. Status code: {response.status_code}')
        print(response.text)

playlist_data = video_data
with open('watch-history.db', 'w') as outfile:
    for video in playlist_data:
        outfile.write(json.dumps(video))
        outfile.write('\n')
