package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/olivere/vite"
)

type AssetsConfig struct {
	AssetsFS     fs.FS
	PublicFS     fs.FS
	ViteFragment *vite.Fragment
}

func DevAssets() (AssetsConfig, error) {
	viteFragment, err := vite.HTMLFragment(vite.Config{
		IsDev:        true,
		ViteTemplate: vite.None,
		ViteEntry:    "assets/app.js",
	})
	if err != nil {
		return AssetsConfig{}, fmt.Errorf("failed to instantiate Vite fragment: %w", err)
	}

	return AssetsConfig{
		AssetsFS:     os.DirFS("./assets"),
		PublicFS:     os.DirFS("./public"),
		ViteFragment: viteFragment,
	}, nil
}

func ProdAssets(distFS embed.FS) (AssetsConfig, error) {
	viteRootFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		return AssetsConfig{}, fmt.Errorf("failed to create sub-filesystem for 'dist' directory: %w", err)
	}
	viteAssetsFS, err := fs.Sub(viteRootFS, "assets")
	if err != nil {
		return AssetsConfig{}, fmt.Errorf("failed to create sub-filesystem for 'dist/assets' directory: %w", err)
	}

	viteFragment, err := vite.HTMLFragment(vite.Config{
		FS:           viteRootFS, // used to read `.vite/manifest.json`
		IsDev:        false,
		ViteTemplate: vite.None,
		ViteEntry:    "assets/app.js",
	})
	if err != nil {
		return AssetsConfig{}, fmt.Errorf("failed to instantiate Vite fragment: %w", err)
	}

	return AssetsConfig{
		AssetsFS:     viteAssetsFS,
		PublicFS:     viteRootFS,
		ViteFragment: viteFragment,
	}, nil
}
