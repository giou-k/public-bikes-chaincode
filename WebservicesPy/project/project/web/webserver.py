# -*- coding: utf-8 -*-
import tornado.auth
import tornado.escape
import tornado.httpserver
import tornado.ioloop
import tornado.options
import tornado.web
import Settings
import sqlite3
import json
import bcrypt
import hashlib
import re
import os, uuid
import csv
from tornado.options import define, options
from tornado_cors import CorsMixin


__UPLOADS__ = "uploads/"
define("port", default=8888, help="run on the given port", type=int)

class BaseHandler(CorsMixin,tornado.web.RequestHandler):
    CORS_ORIGIN = '*'
    CORS_HEADERS = 'Content-Type'
    CORS_METHODS = 'POST', 'GET'
    CORS_CREDENTIALS = True

    def set_default_headers(self):
        print "setting headers!!!"
        self.set_header("Access-Control-Allow-Origin", "*")
        self.set_header("Access-Control-Allow-Headers", "x-requested-with")
        self.set_header('Access-Control-Allow-Methods', 'POST, GET, OPTIONS')
        # Header always set Access-Control-Max-Age "1000"
        # Header always set Access-Control-Allow-Headers "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding"
        # Header always set Access-Control-Allow-Methods "POST, GET, OPTIONS, DELETE, PUT"
    def get_current_user(self):
        return self.get_secure_cookie("user")



class MainHandler(BaseHandler):
    @tornado.web.authenticated
    def get(self):
        conn = sqlite3.connect('dbs/database.db')
        c = conn.cursor()
        username = self.current_user
        #username = username.replace('""', '')
        query = "SELECT apikey FROM users WHERE mail = " + username
        apikey = c.execute(query).fetchone()
        if apikey[0] == '64e947aadea5d17b9da8a7a8194e13175b443cd37e2dc5a8669a9302':
            self.render("index.html")
        else:
            self.write('<html><body> <div> User: '+ username +'</div> <br> <div> My Apikey: '+ apikey[0] +' </div> <br> <div> Stats: </div> </body></html>')


class AuthLoginHandler(BaseHandler):
    def post(self):
        data = json.loads(self.request.body)
        print "POST received: ", data
        username = data['email']
        password = data['password']
        password = hashlib.sha224(password).hexdigest()
        auth = self.check_permission(password, username)
        if auth:
            print('skata')
            self.set_current_user(username)
            self.get_data(usename)
        else:
            self.set_status(400)

    def check_permission(self, password, username):
        conn = sqlite3.connect('dbs/database.db')
        c = conn.cursor()
        tmp_usr = c.execute("SELECT mail FROM users WHERE mail = ?", (username,)).fetchone()
        tmp_pass = c.execute("SELECT password FROM users WHERE mail = ?", (username,)).fetchone()
        print tmp_usr
        print tmp_pass
        #hashed_pass = bcrypt.hashpw(password, bcrypt.gensalt())
        #hashed_pass = hashlib.sha224(password).hexdigest()
        #print tmp_pass[0]
        #print tmp_usr[0]
        #print password
        conn.commit()
        conn.close()
        if tmp_usr is not None and (tmp_pass[0] == password):
            # if username == tmp_usr[0] and password == tmp_pass[0]:
            return True
        return False

    def set_current_user(self, user):
        if user:
            self.set_secure_cookie("user", tornado.escape.json_encode(user))
        else:
            self.clear_cookie("user")

class UserHandler(BaseHandler):
    def get(self):
        username = self.get_argument("email", "/")
        print username
        conn = sqlite3.connect('dbs/database.db')
        c = conn.cursor()
        apikey = c.execute("SELECT apikey FROM users WHERE mail = ?", (username,)).fetchone()
        stations_requests = c.execute("SELECT stations_request FROM users WHERE mail = ?", (username,)).fetchone()
        rupos_requests = c.execute("SELECT rupos_request FROM users WHERE mail = ?", (username,)).fetchone()
        range_requests = c.execute("SELECT range_request FROM users WHERE mail = ?", (username,)).fetchone()
        conn.commit()
        conn.close()
        dict = {'User' +username, 'Apikey' +apikey, 'Stations requests' +stations_requests, 'Rupos requests' +rupos_requests, 'Range requests' +range_requests}
        print dict
        # list[username,apikey[0],stations_request[0],rupos_request[0],range_request[0]]
        # self.redirect('http://localhost:8888/api/?apikey='+apikey[0])
        # return username
        self.write(json.dumps(username,apikey[0],stations_request[0],rupos_request[0],range_request[0]))
        # return stations_request[0]
        # return rupos_request[0]
        # return range_request[0]


class AdminLoginHandler(BaseHandler):
    def post(self):
        data = json.loads(self.request.body)
        print "POST received: ", data
        username = data['email']
        password = data['password']
        password = hashlib.sha224(password).hexdigest()
        auth = self.check_admin_permission(password, username)
        if auth:
            print('skata')
        else:
            self.set_status(400)

    def check_admin_permission(self, password, username):
        conn = sqlite3.connect('dbs/database.db')
        c = conn.cursor()
        tmp_usr = c.execute("SELECT mail FROM users WHERE mail = ?", (username,)).fetchone()
        tmp_pass = c.execute("SELECT password FROM users WHERE mail = ?", (username,)).fetchone()
        conn.commit()
        conn.close()
        if tmp_usr[0] == 'anz@test.com' and tmp_pass[0] == '12345':
            return True
        return False


class AuthLogoutHandler(BaseHandler):
    def get(self):
        self.clear_cookie("user")
        self.redirect(self.get_argument("next", "/"))

class RegisterHandler(BaseHandler):
    def post(self):
        data = json.loads(self.request.body)
        print "POST received: ", data
        email = data['email']
        password = data['password']
        print email
        print password
        conn = sqlite3.connect('dbs/database.db')
        c = conn.cursor()
        c.execute("SELECT * FROM users WHERE mail = ?", (email,))
        if c.fetchone():
            error_msg = u"?error=" + tornado.escape.url_escape("Login name already taken")
            self.redirect(u"/register/" + error_msg)
        else:
            key = hashlib.sha224(email).hexdigest()
            hashed_pass = hashlib.sha224(password).hexdigest()
            #hashed_pass = bcrypt.hashpw(password, bcrypt.gensalt())
            c.execute("INSERT INTO USERS (mail, password, apikey) VALUES (?,?,?)", (email, hashed_pass, key,))
            conn.commit()
            conn.close()
        self.redirect("/auth/login/")

    #    user = {}
    #    user['user'] = email
    #    user['password'] = hashed_pass
    #    auth = self.application.syncdb['users'].save(user)
    #    self.set_current_user(email)

class StationInput(BaseHandler):
    def post(self):
        data = json.loads(self.request.body)
        print "POST received: ", data
        name = data['station']
        passname = data['passname']
        lat = data['lat']
        lng = data['lng']
        conn = sqlite3.connect('dbs/database.db')
        c = conn.cursor()
        c.execute("INSERT INTO STATIONS (name, passname, lat, long) VALUES (?,?,?,?)", (name, passname, lat, lng))
        conn.commit()
        conn.close()


class UploadForm(tornado.web.RequestHandler):
    def get(self):
        self.render("upload.html")


class UploadFile(tornado.web.RequestHandler):
    def post(self):
        conn = sqlite3.connect('dbs/database.db')
        c = conn.cursor()
        file1 = self.request.files['files[]'][0]
        station_name = self.get_argument('station_name', '')
        rupos = self.get_argument('rupos', '')
        print station_name
        original_fname = file1['filename']
        extension = os.path.splitext(original_fname)[1]
        final_filename = str(uuid.uuid4()) + extension
        output_file = open("uploads/" + final_filename, 'w')
        output_file.write(file1['body'])
        self.finish("file" + final_filename + "is uploaded")
        reader = csv.reader(open("uploads/" + final_filename, 'r'), delimiter=',')
        for row in reader:
            date = row[0].split('-')[2] + row[0].split('-')[1] + row[0].split('-')[0]
            print(date)
            print(row[2])
            # c.execute("INSERT INTO MEASUREMENTS (rupos,passname,dt) VALUES (?,?,?)", (rupos,station_name,date))
            # c.execute("INSERT INTO MEASUREMENTS (rupos, passname, dt, one, two, three, four, five, six, seven, eight, nine, ten, eleven, twelve, thirteen, fourteen) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", (rupos,station_name,row[0],row[1],row[2],row[3],row[4],row[5],row[6],row[7],row[8],row[9],row[10],row[11],row[12],row[13],row[14]))
            # c.execute("INSERT INTO MEASUREMENTS (five, six, seven, eight, nine, ten) VALUES (?,?,?,?,?,?)", (row[5],row[6],row[7],row[8],row[9],row[10]))
            # c.execute("INSERT INTO MEASUREMENTS (one,two) VALUES (?,?)", (row[1],row[2],))
        conn.commit()
        conn.close()


class Stations(BaseHandler):
    def get(self):
        conn = sqlite3.connect('dbs/database.db')
        c = conn.cursor()
        rows = c.execute("SELECT * FROM stations")
        data_list = []
        for row in rows:
            temp = {}
            temp['id'] = row[0]
            temp['name'] = row[1]
            temp['passname'] = row[2]
            temp['lat'] = row[3]
            temp['long'] = row[4]
            data_list.append(temp)
        conn.commit()
        conn.close()
        self.write(json.dumps(data_list))



class Users(BaseHandler):
    def get(self):
        conn = sqlite3.connect('dbs/database.db')
        c = conn.cursor()
        rows = c.execute("SELECT * FROM users")
        data_list = []
        for row in rows:
            temp = {}
            temp['id'] = row[0]
            temp['mail'] = row[1]
            temp['stations_request'] = row[4]
            temp['rupos_request'] = row[5]
            temp['range_request'] = row[6]
            data_list.append(temp)
        conn.commit()
        conn.close()
        self.write(json.dumps(data_list))


class ServiceHandler(tornado.web.RequestHandler):
    def check_api(self, apikey):
        conn = sqlite3.connect('dbs/database.db')
        c = conn.cursor()
        tmp_apikey = c.execute("SELECT apikey FROM users WHERE apikey = ?", (apikey,)).fetchone()
        conn.commit()
        conn.close()
        if tmp_apikey is not None:
            return True
        return False

    def check_passname(self, passname):
        conn = sqlite3.connect('dbs/database.db')
        c = conn.cursor()
        tmp_passname = c.execute("SELECT passname FROM stations WHERE passname = ?", (passname,)).fetchone()
        conn.commit()
        conn.close()
        if tmp_passname is not None:
            return True
        return False

    def check_rupos(self, rupos):
        conn = sqlite3.connect('dbs/database.db')
        c = conn.cursor()
        tmp_rupos = c.execute("SELECT rupos FROM measurements WHERE rupos = ?", (rupos,)).fetchone()
        conn.commit()
        conn.close()
        if tmp_rupos is not None:
            return True
        return False

    def check_date(self, date):
        conn = sqlite3.connect('dbs/database.db')
        c = conn.cursor()
        tmp_date = c.execute("SELECT dt FROM measurements WHERE dt = ?", (date,)).fetchone()
        conn.commit()
        conn.close()
        if tmp_date is not None:
            return True
        return False

    def check_time(self, time):
        # conn = sqlite3.connect('dbs/database.db')
        # c = conn.cursor()
        # tmp_time = c.execute("SELECT ? FROM measurements", (time,)).fetchone()
        # conn.commit()
        # conn.close()
        time_list3 = ["one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen",
        "seventeen", "eighteen", "nineteen", "twenty", "twentyone", "twentytwo", "twentythree", "twentyfour"]
        # time_list2 = ["eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen", "twenty", "twentyone", "twentytwo", "twentythree", "twentyfour"]
        # time_list = time_list1 + time_list2
        for i in range (0, 24):
            print("time_list = ", time_list3[i])
            if time == time_list3[i]:
                return True
        return False

# end of checks
# !!!

    def get_all_where(self, query, varwhere, json_str = False):
        conn = sqlite3.connect('dbs/database.db')
        conn.row_factory = sqlite3.Row # This enables column access by name: row['column_name']
        db = conn.cursor()
        rows = db.execute(query,(varwhere,)).fetchall()
        conn.commit()
        conn.close()
        if json_str:
            return json.dumps([dict(ix) for ix in rows], sort_keys=True, indent=4) #CREATE JSON
        return rows

    def get_all_rupos(self, query, varwhere, varupos, json_str = False):
        conn = sqlite3.connect('dbs/database.db')
        conn.row_factory = sqlite3.Row # This enables column access by name: row['column_name']
        db = conn.cursor()
        rows = db.execute(query,(varwhere, varupos,)).fetchall()
        conn.commit()
        conn.close()
        if json_str:
            return json.dumps([dict(ix) for ix in rows], sort_keys=True, indent=4) #CREATE JSON
        return rows

    def get_all_date(self, query, varwhere, varupos, vardate, json_str = False):
        conn = sqlite3.connect('dbs/database.db')
        conn.row_factory = sqlite3.Row # This enables column access by name: row['column_name']
        db = conn.cursor()
        rows = db.execute(query,(varwhere, varupos, vardate,)).fetchall()
        conn.commit()
        conn.close()
        if json_str:
            return json.dumps([dict(ix) for ix in rows], sort_keys=True, indent=4) #CREATE JSON
        return rows

    def get_all_time(self, query, varwhere, varupos, vardate, json_str = False):
        conn = sqlite3.connect('dbs/database.db')
        conn.row_factory = sqlite3.Row # This enables column access by name: row['column_name']
        db = conn.cursor()
        rows = db.execute(query,(varwhere, varupos, vardate,)).fetchall()
        print query
        print varwhere
        print varupos
        print vardate
        conn.commit()
        conn.close()
        if json_str:
            return json.dumps([dict(ix) for ix in rows], sort_keys=True, indent=4) #CREATE JSON
        return rows

    def get_all(self, query, json_str = False):
        conn = sqlite3.connect('dbs/database.db')
        conn.row_factory = sqlite3.Row # This enables column access by name: row['column_name']
        db = conn.cursor()
        rows = db.execute(query).fetchall()
        conn.commit()
        conn.close()
        if json_str:
            return json.dumps([dict(ix) for ix in rows], sort_keys=True, indent=4) #CREATE JSON
        return rows

    def get(self):
        apikey = self.get_argument("apikey", "/")
        auth_api = self.check_api(apikey)
        if auth_api==True:
            self.write("OK")
            passname = self.get_argument("passname", "/")
            rupos = self.get_argument("rupos", "/")
            date = self.get_argument("date", "/")
            time = self.get_argument("time", "/")
            auth_pass = self.check_passname(passname)
            auth_rupos = self.check_rupos(rupos)
            auth_date = self.check_date(date)
            auth_time = self.check_time(time)
            if passname == "stations":
                self.write(self.get_all('SELECT * FROM stations', json_str = True))
                self.write("all stations")
            elif auth_pass and auth_rupos==False:
                self.write(self.get_all_where("SELECT * FROM measurements WHERE passname =?", passname, json_str = True ))
            elif auth_rupos and auth_date==False:
                self.write(self.get_all_rupos("SELECT * FROM measurements WHERE passname = ?  AND rupos = ? ", passname, rupos, json_str = True))
            elif auth_time==False and auth_date==True:
                self.write("mphke")
                self.write(self.get_all_date("SELECT * FROM measurements WHERE passname = ? AND rupos = ? AND dt = ?", passname, rupos, date, json_str = True))
            elif auth_time and auth_date:
                self.write("mphke k edw")
                self.write(self.get_all_time("SELECT %s FROM measurements WHERE passname = ? AND rupos = ? AND dt = ? " % time, passname, rupos, date, json_str = True))
        else:
            self.write("Please enter a valid apikey")

            # self.write("ok")
            # self.write(passname)
            # self.write(date)
            # self.write(time)


class Application(tornado.web.Application):
    def __init__(self):
        handlers = [
            (r"/", MainHandler),
            (r"/auth/login/", AuthLoginHandler),
            (r"/auth/logout/", AuthLogoutHandler),
            (r"/api/", ServiceHandler),
            (r"/register/", RegisterHandler),
            (r"/stationinput/",StationInput),
            (r"/uploadform/",UploadForm),
            (r"/upload/",UploadFile),
            (r"/adminlogin/",AdminLoginHandler),
            (r"/user/",UserHandler),
            (r"/stations/",Stations),
            (r"/users/",Users),
            #(r"/(.*)", tornado.web.StaticFileHandler, {'path': path}),
            ]

        settings = {
           "template_path":Settings.TEMPLATE_PATH,
            "static_path":Settings.STATIC_PATH,
            "debug":Settings.DEBUG,
            "cookie_secret": Settings.COOKIE_SECRET,
            "login_url": "/auth/login/"
        }
        tornado.web.Application.__init__(self, handlers, **settings)

def main():
    tornado.options.parse_command_line()
    http_server = tornado.httpserver.HTTPServer(Application())
    http_server.listen(options.port)
    tornado.ioloop.IOLoop.instance().start()
    database_check()
    data = self.get_argument('body', 'No data received')
    print data

if __name__ == "__main__":
    main()
