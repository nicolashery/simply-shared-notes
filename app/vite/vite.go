package vite

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"os"
	"strings"
)

type Vite struct {
	IsDev      bool
	PublicFS   fs.FS
	AssetsFS   fs.FS
	Manifest   Manifest
	DevPrefix  string
	ProdPrefix string
}

type Manifest map[string]*Chunk

type Chunk struct {
	File string `json:"file"`
}

func New(logger *slog.Logger, isDev bool, distFS embed.FS) (*Vite, error) {
	var publicFS fs.FS
	var assetsFS fs.FS
	var err error
	if isDev {
		publicFS = os.DirFS("./public")
		assetsFS = os.DirFS("./assets")

		logger.Info("using dev assets, make sure Vite is running")
	} else {
		publicFS, err = fs.Sub(distFS, "dist")
		if err != nil {
			return nil, fmt.Errorf("failed to create sub-filesystem for 'dist' directory: %w", err)
		}
		assetsFS, err = fs.Sub(publicFS, "assets")
		if err != nil {
			return nil, fmt.Errorf("failed to create sub-filesystem for 'dist/assets' directory: %w", err)
		}

		logger.Info("using prod assets")
	}

	var manifest Manifest
	if !isDev {
		err = parseManifest(publicFS, &manifest)
		if err != nil {
			return nil, err
		}
	}

	return &Vite{
		IsDev:      isDev,
		PublicFS:   publicFS,
		AssetsFS:   assetsFS,
		Manifest:   manifest,
		DevPrefix:  "http://localhost:5173",
		ProdPrefix: "",
	}, nil
}

func parseManifest(distFS fs.FS, manifest *Manifest) error {
	mf, err := distFS.Open(".vite/manifest.json")
	if err != nil {
		return fmt.Errorf("failed to open manifest file: %w", err)
	}
	defer mf.Close()

	err = json.NewDecoder(mf).Decode(manifest)
	if err != nil {
		return fmt.Errorf("failed to parse manifest file: %w", err)
	}

	return nil
}

func (v *Vite) Fragment(entries []string) template.HTML {
	var sb strings.Builder

	if v.IsDev {
		sb.WriteString(v.makeScriptTag("@vite/client"))

		for _, entry := range entries {
			if isCSS(entry) {
				sb.WriteString(v.makeStylesheetTag(entry))
			} else {
				sb.WriteString(v.makeScriptTag(entry))
			}
		}
	} else {
		for _, entry := range entries {
			chunk, ok := v.Manifest[entry]
			if !ok {
				continue
			}

			if isCSS(entry) {
				sb.WriteString(v.makeStylesheetTag(chunk.File))
			} else {
				sb.WriteString(v.makeScriptTag(chunk.File))
			}
		}
	}

	return template.HTML(sb.String())
}

func (v *Vite) makeStylesheetTag(path string) string {
	var prefix string
	if v.IsDev {
		prefix = v.DevPrefix
	} else {
		prefix = v.ProdPrefix
	}

	return fmt.Sprintf(
		`<link rel="stylesheet" href="%s/%s" />`,
		prefix,
		path,
	)
}

func (v *Vite) makeScriptTag(path string) string {
	var prefix string
	if v.IsDev {
		prefix = v.DevPrefix
	} else {
		prefix = v.ProdPrefix
	}

	return fmt.Sprintf(
		`<script type="module" src="%s/%s"></script>`,
		prefix,
		path,
	)
}

func isCSS(path string) bool {
	return strings.HasSuffix(path, ".css")
}
