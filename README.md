# gstatic
Static file server 

## Serving static file
Given this config:

```toml
[[Static]]
Route="/files"
Folder="data/files"
    [Static.Headers]
    Cache-Control="public,maxage=60,no-cache"
```
Access to `<host>/files/<filename>` mapped to `data/files/<filename>`
with http header `Cache-Control` set to `public,maxage=60,no-cache`

## Transform image
- `<host>/files/<filename>?transform=resize:600x400`

Credit for image lib is https://github.com/nfnt/resize

## Proxy
- `/proxy?link=<target>`

All proxied responses are stored in cache.
