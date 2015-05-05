# picel [![Build Status](http://img.shields.io/travis/henvic/picel/master.svg?style=flat)](https://travis-ci.org/henvic/picel) [![GoDoc](https://godoc.org/github.com/henvic/picel?status.svg)](https://godoc.org/github.com/henvic/picel)

picel is a light-weight, blazing fast REST-ful micro service for image processing with a lean API.

It does one thing and does it well: process images like a UNIX [filter](https://en.wikipedia.org/wiki/Filter_\(software\)) (see [Basics of the UNIX Philosophy](http://www.catb.org/esr/writings/taoup/html/ch01s06.html)).

You can use picel as a middleware between your backend (storage) and your caching layer. User requests can go trough the caching layer and, if the content is not found, it is proxied to the picel microservice, which will then process and try to retrieve the raw data from the backend server, before processing and returning the response. The response SHOULD be cached by the caching layer for further requests.

This way you can simplify your operations environment by avoiding multiple caching layers. Just be aware that due to the very nature of images (large in size, seldom changing) it is wise to consider carefully appropriate caching rules to avoid premature removal.

Also, you want your proxy layer to have protection against abuse ([enhance your calm](http://httpstatusdogs.com/420-enhance-your-calm) to avoid trying to process [too many suspicious requests](http://httpstatusdogs.com/429-too-many-requests)). Refer to [rfc6585#section-4](https://tools.ietf.org/html/rfc6585#section-4) to know more. Modern HTTP servers such as nginx already have options to deal with such attacks.

picel is designed to be safe so that you can try to process untrusted data with no security issues.

## tl;dr
1. Download the picel binary for your platform.
2. Run it with no arguments or with something like `--Backend localhost:8080`
3. Use [picel-js](https://github.com/henvic/picel-js) to encode your images.

The default port is 8123. Change it with `--addr :8000` to listen on port 8000 (for example).

`picel --help` for more help.

## Defaults, performance friendly, and more
By default, picel will try to use webp if the user doesn't explicitly request another format and his client announces it accepts it (Chrome, for example).

Also, JPEG is the default image format for input.

## Dependencies
picel uses [webp](https://developers.google.com/speed/webp/) and [ImageMagick](http://www.imagemagick.org/). At startup it will warn if it doesn't find the binaries for these processes. If you don't have it (or are running old versions) use your operating system package manager system to install the newest versions.

## Protocol
`GET /<backend>/<id><params>.<output>`

* backend: image storage end-point (only when picel server is open / unrestricted)
* id: part of the path before the last '.'
* params: image manipulation parameters / raw
* output: output delivery format

The id and output MUST be escaped by **_** (underscore).

backend is the server backend, to be used if picel was not started with the `--Backend` flag. It should be a host without a leading `http://` or, in case of `https://`, with a `s:` prepended.

Available params are:

1. `raw`
2. `crop {x, y, width, height}`
3. `dimension {width, height}`
4. `extension`

Parameters MUST be given in this order or, otherwise, picel will not recognize them (this is by design on purpose, to avoid having multiple encoding implementations doing things differently / guarantee more cache hits when using a caching layer).

* raw is a parameter without value and MUST NOT be used along others. It implies that picel SHOULD return the original file from the backend. This option might not be available.
* crop MUST be given using the format `<x>x<y>:<width>x<height>` as in `0x0:100x200`
* width and height are pixel integers using the format `<width>x<height>`, when one is neglected the resizing is made proportional
* extension is a string, when it's the same as of the output it is discarded

All parameters are prefixed by a **_** (underscore).

### ?explain
To help debugging you can use append the ?explain to a URL in order to get a JSON response that will tell you how a image was transformed (or failed to be).

Example:
`curl https://localhost:8123/s:example.net/foo_137x0:737x450_800x600_jpg.webp` will return the requested image (if it exists).

`curl https://localhost:8123/s:example.net/foo_137x0:737x450_800x600_jpg.webp?explain` will return an object like the following:

```
{
    "message": "Success. Image path parsed and decoded correctly",
    "transform": {
        "image": {
            "id": "foo",
            "extension": "jpg",
            "source": "https://example.net/foo.jpg"
        },
        "original": false,
        "width": 800,
        "height": 600,
        "crop": {
            "x": 137,
            "y": 0,
            "width": 737,
            "height": 450
        },
        "output": "webp"
    },
    "errors": null
}
```

Please note that ?explain only tell if a request is **not bad** and does **NOT** verify if processing works or even if an image exists on the backend server. If you just need to verify it process correctly you can judge by getting the Content-Length from a 200 OK'ed HEAD request.

## Encoding libraries
Currently you can encode urls using either the go package or using the auxiliary [JavaScript encoding library](https://github.com/henvic/picel-js) that can be conveniently installed with npm or bower (package name is picel). If you need the encoder library available for another language let me know.

If you need to write URLs by hand take a look at the [examples for the JS encoder](https://github.com/henvic/picel-js#examples) and [its source code](https://github.com/henvic/picel-js/blob/master/picel.js).

## test_assets branch
The test_assets branch exists to the sole purpose of serving as a branch for binary image files for the integration tests. It only contains binary files (and nothing else) and maybe rebased at any time since its history doesn't matter. It's used to checkout the `test_assets` directory whenever the tests are run.

Currently the images on the test_assets branch, despite [AdditiveColor](https://commons.wikimedia.org/wiki/File:AdditiveColor.svg), are thw following photos I've taken and published on [my Flickr account](https://www.flickr.com/photos/henriquev):

* [Golden Gate Bridge after the sunset](https://www.flickr.com/photos/henriquev/8872926264)
* [Rocks & waves @ Big Sur #1](https://www.flickr.com/photos/henriquev/11274440243)
* [Rocks & waves @ Big Sur #2](https://www.flickr.com/photos/henriquev/14085712935)
* [Big Sur, CA #3](https://www.flickr.com/photos/henriquev/9741409523)
* [Insects](https://www.flickr.com/photos/henriquev/8544618839)
* [Raccoons](https://www.flickr.com/photos/henriquev/16100340385)

