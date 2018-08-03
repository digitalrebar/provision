package cli

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/digitalrebar/provision/api"
	"github.com/digitalrebar/provision/models"
	"github.com/groob/plist"
	gzip "github.com/klauspost/pgzip"
	"github.com/spf13/cobra"
)

// Args will have been validated by the time this gets called.
// args[0] will have a valid path that contains an
// NBImageInfo.plist file.
func genEnvAndArchiveFromAppleNBI(c *cobra.Command, args []string) error {
	pListFile := path.Join(args[0], "NBImageInfo.plist")
	info, err := os.Open(pListFile)
	if err != nil {
		return err
	}
	defer info.Close()
	bsdpInfo := &models.BsdpBootOption{}
	dec := plist.NewXMLDecoder(info)
	if err := dec.Decode(bsdpInfo); err != nil {
		return err
	}
	env := &models.BootEnv{}
	env.Fill()
	osName := bsdpInfo.OSName() + "-" +
		bsdpInfo.OSVersion + "-" +
		bsdpInfo.InstallType()
	env.Name = osName + "-install"
	env.OS.Name = osName
	if _, err := os.Stat(path.Join(args[0], "i386", bsdpInfo.Booter)); err != nil {
		return fmt.Errorf("NBI missing Booter i386/%s: %v", bsdpInfo.Booter, err)
	}
	if _, err := os.Stat(path.Join(args[0], bsdpInfo.RootPath)); err != nil {
		return fmt.Errorf("NBI missing RootPath %s: %v", bsdpInfo.RootPath, err)
	}
	env.Kernel = path.Join("i386", bsdpInfo.Booter)
	env.OS.IsoFile = osName + ".tar.gz"
	env.Meta["AppleBsdp"] = bsdpInfo.String()
	env.Meta["KernelIsLoader"] = "true"
	log.Printf("Creating bootenv archive %s.tar.gz (this may take awhile)", osName)
	outArchive, err := os.Create(osName + ".tar.gz")
	if err != nil {
		return fmt.Errorf("Error creating new archive %s.tar.gz: %v", osName, err)
	}
	defer outArchive.Close()
	gzWrite, err := gzip.NewWriterLevel(outArchive, gzip.BestCompression)
	if err != nil {
		return fmt.Errorf("Error creating gzip Writer: %v", err)
	}
	defer gzWrite.Close()
	gzWrite.SetConcurrency(1<<20, runtime.NumCPU()*2)
	tarWrite := tar.NewWriter(gzWrite)
	defer tarWrite.Close()
	err = filepath.Walk(args[0], func(item string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, item)
		if err != nil {
			return err
		}
		if header.Typeflag == tar.TypeSymlink {
			header.Linkname, _ = os.Readlink(item)
			if !filepath.IsAbs(header.Linkname) {
				header.Linkname = path.Join(path.Dir(item), header.Linkname)
			}
			header.Linkname, err = filepath.Rel(path.Dir(item), header.Linkname)
			if err != nil {
				return err
			}
		}
		header.Name, err = filepath.Rel(args[0], item)
		if err != nil {
			return err
		}
		if header.Name == "" || header.Name == "." {
			return nil
		}
		if err := tarWrite.WriteHeader(header); err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		f, err := os.Open(item)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(tarWrite, f); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error creating drp boot archive: %v", err)
	}
	log.Printf("Creating bootenv %s.yaml", osName)
	outEnv, err := os.Create(osName + ".yaml")
	if err != nil {
		return fmt.Errorf("Errorf creating %s.yaml: %v", osName, err)
	}
	defer outEnv.Close()
	buf, err := api.Pretty("yaml", env)
	if err != nil {
		return fmt.Errorf("Error marshalling new BootEnv: %v", err)
	}
	if _, err := outEnv.Write(buf); err != nil {
		return fmt.Errorf("Error writing new BootEnv: %v", err)
	}
	return nil
}
