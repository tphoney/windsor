import urllib, urllib2, cookielib, datetime, sys, time
from optparse import OptionParser
#crontab example
#0 0 * * 1 python /home/tp/windsor.py
month = "04"
date = "22"
hour = "20"
now = datetime.datetime.now()
user = "joe bloggs"
#rooms = "28"  #dome 5
#participant = "emma+duggan"
#rooms = "27"   #dome 4

parser = OptionParser()
parser.add_option("-p", "--participant", dest="participant", default="paul+gibson", help="participant name, eg: joe+bloggs")
parser.add_option("-d", "--dayofweek", dest="dayofweek", default=0, type="int", help="day of week, ie 0 = Monday, 1=Tuesday, 2=Wednesday...")
(options, args) = parser.parse_args()
#time.sleep(9)

def goGetEm (participant, month, date, hour, rooms):
	print "time now", datetime.datetime.now()
	print "participant", participant
	print "month", "%02d" % month
	print "date", "%02d" % date
	print "hour", hour
	print "rooms", rooms

	cookie_jar = cookielib.CookieJar()
	opener = urllib2.build_opener(urllib2.HTTPCookieProcessor(cookie_jar))
	urllib2.install_opener(opener)

	# do POST
	url_1 = 'http://www.windsortennis.co.uk/courtbooker/day.php?day=07&month=03&year=2015&area=0&room=3&returl=http%3A%2F%2Fwww.windsortennis.co.uk%2Fcourtbooker%2Fday.php%3Fday%3D07%26month%3D03%26year%3D2015%26area%3D0%26room%3D3'
	values = dict(NewUserName=str(user), NewUserPassword='password', returl='http://www.windsortennis.co.uk/courtbooker/day.php?day=07&month=03&year=2015&area=0&room=3&returl=http://www.windsortennis.co.uk/courtbooker/day.php?day=07&month=03&year=2015&area=0&room=3', TargetURL='day.php?day=07&month=03&year=2015&area=0&room=3&returl=http%3A%2F%2Fwww.windsortennis.co.uk%2Fcourtbooker%2Fday.php%3Fday%3D07%26month%3D03%26year%3D2015%26area%3D0%26room%3D3', Action='SetName')
	data = urllib.urlencode(values)
	req = urllib2.Request(url_1, data)
	rsp = urllib2.urlopen(req)
	content = rsp.read()
	if "You are" in content:
		print ("Logged in")
		#book the damn thing
		url_book = 'http://www.windsortennis.co.uk/courtbooker/edit_entry_handler.php'
		values = dict(name=str(user), 
			participant_1= str(participant),
			participant_2="",participant_3="",
			participant_4="", participant_5="", participant_6="", participants_others="", 
			description="",
			start_day= str(date), 
			start_month= str(month) ,
			start_year=2018,
			start_seconds=72000,
			end_day= str(date), 
			end_month= str(month) ,
			end_year=2018,
			end_seconds=72000,
			area=13,
			rooms=str(rooms),
			type='A',
			returl='http://www.windsortennis.co.uk/courtbooker/day.php?year=2015&month=03&day=9&area=13&room=27',
			create_by=str(user),
			rep_id=0,
			edit_type='series')

		data = urllib.urlencode(values)
		req = urllib2.Request(url_book, data)
		rsp = urllib2.urlopen(req)
		content = rsp.read()
	return content

def next_weekday(d, weekday):
	days_ahead = weekday - d.weekday()
	if days_ahead <= 0: # Target day already happened this week
		days_ahead += 14
	return d + datetime.timedelta(days_ahead)

def is_success(content):
	returnVal = True
	if 'conflict' in content:
		returnVal = False
	return returnVal
		

#test
#goGetEm ("paul+gibson", "8", "10", "15", "27");
print "today is month:", now.month, "day:", now.day
next_date = next_weekday(now, options.dayofweek)
print "next date:", next_date

success = False
pageContent = "nothing"
for i in range (0, 80):
		pageContent = goGetEm (options.participant, next_date.month, next_date.day, hour, "28");
		if is_success(pageContent):
			print "success 28!!!"
			break;
		pageContent = goGetEm (options.participant, next_date.month, next_date.day, hour, "27");
		if is_success(pageContent):
			print "success 27!!!"
			break;
		print "wait 1"
		time.sleep(1)
print "Done."
