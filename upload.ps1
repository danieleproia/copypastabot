scp .\data.json eleven:/root/scripts/copypastabot/
# launch command via ssh
ssh eleven "supervisorctl restart copypastabot"