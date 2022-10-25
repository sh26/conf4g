// Copyright © 2022 Park Seong Ho <sh26@kakao.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// conf4g provides a simple parser for reading/writing configuration (INI) files.
//
package conf4g

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/alyu/configparser"
)

type section struct {
	name string
	data map[string]string
}

type Configuration struct {
	confpath string
	sections map[string]section

	mu *sync.Mutex
}

// MakeConfig 함수는 Configuration 구조체의 생성자 함수입니다.
// 새로운 Configuration 구조체 타입 포인터 변수의 메모리 주소를 반환합니다.
// Configuration 구조체의 구조는 다음과 같습니다.
// =======================================
//
// confpath	string
// sections	map[string]section{name string, data map[string]{string}}
// mu		*sync.Mutex
//
// confpath			: configuration 파일의 위치입니다.
// sections 		: configuration 파일의 구조를 저장하는 변수입니다.
//          		  각 section은 name과 data로 구성되어져 있습니다.
// sections - name	: section의 이름입니다. section마다 하나만 존재할 수 있습니다.
// sections - data	: section의 내용입니다. 여러개의 [key=value]로 구성되어져 있습니다.
//
// =======================================
func MakeConfig() *Configuration { return &Configuration{} }

// Initialize 함수는 configuration 변수의 내부 값들을 초기화합니다.
// =======================================
//
// levelMax 	: INFO
// dirPath		: {current-work-path}/config
// filename		: {os.Args[0] -> split[0] + .ini}
//
// =======================================
func (conf *Configuration) Initialize(path ...interface{}) error {
	conf.sections = make(map[string]section)

	target, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	conf.mu = &sync.Mutex{}

	if strings.Contains(target, "go-build") {
		target, _ = os.Getwd()
	}

	if path == nil {

		base := filepath.Base(os.Args[0])

		if strings.Contains(base, ".") {
			bound := strings.Split(base, ".")
			base = bound[0] + ".ini"
		}

		conf.confpath = filepath.Join(target, "config/", base)
	} else {
		if reflect.TypeOf(path[0]).Kind() == reflect.String {
			conf.confpath = filepath.Join(target, filepath.Clean(fmt.Sprint(path...)))
		} else {
			return errors.New("Initialize : invalid parameter")
		}
	}

	return nil
}

// GetCurrentPath 함수는 config 파일의 경로를 반환합니다.
// 파일 경로가 정의되지 않았을 경우 에러를 반환합니다.
func (conf *Configuration) GetCurrentPath() (string, error) {
	if conf.confpath == "" {
		return "", errors.New("GetCurrentPath : no path specified")
	}
	return conf.confpath, nil
}

// Read 함수는 config 파일의 내용을 변수에 갱신합니다.
// 파일 경로가 정의되지 않았을 경우 에러를 반환하며 refresh 내부 함수를 호출합니다.
func (conf *Configuration) Read() error {
	if conf.confpath == "" {
		return errors.New("Read : missing configuration path")
	}
	conf.refresh()
	return nil
}

// Write 함수는 config 파일에 내용을 추가 및 갱신합니다.
// 작성 중 mutex의 Lock 함수를 사용하여 동기 처리를 합니다.
// 인자값 중 하나라도 값이 없을 시 에러를 반환합니다.
// 폴더와 파일을 경로에 위치하지 않을 경우, 해당 폴더와 파일을 신규로 생성합니다.
// config 내용의 기록은 다음의 라이브러리를 사용합니다.
//
// https://github.com/alyu/configparser
// =======================================
// config 내용은 다음과 같게 작성됩니다.
//
// [section]
// key=value
//
// section	: Print
// key		: Hello
// value	: World
// -->
// [Print]
// Hello=World
// =======================================
func (conf *Configuration) Write(section, key, value string) (err error) {
	conf.Read()
	conf.mu.Lock()

	defer func() {
		conf.mu.Unlock()
		conf.Read()
	}()

	if section == "" {
		return errors.New("Write : missing section")
	}
	if key == "" {
		return errors.New("Write : missing key")
	}
	if value == "" {
		return errors.New("Write : missing value")
	}

	if ftype, fileerr := exists(conf.confpath); fileerr != nil {
		if _, direrr := exists(filepath.Dir(conf.confpath)); direrr != nil {
			os.MkdirAll(filepath.Dir(conf.confpath), os.ModePerm)
		}

		fi, ferr := os.Create(conf.confpath)
		if ferr != nil {
			return errors.New(fmt.Sprint("Write : cannot create configuration ", ferr))
		}
		fi.Close()
	} else {
		if ftype == 0 {
			return errors.New("Write : target is directory")
		}
	}

	con, cerr := configparser.Read(conf.confpath)
	if cerr != nil {
		return errors.New("Write : cannot read configuration")
	}

	sec, serr := con.Section(section)
	if serr != nil {
		sec = con.NewSection(section)
	}

	if !sec.Exists(key) {
		sec.Add(key, value)
	} else {
		sec.SetValueFor(key, value)
	}

	// 2008 32bit 백업 에러, 추후 원인 분석
	os.Remove(conf.confpath + ".bak")

	err = configparser.Save(con, conf.confpath)

	return nil
}

// DeleteSection 함수는 config 파일에서 section을 삭제합니다.
// section이 지정되지 않을 시 에러를 반환합니다.
func (conf *Configuration) DeleteSection(section string) error {
	conf.Read()
	conf.mu.Lock()

	defer func() {
		conf.mu.Unlock()
		conf.Read()
	}()

	if section == "" {
		return errors.New("DeleteSection : missing section")
	}

	con, cerr := configparser.Read(conf.confpath)

	if cerr != nil {
		return errors.New(fmt.Sprint("DeleteSection : cannot read configuration", cerr))
	} else {
		if _, derr := con.Delete(section); derr != nil {
			return errors.New(fmt.Sprint("DeleteSection : cannot delete section", derr))
		}
	}

	if serr := configparser.Save(con, conf.confpath); serr != nil {
		return errors.New(fmt.Sprint("DeleteSection : cannot save configuration", serr))
	}

	return nil
}

// DeleteValue 함수는 config 파일에서 value를 삭제합니다.
// section과 key가 지정되지 않을 시 에러를 반환합니다.
func (conf *Configuration) DeleteValue(section string, key string) error {
	conf.Read()
	conf.mu.Lock()

	defer func() {
		conf.mu.Unlock()
		conf.Read()
	}()

	if section == "" {
		return errors.New("DeleteValue : missing section")
	}
	if key == "" {
		return errors.New("DeleteValue : missing key")
	}

	con, cerr := configparser.Read(conf.confpath)
	if cerr != nil {
		return errors.New(fmt.Sprint("DeleteValue : cannot read configuration", cerr))
	}

	sec, serr := con.Section(section)
	if serr != nil {
		return errors.New(fmt.Sprint("DeleteValue : cannot load section", cerr))
	}

	sec.Delete(key)

	if serr2 := configparser.Save(con, conf.confpath); serr2 != nil {
		return errors.New(fmt.Sprint("DeleteValue : cannot save configuration", serr2))
	}

	return nil
}

// ExistSection 함수는 config 파일에서 section의 존재여부를 확인합니다.
// section이 지정되지 않을 시 에러를 반환합니다.
func (conf *Configuration) ExistSection(section string) (*section, error) {
	conf.Read()
	if targetsection, ok := conf.sections[section]; ok {
		return &targetsection, nil
	}
	return nil, errors.New("ExistSection : cannot find section")
}

// ExistValue 함수는 config 파일에서 지정 된 section의 value에 대한 존재여부를 확인합니다.
// section과 key가 지정되지 않을 시 에러를 반환합니다.
func (conf *Configuration) ExistValue(section, key string) (string, error) {
	conf.Read()

	if targetsection, serr := conf.ExistSection(section); serr == nil {
		if targetvalue, ok := targetsection.data[key]; ok {
			return targetvalue, nil
		} else {
			return "", errors.New(fmt.Sprint("ExistValue : cannot find value"))
		}
	} else {
		return "", serr
	}
}

// GetSectionList 함수는 config 파일의 모든 section을 string array로 반환합니다.
// section이 존재하지 않을 경우 nil을 반환합니다.
func (conf *Configuration) GetSectionList() []string {
	conf.Read()
	if len(conf.sections) == 0 {
		return nil
	}

	var sectionlist []string

	for name, _ := range conf.sections {
		sectionlist = append(sectionlist, name)
	}

	return sectionlist
}

// GetKeyList 함수는 config 파일의 지정된 section의 모든 key를 string array로 반환합니다
// section이 존재하지 않을 경우 nil을 반환합니다.
func (conf *Configuration) GetKeyList(section string) []string {
	conf.Read()
	if len(conf.sections) == 0 {
		return nil
	}

	if targetsection, serr := conf.ExistSection(section); serr == nil {
		if len(targetsection.data) == 0 {
			return nil
		}
		var keylist []string
		for name, _ := range targetsection.data {
			keylist = append(keylist, name)
		}
		return keylist
	} else {
		return nil
	}
}

// Find 함수는 config 파일의 지정된 section과 key에 대한 value 값을 반환합니다.
// value가 존재하지 않을 경우 공백값을 반환합니다.
func (conf *Configuration) Find(section, key string) string {
	conf.Read()

	if targetsection, sok := conf.sections[section]; sok {
		if targetvalue, vok := targetsection.data[key]; vok {
			return targetvalue
		}
	}
	return ""
}

// Clear 함수는 config 파일의 모든 내용을 삭제합니다.
// 내부 함수인 clear 함수를 호출합니다.
func (conf *Configuration) Clear() error { return conf.clear() }

// Status 함수는 config 파일의 존재 여부를 반환합니다.
// 파일 경로가 정의되지 않았을 경우 에러를 반환합니다.
func (conf *Configuration) Status() error {
	if conf.confpath == "" {
		return errors.New("Status : config cannot read")
	}
	_, fileerr := exists(conf.confpath)
	return fileerr
}

// clear 함수는 config 파일의 모든 내용을 삭제합니다.
// 삭제 도중 치명적인 문제가 발생할 경우 에러를 반환합니다.
func (conf *Configuration) clear() error {
	conf.Read()

	con, cerr := configparser.Read(conf.confpath)
	if cerr != nil {
		return errors.New(fmt.Sprint("clear : config cannot read,", cerr))
	}

	sec, serr := con.AllSections()
	if serr != nil {
		return errors.New(fmt.Sprint("clear : section cannot read,", serr))
	}

	sec = sec[1:]
	for _, tempsec := range sec {
		if derr := conf.DeleteSection(tempsec.Name()); derr != nil {
			return errors.New(fmt.Sprint("clear : ", derr))
		}
	}
	return nil
}

// refresh 함수는 config 파일 내용을 변수에 갱신합니다
// 변수 내용 작성 중 mutex의 Lock 함수를 사용하여 동기 처리를 합니다.
func (conf *Configuration) refresh() (ret error) {
	conf.mu.Lock()

	defer func() {
		conf.mu.Unlock()
		if err := recover(); err != nil {
			// error
			ret = errors.New(fmt.Sprint("refresh : ", err))
		}
	}()

	if _, fileerr := exists(conf.confpath); fileerr != nil {
		ret = errors.New(fmt.Sprint("refresh : ", fileerr))
	}

	conf.sections = map[string]section{}
	con, cerr := configparser.Read(conf.confpath)
	if cerr != nil {
		ret = errors.New(fmt.Sprint("refresh : config cannot read,", cerr))
	}

	sec, serr := con.AllSections()
	if serr != nil {
		ret = errors.New(fmt.Sprint("refresh : section cannot read,", serr))
	}

	sec = sec[1:]
	for _, tempsec := range sec {
		conf.sections[tempsec.Name()] = section{
			name: tempsec.Name(),
			data: tempsec.Options(),
		}
	}

	return
}

func Exists(target string) (int, error) {
	return exists(target)
}

func exists(target string) (int, error) {
	cleanpath := filepath.Clean(target)

	fi, err := os.Stat(cleanpath)
	if err != nil {
		return 4, errors.New("Exists : invalid filepath")
	}

	switch ftype := fi.Mode(); {
	case ftype.IsDir():
		// directory
		return 0, nil
	case ftype.IsRegular():
		// file
		return 1, nil
	default:
		return 2, nil
	}
}
