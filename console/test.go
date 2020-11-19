package console

import (
	"github.com/sirupsen/logrus"
	"time"
)

func TestConsole()  {


	logrus.Debug("TestConsole，执行中：", time.Now()) //注意，如果是每秒执行的就不要调用这个方法，不然导致每秒打印log文件内容过多
	//调用任务module方法处理业务逻辑


}
