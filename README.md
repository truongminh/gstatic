# gstatic
Static file server 

## Serving static file

```toml
[[Static]]
Route="/files"
Folder="data/files"
    [Static.Headers]
    Cache-Control="public,maxage=60,no-cache"
```

## Transform image
- /files/<filename>?transform=resize:600x400
Credit for image lib is https://github.com/nfnt/resize

## Proxy
- /proxy?link=<target>
