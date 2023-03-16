import csv
import requests
from datetime import datetime
from tqdm import tqdm
import json

invidious_api = input('Enter the Invidious API url: ')+'api/v1/videos/'
csv_file=input('Enter the csv file path: ')

with open(csv_file, 'r') as csvfile:
  reader = csv.reader(csvfile)
  video_ids=[]
  next(reader)
  for row in reader:
    if len(row)>0:
      v_id=row[0].strip()
      video_ids.append(v_id)
  video_data = []
  for video_id in tqdm(video_ids, desc='Fetching video details', unit='videos'):
    response = requests.get(invidious_api + video_id)
    if response.status_code == 200:
      video = response.json()
      video_data.append({"videoId":video["videoId"],
      "title":video["title"],
      "author":video["author"],
      "authorId":video["authorId"],
      "published":"",
      "description":"",
      "viewCount":video["viewCount"],
      "lengthSeconds":video["lengthSeconds"],
      "timeAdded":int(datetime.now().replace(second=0, microsecond=0).timestamp()) * 1000,
      "isLive":False,
      "paid":False,
      "type":"video"})
    else:
      print(f'Failed to fetch video data for video ID {video_id}. Status code: {response.status_code}')
      print(response.text)

  playlist_data = {"playlistName":"Favorites",
  "videos":video_data}
  with open('playlist_data.db', 'w') as outfile:
    outfile.write(json.dumps([playlist_data]))
