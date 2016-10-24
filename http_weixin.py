#coding=UTF-8
import BaseHTTPServer
import urlparse
import time
import re
from SocketServer import ThreadingMixIn
import threading
import sys
import urllib
import urllib2
import time
import json
import simplejson
from optparse import OptionParser

reload(sys)
sys.setdefaultencoding('utf-8')


class Token(object):
    def __init__(self, corpid, corpsecret):
        self.baseurl = 'https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid={0}&corpsecret={1}'.format(corpid, corpsecret)
        self.expire_time = sys.maxint
       # print "bbb"
       ## request = urllib2.Request(self.baseurl)
       # response = urllib2.urlopen(request)
       # print response.read()
       #get_token(self)

    def get_token(self):
        if self.expire_time > time.time():
            request = urllib2.Request(self.baseurl)
            response = urllib2.urlopen(request)
            ret = response.read().strip()
            ret = json.loads(ret)
            if 'errcode' in ret.keys():
                print >> ret['errmsg'],sys.stderr
                sys.exit(1)
            self.expire_time = time.time() + ret['expires_in']
            self.access_token = ret['access_token']
        return self.access_token

result= Token('xxxxxxxxxxxxx','xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx')

def send_message(wx_user,text_text):
#微信api消息类型:"text消息" （http://qydev.weixin.qq.com/wiki/index.php?title=%E6%B6%88%E6%81%AF%E7%B1%BB%E5%9E%8B%E5%8F%8A%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F）
    data = {}
    data['touser'] = wx_user
    data['toparty'] = 'text'
    data['agentid'] = 7
    data['msgtype'] = 'text'
    data['safe']="safe"
    data['totag']=""
    data['text'] = {"content": text_text}

    data = simplejson.dumps(data,ensure_ascii=False)
    req = urllib2.Request('https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s'%result.get_token())
    content = urllib2.urlopen(req,data).read()
    print content


class WebRequestHandler(BaseHTTPServer.BaseHTTPRequestHandler):

    def do_POST(self):
        path = self.path
        print path
        #获取post提交的数据
        datas = self.rfile.read(int(self.headers['content-length']))
        datas = urllib.unquote(datas).decode("utf-8", 'ignore')
        im = datas.split('&')[0].replace("tos=","").replace('"','')
        content = datas.split('&')[1].replace("content=","").replace('"','')
        mm = im.split(',')
        for i in mm:
                print i
                send_message(i,content)
                self.send_response(200)
                self.end_headers()

                buf = '''''<!DOCTYPE HTML>
                <html>
                <head><title>Post page</title></head>
                <body>Tos:%s  <br />Content:%s</body>
                </html>'''% (im, content)
                self.wfile.write(buf)

class ThreadingHttpServer( ThreadingMixIn, BaseHTTPServer.HTTPServer ):
    pass

if __name__ == '__main__':
    #服务端口4041，可自行修改，冲突就可以了
    server = ThreadingHttpServer(('0.0.0.0',4041), WebRequestHandler)
    ip, port = server.server_address
    # Start a thread with the server -- that thread will then start one
    # more thread for each request
    server_thread = threading.Thread(target=server.serve_forever)
    # Exit the server thread when the main thread terminates
    server_thread.setDaemon(True)
    server_thread.start()
    print "Server loop running in thread:", server_thread.getName()
    while True:
        pass
