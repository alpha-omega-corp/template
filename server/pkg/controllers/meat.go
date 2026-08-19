package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"template/pkg/handlers"
	"template/pkg/models"

	"github.com/alpha-omega-corp/_core/app"
	"github.com/alpha-omega-corp/_core/helpers"
	"github.com/uptrace/bunrouter"
)

// StorageDir is where uploaded meat images are written, relative to the
// server's working directory.
const StorageDir = "storage"

var allowedImageExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
}

type MeatController struct {
	app.BaseController[*handlers.MeatHandler]
}

var Meat = &MeatController{}

func (c *MeatController) GetAll() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		return helpers.Response(w, func() (*[]models.Meat, error) {
			return c.Handler().GetAll(req.Context())
		})
	}
}

func (c *MeatController) GetOne() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		params := helpers.GetParams[helpers.GetOne](w, req)
		return helpers.Response(w, func() (*models.Meat, error) {
			return c.Handler().GetOne(req.Context(), params.Id)
		})
	}
}

func (c *MeatController) Create() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		data := helpers.GetBody[models.Meat](w, req)
		return helpers.Response(w, func() (*models.Meat, error) {
			return c.Handler().Create(req.Context(), data)
		})
	}
}

func (c *MeatController) Update() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		data := helpers.GetBody[models.Meat](w, req)
		return helpers.Response(w, func() (*models.Meat, error) {
			return c.Handler().Update(req.Context(), data)
		})
	}
}

func (c *MeatController) Delete() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		params := helpers.GetParams[helpers.GetOne](w, req)
		return helpers.Response(w, func() (*models.Meat, error) {
			return nil, c.Handler().Delete(req.Context(), params.Id)
		})
	}
}

func (c *MeatController) UploadImage() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		params := helpers.GetParams[helpers.GetOne](w, req)

		if err := req.ParseMultipartForm(10 << 20); err != nil {
			helpers.Error(w, err, http.StatusBadRequest)
			return nil
		}

		file, header, err := req.FormFile("image")
		if err != nil {
			helpers.Error(w, err, http.StatusBadRequest)
			return nil
		}
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(filepath.Base(header.Filename)))
		if !allowedImageExt[ext] {
			helpers.Error(w, fmt.Errorf("unsupported image type %q", ext), http.StatusBadRequest)
			return nil
		}

		filename := fmt.Sprintf("%d%s", params.Id, ext)
		dst, err := os.Create(filepath.Join(StorageDir, filename))
		if err != nil {
			helpers.Error(w, err, http.StatusInternalServerError)
			return nil
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			helpers.Error(w, err, http.StatusInternalServerError)
			return nil
		}

		return helpers.Response(w, func() (*models.Meat, error) {
			return c.Handler().SetImage(req.Context(), params.Id, filename)
		})
	}
}
