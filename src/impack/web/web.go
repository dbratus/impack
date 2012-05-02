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
	http.HandleFunc("/cut", cut)
	http.HandleFunc("/blob/", blob)
}

func index(resp http.ResponseWriter, req *http.Request) {
	http.ServeFile(resp, req, "index.html")
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

	if _, err := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "SpriteInfo", nil), spriteInfo); err != nil {
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

func cut(resp http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	minuteBack, _ := time.ParseDuration("-1m")
	cutPoint := time.Now().Add(minuteBack)
	
	q := datastore.NewQuery("SpriteInfo").
		Filter("Timestamp <", cutPoint)
	
	itr := q.Run(ctx)

	var spriteInfo SpriteInfo
	
	cnt := 0
	var k *datastore.Key
	var itrErr error
	
	for k, itrErr = itr.Next(&spriteInfo); itrErr == nil; k, itrErr = itr.Next(&spriteInfo) {
		if err := blobstore.Delete(ctx, spriteInfo.CssId); err != nil {
			ctx.Errorf("Error deleting CSS %s %s.", spriteInfo.CssId, err)
		}
		
		if err := blobstore.Delete(ctx, spriteInfo.ImageId); err != nil {
			ctx.Errorf("Error deleting image %s %s.", spriteInfo.ImageId, err)
		}
		
		if err := datastore.Delete(ctx, k); err != nil {
			ctx.Errorf("Error deleting sprite %s %s.", k, err)
		}
		
		cnt++
	}
	
	ctx.Infof("Cut stoped with result %v. %d sprites deleted", itrErr, cnt)
}
