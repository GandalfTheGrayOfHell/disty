import requests
import os


FILEPATH = "./test/src/main.c"

file = open(FILEPATH)
data = file.read()

mtime = os.path.getmtime(FILEPATH)

r = requests.get("http://localhost:3000/push?project=test&filename=/src/main.c&track=TRACKED&modtime=" + str(int(mtime)), data=data)
file.close()

print("/push | Status:", r.status_code)

exit(0)