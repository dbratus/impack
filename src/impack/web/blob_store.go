// ImPack - CSS sprites maker
// Copyright (C) 2012 Dmitry Bratus
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
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