# Python｜发送邮件



[TOC]



## 一、 smtplib模块

主要通过SMTP类与邮件系统进行交互。


- 实例化一个SMTP对象：  
s = smtplib.SMTP(邮件服务地址)  
s = smtplib.SMTP_SSL(邮件服务地址，端口号)  


- 登陆邮件，权限验证：  
s.login(用户名，密码)


- 发送邮件：  
s.sendmail(发件人邮箱，收件人邮箱，发送内容)


- 断开连接：  
s.close()



## 二、email模块



支持发送的邮件内容为纯文本、HTML内容、图片、附件。


- MIMEText：（MIME媒体类型）内容形式为纯文本、HTML页面。

from email.mime.text import MIMEText

**MIMEText(msg,type,chartset)**

msg：文本内容

type：文本类型默认为plain（纯文本），发送HTML格式的时候，修改为html，但同时要求msg的内容也是html的格式。

chartset：文本编码，中文为“utf-8”

**构造TEXT格式的消息**

msg = MIMEText("hello.text","plain","utf-8")

msg["Subject"] = "xxxxx"

msg["From"] = "xxxx"

msg["To"] = "xxxx"

**发送以上构造的邮件内容要使用as_string将构造的邮件内容转换为string形式。**

s.sendmail("xxx","xxx",msg.as_string)


- MIMEImage：内容形式为图片。

from email.mime.image import MIMEImage


- MIMEMultupart：多形式组合，可包含文本和附件。

from email.mime.multipart import MIMEMultipart

**MIMEMultipart()**

**构造邮件正文**

msg_sub = MIMEText("hello.text","plain","utf-8")

msg.attach(msg_sub)

**添加邮件附件**

msg_img = MIMEImage(open(os.getcwd()+ "/reports/xxxx.png","rb").read())

msg_img.add_header('Content-Disposition','attachment', filename = "xxxx.png" )

msg_img.add_header('Content-ID','<0>')

msg.attach(mag_img)



## 三、函数



```python

#!/usr/bin/python2
# -*- coding: UTF-8 -*-

# 实现发送邮件功能
import smtplib
from email.header import Header
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart


class Sender:
    def __init__(self, host, user, password, receivers):
        # host 邮件服务地址
        self.host = host
        # user 发送人邮箱
        # password 发送人邮箱授权码
        self.user = user
        self.password = password
        # receivers 收件人邮箱
        self.receivers = receivers

    #
    # 创建普通邮件正文
    #
    def setmsg(self, _subject='主题', _text='邮件正文',
               _from=None, _to=None, _cc=None,
               _subtype='plain', _charset='utf-8'):
        # subtype 邮件内容类型{'plain' 文本格式,'html' H5格式}
        msg = MIMEText(_text, _subtype, _charset)
        msg['Subject'] = Header(_subject, 'utf-8')
        msg['From'] = _from
        msg['To'] = _to
        msg['Cc'] = _cc
        # msg 邮件内容
        self.msg = msg
        return self.msg

    #
    # 创建包含附件的邮件正文
    #
    def setmsg_attach(self, _subject='主题', _text='邮件正文', file_list=None,
                      _from=None, _to=None, _cc=None,
                      _subtype='plain', _charset='utf-8'):
        msgs = MIMEMultipart()
        msgs['Subject'] = Header(_subject, 'utf-8')
        msgs['From'] = _from
        msgs['To'] = _to
        msgs['Cc'] = _cc
        # msgs 邮件内容
        msgs.attach(MIMEText(_text, _subtype, _charset))

        # 生成邮件附件内容
        def add_attach(file_name):
            att1 = MIMEText(open(file_name, 'rb').read(), 'base64', 'utf-8')
            att1["Content-Type"] = 'application/octet-stream'
            att1["Content-Disposition"] = 'attachment; filename="%s"' % file_name
            return att1

        # file_list 需要添加进邮件的文件列表
        for file_name in file_list:
            msgs.attach(add_attach(file_name))

        self.msg = msgs
        return self.msg

    #
    # 登录邮箱，发送邮件 【设置-账户-POP/SMTP：开启】
    #
    def sendmail(self):
        if not self.msg:
            print "message is nil"
            return False

        try:
            smtpObj = smtplib.SMTP(self.host)
            smtpObj.set_debuglevel(1)  # 输出发送邮件详细过程
            smtpObj.login(self.user, self.password)
            smtpObj.sendmail(self.user, self.receivers,
                             self.msg.as_string())
            print "邮件发送成功"
            return True
        except Exception as e:
            print e


if __name__ == '__main__':
    print "Test"

    send = Sender(host="smtp.qq.com",
                  user="1837565816@qq.com",
                  password="suwvrojxmubscieh",  # 授权码
                  receivers="xuefeng.han@mintegral.com")

    content = """         #内容，HTML格式
    <p>Python 邮件发送测试...</p>
    <p><a href="http://www.baidu.com">这是一个链接</a></p>
    """
    send.setmsg("主题<测试>", content,
                None, None, None,
                'html')
    send.sendmail()
    # Output:
    # 邮件发送成功
 

```