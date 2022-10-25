package conf4g

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMakeConfigFunction(t *testing.T) {

	/*
		MakeConfig()

		--> Configuration{&sync.Mutex{}, configpath, map[string]section{}}
	*/

	Convey("MakeConfig Function", t, func() {
		Convey("MakeConfig Default", func() {
			So(MakeConfig(), ShouldNotBeNil)
		})
	})
}

func TestInitializeFunction(t *testing.T) {

	/*
		Initialize(...path string)

		filepath : /root/work/config
		initargs : /root/work/test
		Initialize()

		--> config : /root/work/config/test.ini

		Initialize("config/master.ini")

		--> config : /root/work/config/master.ini
	*/

	Convey("Initialize Function", t, func() {
		Convey("Initialize Default Path", func() {
			conf := MakeConfig()
			So(conf.Initialize(), ShouldBeNil)
			So(conf.Write("Initialize", "DefaultPath", "True"), ShouldBeNil)
		})

		Convey("Initialize Other Path", func() {
			conf := MakeConfig()
			So(conf.Initialize("config/master.ini"), ShouldBeNil)
			So(conf.Write("Initialize", "OtherPath", "True"), ShouldBeNil)

			// remove master.ini
			conf.Clear()
			os.RemoveAll(conf.confpath)
			os.RemoveAll(conf.confpath + ".bak")
		})

		Convey("Initialize Wrong Path", func() {
			conf := MakeConfig()
			So(conf.Initialize("//"), ShouldBeNil)
			So(conf.Write("Initialize", "WrongPath", "True"), ShouldNotBeNil)
		})

		Convey("Initialize Invalid Parameter", func() {
			conf := MakeConfig()
			So(conf.Initialize(true), ShouldNotBeNil)
		})
	})
}

func TestGetCurrentPathFunction(t *testing.T) {

	/*
		variable.GetCurrentPath()

		filepath : /root/work/config/test.ini
		variable.GetCurrentPath()

		--> /root/work/config/test.ini
	*/

	Convey("GetCurrentPath Function", t, func() {
		Convey("GetCurrentPath Path Exist", func() {
			conf := MakeConfig()
			conf.confpath = "C:\\GoCode\\src\\github.com\\sh26\\sh26lib-go\\conf4g\\config\\conf4g.ini"
			_, err := conf.GetCurrentPath()

			So(err, ShouldBeNil)
		})

		Convey("GetCurrentPath Empty", func() {
			conf := MakeConfig()
			conf.confpath = ""
			_, err := conf.GetCurrentPath()
			So(err, ShouldNotBeNil)
		})
	})
}

func TestReadFunction(t *testing.T) {

	/*
		variable.Read()

		filepath : /root/work/config/conf4g.ini

		conf4g.Initialize()
		conf4g.Read()

		--> variable{[{config data}]}
	*/

	Convey("Read Function", t, func() {
		Convey("Read Path Exist", func() {
			conf := MakeConfig()
			conf.Initialize()

			So(conf.Read(), ShouldBeNil)
		})

		Convey("Read Path Empty", func() {
			conf := MakeConfig()

			So(conf.Read(), ShouldNotBeNil)
		})
	})

}

func TestWriteFunction(t *testing.T) {

	/*
		variable.Write("section", "key", "value")

		filepath : /root/work/config/conf4g.ini

		conf4g.Initialize()
		conf4g.Write("Write", "Hello", World)

		-->
		[Write]
		Hello=World
	*/

	Convey("Write Function", t, func() {
		Convey("Write Something", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			So(conf.Write("Write", "Hello", "World"), ShouldBeNil)
		})

		Convey("Write Section Empty", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			So(conf.Write("", "Hello", "World"), ShouldNotBeNil)
		})

		Convey("Write Key Empty", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			So(conf.Write("Write", "", "World"), ShouldNotBeNil)
		})

		Convey("Write Value Empty", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			So(conf.Write("Write", "Hello", ""), ShouldNotBeNil)
		})

		Convey("Write Value Update", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Write", "Hello", "World")
			conf.Write("Write", "Bye", "World")
			So(conf.Find("Write", "Bye"), ShouldEqual, "World")
		})

		Convey("Write File Create", func() {
			conf := MakeConfig()
			conf.Initialize("config/create.ini")

			So(conf.Write("New", "File", "Create"), ShouldBeNil)

			// remove create.ini
			conf.Clear()
			os.RemoveAll(conf.confpath)
			os.RemoveAll(conf.confpath + ".bak")
		})

		Convey("Write Directory Create", func() {
			conf := MakeConfig()
			conf.Initialize("config/new/create.ini")

			So(conf.Write("New", "File", "Create"), ShouldBeNil)

			// remove new directory
			conf.Clear()
			os.RemoveAll(filepath.Dir(conf.confpath))
		})
	})
}

func TestDeleteSectionFunction(t *testing.T) {

	/*
		variable.DeleteSection(section)

		configdata :

		[Create]
		Hello=World
		[Delete]
		Bye=World

		variable.DeleteSection("Create")

		-->
		[Delete]
		Bye=World
	*/

	Convey("DeleteSection Function", t, func() {
		Convey("DeleteSection Delete", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")
			conf.Write("Section002", "Key001", "Value001")
			conf.Write("Section002", "Key002", "Value002")

			So(conf.DeleteSection("Section001"), ShouldBeNil)
		})

		Convey("DeleteSection Empty", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")
			conf.Write("Section002", "Key001", "Value001")
			conf.Write("Section002", "Key002", "Value002")

			So(conf.DeleteSection(""), ShouldNotBeNil)
		})

		Convey("DeleteSection Wrong", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")
			conf.Write("Section002", "Key001", "Value001")
			conf.Write("Section002", "Key002", "Value002")

			So(conf.DeleteSection("Section003"), ShouldBeNil)
		})
	})
}

func TestDeleteValueFunction(t *testing.T) {

	/*
		variable.DeleteValue(section, key)

		configdata :

		[Create]
		Hello=World
		Start=World
		[Delete]
		Bye=World
		End=World

		variable.DeleteSection("Create", "Hello")
		variable.DeleteSection("Delete", "End")

		-->
		[Create]
		Start=World
		[Delete]
		Bye=World
	*/

	Convey("DeleteValue Function", t, func() {
		Convey("DeleteValue Delete", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")

			So(conf.DeleteValue("Section001", "Key001"), ShouldBeNil)
		})

		Convey("DeleteSection Section Empty", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")

			So(conf.DeleteValue("", "Key001"), ShouldNotBeNil)
		})

		Convey("DeleteSection Key Empty", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")

			So(conf.DeleteValue("Section001", ""), ShouldNotBeNil)
		})

		Convey("DeleteSection Wrong", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")

			So(conf.DeleteValue("Section001", "Key003"), ShouldBeNil)
		})
	})
}

func TestExistSectionFunction(t *testing.T) {

	/*
		variable.ExistSection(section)

		configdata :

		[Create]
		Hello=World
		[Delete]
		Bye=World
		[Remove]
		Re=World

		variable.ExistSection("Remove")

		--> *section
		[Remove]
		Re=World
	*/

	Convey("ExistSection Function", t, func() {
		Convey("ExistSection Exist", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section002", "Key001", "Value001")

			sec, err := conf.ExistSection("Section001")

			So(sec, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})

		Convey("ExistSection Empty", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section002", "Key001", "Value001")

			sec, err := conf.ExistSection("")

			So(sec, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		Convey("ExistSection Wrong", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section002", "Key001", "Value001")

			sec, err := conf.ExistSection("Section003")

			So(sec, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestExistValueFunction(t *testing.T) {

	/*
		variable.ExistValue(section, key)

		configdata :

		[Create]
		Hello=unus
		Start=duo
		[Delete]
		Bye=tres
		End=quattuor

		variable.ExistValue("Delete", "End")

		--> quattuor
	*/

	Convey("ExistSection Function", t, func() {
		Convey("ExistSection Exist", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")

			val, err := conf.ExistValue("Section001", "Key001")

			So(val, ShouldNotBeEmpty)
			So(err, ShouldBeNil)
		})

		Convey("ExistSection Section Empty", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")

			val, err := conf.ExistValue("", "Key001")

			So(val, ShouldBeEmpty)
			So(err, ShouldNotBeNil)
		})

		Convey("ExistSection Key Empty", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")

			val, err := conf.ExistValue("Section001", "")

			So(val, ShouldBeEmpty)
			So(err, ShouldNotBeNil)
		})

		Convey("ExistSection Wrong", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")

			val, err := conf.ExistValue("Section002", "Key001")

			So(val, ShouldBeEmpty)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestGetSectionListFunction(t *testing.T) {

	/*
		variable.GetSectionList()

		configdata :

		[Create]
		one=unus
		[Delete]
		two=duo
		[Update]
		three=tres

		variable.GetSectionList()

		--> [Create, Delete, Update]
	*/

	Convey("GetSectionList Function", t, func() {
		Convey("GetSectionList All", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section002", "Key002", "Value002")
			conf.Write("Section003", "Key003", "Value003")

			So(conf.GetSectionList(), ShouldNotBeNil)
		})

		Convey("GetSectionList Empty", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			// conf.Write("Section001", "Key001", "Value001")
			// conf.Write("Section002", "Key002", "Value002")
			// conf.Write("Section003", "Key003", "Value003")

			So(conf.GetSectionList(), ShouldBeNil)
		})
	})
}

func TestGetKeyListFunction(t *testing.T) {

	/*
		variable.GetKeyList(section)

		configdata :

		[Create]
		one=unus
		two=duo
		three=tres

		variable.GetKeyList()

		--> [Create, Delete, Update]
	*/

	Convey("GetKeyList Function", t, func() {
		Convey("GetKeyList Select", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")
			conf.Write("Section001", "Key003", "Value003")

			So(conf.GetKeyList("Section001"), ShouldNotBeNil)
		})

		Convey("GetSectionList Section Empty", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()

			So(conf.GetKeyList(""), ShouldBeNil)
		})

		Convey("GetSectionList Key Empty", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")

			conf.DeleteValue("Section001", "Key001")

			So(conf.GetKeyList("Section001"), ShouldBeNil)
		})

		Convey("GetSectionList Section Wrong", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")
			conf.Write("Section001", "Key003", "Value003")

			So(conf.GetKeyList("Section002"), ShouldBeNil)
		})
	})
}

func TestFindFunction(t *testing.T) {

	/*
		variable.Find(section, key)

		configdata :

		[Create]
		one=unus
		two=duo
		[Delete]
		three=tres
		four=quattuor

		variable.Find("Create", "two")

		--> duo
	*/

	Convey("Find Function", t, func() {
		Convey("Find Value", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")

			So(conf.Find("Section001", "Key001"), ShouldNotBeEmpty)
		})

		Convey("Find Section Empty", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")

			So(conf.Find("", "Key001"), ShouldBeEmpty)
		})

		Convey("Find Key Empty", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")

			So(conf.Find("Section001", ""), ShouldBeEmpty)
		})

		Convey("Find Wrong", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Clear()
			conf.Write("Section001", "Key001", "Value001")
			conf.Write("Section001", "Key002", "Value002")

			So(conf.Find("Section002", "Key001"), ShouldBeEmpty)
		})
	})

}

func TestClearFunction(t *testing.T) {
	/*
		variable.Clear()

		configdata :

		[Create]
		one=unus
		[Delete]
		two=duo

		variable.Clear()

		--> {empty}
	*/

	Convey("Clear Function", t, func() {
		Convey("Clear Data", func() {
			conf := MakeConfig()
			conf.Initialize()

			conf.Write("Section001", "Key001", "Value001")
			conf.Clear()

			So(conf.Find("Section001", "Key001"), ShouldBeEmpty)
		})
	})
}

func TestStatusFunction(t *testing.T) {

	/*
		variable.Status()

		filepath : /root/work/config/conf4g.ini

		variable.Status()

		--> nil
	*/

	Convey("Status Function", t, func() {
		Convey("Status Exist", func() {
			conf := MakeConfig()
			conf.Initialize()

			So(conf.Status(), ShouldBeNil)
		})
	})

	Convey("Status Function", t, func() {
		Convey("Status Path Empty", func() {
			conf := MakeConfig()

			So(conf.Status(), ShouldNotBeNil)
		})
	})
}
