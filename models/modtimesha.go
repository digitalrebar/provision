package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/pkg/xattr"
)

type ModTimeSha struct {
	ModTime time.Time
	ShaSum  []byte
}

func (m *ModTimeSha) MarshalBinary() ([]byte, error) {
	res := []byte{}
	res = append(res, byte(len(m.ShaSum)))
	res = append(res, m.ShaSum...)
	t, err := m.ModTime.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return append(res, t...), nil
}

func (m *ModTimeSha) UnmarshalBinary(buf []byte) error {
	hStop := int(buf[0] + 1)
	m.ShaSum = append([]byte{}, buf[1:hStop]...)
	return m.ModTime.UnmarshalBinary(buf[hStop:])
}

func (m *ModTimeSha) String() string {
	return hex.EncodeToString(m.ShaSum)
}

func (m *ModTimeSha) UpToDate(fi *os.File) bool {
	stat, err := fi.Stat()
	return err == nil && !stat.IsDir() && stat.ModTime().Equal(m.ModTime)
}

func (m *ModTimeSha) Generate(fi *os.File) error {
	stat, err := fi.Stat()
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return fmt.Errorf("Cannot generate modtimesha on a directory")
	}
	mtime := stat.ModTime()
	if _, err := fi.Seek(0, io.SeekStart); err != nil {
		return err
	}
	shasum := sha256.New()
	sz, err := io.Copy(shasum, fi)
	fi.Seek(0, io.SeekStart)
	if err != nil || sz != stat.Size() {
		return fmt.Errorf("Failed to calculate shasum for %s", fi.Name())
	}
	m.ModTime = mtime
	m.ShaSum = shasum.Sum(nil)
	return nil
}

func (m *ModTimeSha) ReadFromXattr(fi *os.File) error {
	var buf []byte
	var err error
	if runtime.GOOS != "windows" {
		buf, err = xattr.FGet(fi, "user.drpetag")
	} else {
		buf, err = ioutil.ReadFile(fi.Name() + ":drpEtag")
	}
	if err != nil {
		return err
	}
	return m.UnmarshalBinary(buf)
}

func (m *ModTimeSha) SaveToXattr(fi *os.File) error {
	xb, _ := m.MarshalBinary()
	var err error
	if runtime.GOOS != "windows" {
		err = xattr.FSet(fi, "user.drpetag", xb)
	} else {
		err = ioutil.WriteFile(fi.Name()+":drpEtag", xb, 0600)
	}
	return err

}
