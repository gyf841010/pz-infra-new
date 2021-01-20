1. database 数据库(基于gorm)
- 更新说明:2021-1-16,gorm升级为v2,可使用批量插入等新特性,参见 :https://gorm.io/zh_CN/docs/changelog.html
  
- 批量插入,gorm 文档: https://gorm.io/zh_CN/docs/create.html
  ``` 
    // 将切片数据传递给 Create 方法，GORM 将生成一个单一的 SQL 语句来插入所有数据，
    // 并回填主键的值，钩子方法也会被调用。
    db := database.GetDB()
    pa := []PublicAccount{
			{
				Token:  "token_______1",
				Remark: "remard______1",
				AppID:  "app_id______1",
			},
			{
				Token:  "token_______2",
				Remark: "remard______2",
				AppID:  "app_id______2",
			},
			{
				Token:  "token_______3",
				Remark: "remard______3",
				AppID:  "app_id______3",
			},
		}
	err := db.Create(&pa).Error
	if err != nil {
		// TODO:
	}
    // 插入成功后 pa已有每条记录的ID(表主键ID)
    for _,p:=range pa{
        fmt.Println(p.ID)
    }
  ```
- 链式方法影响,
  1. 参考: https://gorm.io/zh_CN/docs/method_chaining.html
  2. https://gorm.io/zh_CN/docs/method_chaining.html#goroutine_safe
  3. https://gorm.io/zh_CN/docs/session.html
  ```
  db := database.GetDB()
  driver111 := &Driver{}
  db.First(driver111, 10000)
  db.First(driver111, 10001)
  // 第一条执行的sql语句会影响后面未开新Statement的sql语句执行
  // 第二条sql如下:
  2021/01/16 16:22:25  record not found
    [0.470ms] [rows:0] 
    SELECT * FROM `driver` WHERE `driver`.`id` = 10001 AND `driver`.`id` = 10000 ORDER BY `driver`.`id` LIMIT 1
  
  ```
2. logging 日志
- 目标:日志只打一次,避免一个错误多次输出
- 注意
   1. 通过errors.New()/WithStack()/Wrap()等处理的错误已包含堆栈信息,不要使用Log.Error()记录日志(会输出双份的堆栈信息,数据冗余)
   ```
   e := errors.New("test_error")		
   Log.ErrorWithStack(e, "test_args")

   // 会输出冗余信息
   Log.Error("t", WithError(e))
   ```
- 样例 1).结合 git .com/pkg/errors 包实现堆
    ```
    // contrller 记录下层返回的错误到日志
    func (controller *Controller)Add(){
        ....
        resp,err:=service.Add(req)
        if err!=nil {
            Log.ErrorWithStack(err,"错误描述,可选")
        }
    }

    ```
    ```
    //service 2种处理方式:1)吞掉下层错误记录日志 2)往上抛
    func Add(req *Req)(*Resp,error){
        ...
        err:=model.Add()
        //1)如果错误需要在这一层吞掉不返回给客户端,需要输出到日志
        if err!=nil{
            ...
            Log.ErrorWithStack(err,"错误描述,可选")
            ...

            return resp,nil
        }
        //2)往controller层抛(可选: 将model层错误转换为业务层的错误)  
        // 往上抛时也可以通过errors.WithMessage(err,"附加描述信息,可选")           
        reutrn nil,err
    }
    ```
    ```
    // model
    // model层发生错误不需要在这一层输出到日志文件,log带上堆栈信息往上一层抛
    func Add()error{
        ...
        err:=db.Create(...).Error
        return errors.WithStack(err)
    }
    ```



- 日志记录样例
  1. 说明:$GOPATH不是真实路径(为方便说明做了处理),实际日志会输出实际路径
  2. 代码
   ```
   package logging

    import (
    	"fmt"
    	"testing"

    	"github.com/gyf841010/pz-infra-new/errorUtil"
    	"github.com/pkg/errors"
    )

    func aTestError() error {
    	err := errors.New("test_error")
    	return err
    }
    func TestLogging(t *testing.T) {
    	InitLogger("test")
    	err := aTestError()
    	// 3. 使用Log.ErrorWithStack()输出带堆栈信息的错误
        Log.ErrorWithStack("my", err, "", With("with msg", "test"), With("struct", struct   {
    		A string `json:"a"`
    	}{A: "some msg"}))

        // 4. 使用Log.Error()输出带堆栈信息错误
    	Log.Error("test", WithError(err))
    }
   ```
  3. 使用Log.ErrorWithStack()输出带堆栈信息的错误
  ```
      ERROR[2021-01-18T20:42:41+08:00] my test_error
    github.com/gyf841010/pz-infra-new/logging.aTestError
    	$GOPATH/pz-infra-new/logging/log_test.go:12
    github.com/gyf841010/pz-infra-new/logging.TestLogging
    	$GOPATH/pz-infra-new/logging/log_test.go:17
    testing.tRunner
    	/usr/local/go/src/testing/testing.go:1123
    runtime.goexit
    	/usr/local/go/src/runtime/asm_amd64.s:1374 [ {with msg test} {struct {some msg}}] 
  ```
  4. 使用Log.Error()输出带堆栈信息错误,会输出冗余信息(Error方法实现中有获取堆栈信息的过程,这里输出,2遍堆栈信息)
   ```
       ERROR[2021-01-18T20:42:42+08:00] test                                            error=test_error
    github.com/gyf841010/pz-infra-new/logging.aTestError
    	$GOPATH/pz-infra-new/logging/log_test.go:12
    github.com/gyf841010/pz-infra-new/logging.TestLogging
    	$GOPATH/pz-infra-new/logging/log_test.go:17
    testing.tRunner
    	/usr/local/go/src/testing/testing.go:1123
    runtime.goexit
    	/usr/local/go/src/runtime/asm_amd64.s:1374 component=test stacktrace=goroutine 21   [running]:
    github.com/gyf841010/pz-infra-new/logging.takeStacktrace(0x419f00, 0x0, 0x0)
    	$GOPATH/pz-infra-new/logging/field.go:54 +0x91
    github.com/gyf841010/pz-infra-new/logging.Stacktrace(0x0, 0x0, 0x0, 0x0)
    	$GOPATH/pz-infra-new/logging/field.go:48 +0x4d
    github.com/gyf841010/pz-infra-new/logging.(*logrusLogger).addFields(0xc000296120,   0xc000203aa0, 0x1, 0x1, 0xc000203a01, 0x0)
    	$GOPATH/pz-infra-new/logging/logrus.go:128 +0x325
    github.com/gyf841010/pz-infra-new/logging.(*logrusLogger).Error(0xc000296120,   0xdf96f9, 0x4, 0xc000203aa0, 0x1, 0x1, 0x0, 0x0)
    	$GOPATH/pz-infra-new/logging/logrus.go:174 +0x77
    github.com/gyf841010/pz-infra-new/logging.TestLogging(0xc000294300)
    	$GOPATH/pz-infra-new/logging/log_test.go:22 +0x464
    testing.tRunner(0xc000294300, 0xe2c300)
    	/usr/local/go/src/testing/testing.go:1123 +0x1a3
    created by testing.(*T).Run
    	/usr/local/go/src/testing/testing.go:1168 +0x648
     fileFunc=log_test.go:22:TestLogging
   ```