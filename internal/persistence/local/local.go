package local

import (
	"fmt"
	"github.com/dagu-dev/dagu/internal/dag"
	"github.com/dagu-dev/dagu/internal/grep"
	"github.com/dagu-dev/dagu/internal/persistence"
	"github.com/dagu-dev/dagu/internal/utils"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type dagStoreImpl struct {
	dir string
}

func NewDAGStore(dir string) persistence.DAGStore {
	return &dagStoreImpl{
		dir: dir,
	}
}

func (d *dagStoreImpl) Create(name string, tmpl []byte) (string, error) {
	if err := d.ensureDirExist(); err != nil {
		return "", fmt.Errorf("failed to create DAGs directory %s", d.dir)
	}
	loc, err := d.fileLocation(name)
	if err != nil {
		return "", fmt.Errorf("failed to create DAG file: %s", err)
	}
	if d.exists(loc) {
		return "", fmt.Errorf("the DAG file %s already exists", loc)
	}
	return name, os.WriteFile(loc, tmpl, 0644)
}

func (d *dagStoreImpl) exists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func (d *dagStoreImpl) fileLocation(name string) (string, error) {
	loc := path.Join(d.dir, name)
	return d.normalizeFilename(loc)
}

func (d *dagStoreImpl) normalizeFilename(file string) (string, error) {
	f, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	a := strings.TrimSuffix(f, ".yaml")
	a = strings.TrimSuffix(f, ".yml")
	return fmt.Sprintf("%s.yaml", a), nil
}

func (d *dagStoreImpl) ensureDirExist() error {
	if !d.exists(d.dir) {
		if err := os.MkdirAll(d.dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func (d *dagStoreImpl) List() ([]dag.DAG, error) {
	//TODO implement me
	panic("implement me")
}

func (d *dagStoreImpl) Grep(pattern string) (ret []*persistence.GrepResult, errs []string, err error) {
	if err = d.ensureDirExist(); err != nil {
		errs = append(errs, fmt.Sprintf("failed to create DAGs directory %s", d.dir))
		return
	}

	fis, err := os.ReadDir(d.dir)
	dl := &dag.Loader{}
	opts := &grep.Options{
		IsRegexp: true,
		Before:   2,
		After:    2,
	}

	utils.LogErr("read DAGs directory", err)
	for _, fi := range fis {
		if utils.MatchExtension(fi.Name(), dag.EXTENSIONS) {
			fn := filepath.Join(d.dir, fi.Name())
			utils.LogErr("read DAG file", err)
			m, err := grep.Grep(fn, fmt.Sprintf("(?i)%s", pattern), opts)
			if err != nil {
				errs = append(errs, fmt.Sprintf("grep %s failed: %s", fi.Name(), err))
				continue
			}
			d, err := dl.LoadMetadataOnly(fn)
			if err != nil {
				errs = append(errs, fmt.Sprintf("check %s failed: %s", fi.Name(), err))
				continue
			}
			ret = append(ret, &persistence.GrepResult{
				Name:    strings.TrimSuffix(fi.Name(), path.Ext(fi.Name())),
				DAG:     d,
				Matches: m,
			})
		}
	}
	return ret, errs, nil
}

func (d *dagStoreImpl) Load(name string) (*dag.DAG, error) {
	//TODO implement me
	panic("implement me")
}

func (d *dagStoreImpl) Rename(oldDAGPath, newDAGPath string) error {
	oldLoc, err := d.fileLocation(oldDAGPath)
	if err != nil {
		return fmt.Errorf("invalid old name: %s", oldDAGPath)
	}
	newLoc, err := d.fileLocation(newDAGPath)
	if err != nil {
		return fmt.Errorf("invalid new name: %s", newDAGPath)
	}
	if err := os.Rename(oldLoc, newLoc); err != nil {
		return err
	}
	return nil
}
