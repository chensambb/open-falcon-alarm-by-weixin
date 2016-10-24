#coding=utf-8
'''
Created on 2016-10-17

@author: chenguomin  qq:29235373

@explain: 实现GET方法和POST方法请求
'''
from  BaseHTTPServer import HTTPServer,BaseHTTPRequestHandler
import urllib
import urllib2

def send_message(tos,txt):
'''
此函数主要用于调用我公司内部的短信API(get方式)，大伙可自行修改。
比如：需要加用户名密码,MD5，token的自行查看短信平台提供商，如何调用短信API。
下面提供一个url供大家参考。
'''
    #url  =  "http://www.sms.com:4000/sendsms?uid=16888&pwd=123456&mobile=%s&msg=%s" % (tos, txt)         
    url = "http://192.168.20.88:8080/SendSms/sendsms?mobile=%s&content=%s" % (tos, txt)
    req = urllib2.Request(url)
    res_data = urllib2.urlopen(req)
    res = res_data.read()
    print res


class ServerHTTP(BaseHTTPRequestHandler):
    def do_GET(self):
        path = self.path
        print path
        #拆分url(也可根据拆分的url获取Get提交才数据),可以将不同的path和参数加载不同的html页面，或调用不同的方法返回不同的数据，来实现简单的网站或接口
        query = urllib.splitquery(path)
        print query
        self.send_response(200)
        self.send_header("Content-type","text/html")
        self.send_header("test","This is test!")
        self.end_headers()
        buf = '''''<!DOCTYPE HTML>
                <html>
                <head><title>Get page</title></head>
                <body>

                <form action="post_page" method="post">
                  mobile: <input type="text" name="tos" /><br />
                  content: <input type="text" name="content" /><br />
                  <input type="submit" value="POST" />
                </form>

                </body>
                </html>'''
        self.wfile.write(buf)

    def do_POST(self):
        path = self.path
        print path
        #获取post提交的数据
        datas = self.rfile.read(int(self.headers['content-length']))
        datas = urllib.unquote(datas).decode("utf-8", 'ignore')
        mobile = datas.split('&')[0].replace("tos=","").replace('"','')
        content = datas.split('&')[1].replace("content=","").replace('"','')
        mm = mobile.split(',')
        for i in mm:
                print i
                send_message(i,content)
                self.send_response(200)
                self.end_headers()

                buf = '''''<!DOCTYPE HTML>
                <html>
                <head><title>Post page</title></head>
                <body>Tos:%s  <br />Content:%s</body>
                </html>'''% (mobile, content)
                self.wfile.write(buf)


def start_server(port):
    http_server = HTTPServer(('', int(port)), ServerHTTP)
    http_server.serve_forever()

if __name__ == "__main__":
    #端口可自定义，不冲突就可以，这里设置为：4040
    start_server(4040)
