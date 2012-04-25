package web

import (
	"appengine"
	"appengine/blobstore"
)

func saveBlob(context appengine.Context, mimeType string, saveFunc func(*blobstore.Writer) (error)) (appengine.BlobKey, error) {
	var blobWriter *blobstore.Writer
	
	if w, err := blobstore.Create(context, mimeType); err == nil {
		blobWriter = w
	} else {
		return "", err
	}

	if err := saveFunc(blobWriter); err != nil {
		return "", err
	}

	if err := blobWriter.Close(); err != nil {
		return "", err
	}
	
	if k, err := blobWriter.Key(); err == nil {
		return k, nil
	} else {
		return "", err
	}
	
	return "", nil
}