package mongo

import (
	"fmt"
	mgo "gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	p "org/shoper/app/utils/prop"
	"strconv"
	"strings"
	"time"
)

var (
	info *mgo.DialInfo
)

func init() {
	prop, err := p.Load("/opt/test.properties")
	if err != nil {
		fmt.Println(err.Error())
	}
	url := prop.GetData("url")
	timeout, e := strconv.Atoi(prop.GetData("timeout"))
	if e != nil {
		panic(e)
	}
	user := prop.GetData("user")
	password := prop.GetData("password")
	if user != "" && password != "" {
		info.Username = user
		info.Password = password
	}
	urls := strings.Split(url, ";")
	info = new(mgo.DialInfo)
	info.Addrs = urls
	info.Timeout = time.Second * time.Duration(timeout/1000)
	//Seconds * timeout / parseInt(1000)
	/*info := mgo.DialInfo{
		Addrs: urls,
		Timeout:time.Duration.Seconds()
	}*/

	/*session, err = Connection(host, port)
	if err != nil {
		panic(err.Error())
	}*/
}
func config() {
	mgo.DialWithInfo(info)
}
func GetCoon() (mgoOpera *MgoOpera, err error) {
	mgoOpera = &MgoOpera{}
	mgoOpera.session, err = mgo.DialWithInfo(info)
	if err != nil {
		return nil, err
	}
	return mgoOpera, nil
}
