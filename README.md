# webwatch - command line util to check website changes

./webwatch load list of urls and css selectors from urls.txt compare content with webwatch.db and output changes.

You can setup notification with own tools, I do that with notify_me shell script ```./webwatch | ./notify_me```
```bash
#!/bin/sh

while read x; do 
  if [ -n "$x" ]; then # ignore blank
		curl -i -X GET "https://api.telegram.org/BOTID:TOKEN/sendMessage?chat_id=CHAT_ID&text=$x"
  fi
done

```


# Alternatives
https://github.com/thp/urlwatch

https://github.com/JNRowe/cupage
