package web

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"archive/zip"
	"fmt"
	"html/template"
	"image"
	"image/draw"
	"image/png"
	"impack"
	"impack/loading"
	"mime/multipart"
	"net/http"
	"path"
	"time"
)

type SpriteInfo struct {
	ImageId   appengine.BlobKey
	CssId     appengine.BlobKey
	Timestamp time.Time
}

type SpriteView struct {
	Sprite  SpriteInfo
	Classes []string
}

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/blob/", blob)
	http.HandleFunc("/js/", static)
}

func index(resp http.ResponseWriter, req *http.Request) {
	http.ServeFile(resp, req, "index.html")
}

func static(resp http.ResponseWriter, req *http.Request) {
	http.ServeFile(resp, req, req.URL.Path)
}

func blob(resp http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	_, id := path.Split(req.URL.Path)
	ctx.Infof("BLOB Id %s parsed.", id)

	blobstore.Send(resp, appengine.BlobKey(id))
}

func upload(resp http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	var arch multipart.File

	if f, _, err := req.FormFile("zip"); err == nil {
		arch = f
	} else {
		ctx.Errorf("%s", err.Error())
		http.Error(resp, "Server error", http.StatusInternalServerError)
		return
	}

	var fileSize int64

	if sz, err := arch.Seek(0, 2); err == nil {
		fileSize = sz
	} else {
		ctx.Errorf("%s", err.Error())
		http.Error(resp, "Server error", http.StatusInternalServerError)
		return
	}

	var images []loading.Image

	if r, err := zip.NewReader(arch, fileSize); err == nil {
		images = loading.LoadImagesFromZip(r)
	} else {
		ctx.Errorf("%s", err.Error())
		http.Error(resp, "Server error", http.StatusInternalServerError)
		return
	}

	rects := make([]image.Rectangle, len(images))

	for i := 0; i < len(rects); i++ {
		rects[i] = images[i].Data.Bounds()
	}

	union := impack.Arrange(rects)
	dest := image.NewNRGBA(image.Rect(0, 0, union.Dx(), union.Dy()))

	for i := 0; i < len(rects); i++ {
		draw.Draw(dest, rects[i], images[i].Data, image.Pt(0, 0), draw.Src)
	}

	spriteInfo := &SpriteInfo{Timestamp: time.Now()}

	if k, err := saveBlob(ctx, "image/png", func(w *blobstore.Writer) error { return png.Encode(w, dest) }); err == nil {
		ctx.Infof("Image saved with key %s.", k)
		spriteInfo.ImageId = k
	} else {
		ctx.Errorf("%s", err.Error())
		http.Error(resp, "Server error", http.StatusInternalServerError)
		return
	}

	if k, err := saveBlob(ctx, "text/css", func(w *blobstore.Writer) error { return generateCss(w, spriteInfo.ImageId, rects, images) }); err == nil {
		ctx.Infof("CSS saved with key %s.", k)
		spriteInfo.CssId = k
	} else {
		ctx.Errorf("%s", err.Error())
		http.Error(resp, "Server error", http.StatusInternalServerError)
		return
	}

	if _, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "spriteInfo", nil), spriteInfo); err != nil {
		ctx.Errorf("%s", err.Error())
		http.Error(resp, "Server error", http.StatusInternalServerError)
		return
	}

	if t, err := template.ParseFiles("result.html"); err == nil {
		classes := make([]string, len(rects))
		for i := 0; i < len(classes); i++ {
			classes[i] = fmt.Sprintf("cls%d", i+1)
		}

		spriteView := &SpriteView{Sprite: *spriteInfo, Classes: classes}

		if err = t.Execute(resp, spriteView); err != nil {
			ctx.Errorf("%s", err.Error())
			http.Error(resp, "Server error", http.StatusInternalServerError)
			return
		}
	} else {
		ctx.Errorf("%s", err.Error())
		http.Error(resp, "Server error", http.StatusInternalServerError)
		return
	}
}
