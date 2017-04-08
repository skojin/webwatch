# webwatch - command line util to check website changes

./webwatch load list of URLs and CSS selectors from urls.txt, check for updates and output if have page was changed.

urls.txt is plain text file with following format
```text
  comment starts with space, empty lines ignored

http://target_url
#filter .as.css.selector

http://www.amazon.com/dp/B00GDQ0RMG/ref=ods_gw_d_h1_s
#priceblock_ourprice

  filter started with > is shell command with page body as input, useful with https://github.com/ericchiang/pup
https://news.ycombinator.com
> sort | uniq | wc -l

https://news.ycombinator.com
> pup 'table table tr:nth-last-of-type(n+2) td.title a attr{href}'
```

# Notifications

You can setup notification with your own tools, I do that with notify_me shell script ```./webwatch | ./notify_me```
```bash
#!/bin/sh

while read x; do 
  if [ -n "$x" ]; then # ignore blank
    curl -i -X GET "https://api.telegram.org/BOTID:TOKEN/sendMessage" -F "chat_id=CHAT_ID" -F "text=$x"
  fi
done
```

# Download

Download linux and OSX [binaries](https://github.com/skojin/webwatch/releases)

# Alternatives
https://github.com/thp/urlwatch

https://github.com/JNRowe/cupage
