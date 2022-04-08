#!/usr/bin/env python3

import os
import time
import math


import pyroscope

pyroscope.configure(
	app_name       = "broken.python.app",
	server_address = "http://pyroscope:4040",
	tags           = {
    "forked": "false",
	}
)

def work(label, seconds):
	i = 0
	st = int(time.time())
	last_tag_update = -1
	while st+seconds > int(time.time()):
		second = math.floor((int(time.time()) - st))
		if last_tag_update != second:
			print(f'updating tags for {label}, second {second}')
			# I call it a 500 times because it doesn't hang consistently enough
			for i in range(500):
				pyroscope.tag({ "second": str(second) })
			last_tag_update = second
		i += 1

def forked_process():
  print("forked_process work start")
  work("forked_process", 15)
  print("forked_process work end")

def old_process():
  print("old_process work start")
  work("old_process", 15)
  print("old_process work end")

if __name__ == "__main__":
	# print multiline string
	print("""
	this program should:
	* create a fork of itself
	* old process will run function `bar` for 15 seconds
	* forked process will change its tags to 'forked => true'
	* forked process will run function `foo` for 15 seconds
	* both will add tags { "second" => X } for every second X
	* results in pyroscope should show 15 seconds for old process and 15 seconds for new process
	* there should be no double-counting
	---

	currently the there's a couple of problems:
	* forked process eventually hangs when you call Pyroscope.tag enough times — big issue
	* the data from forked process doesn't show up in pyroscope — this is a smaller issue (but we should still fix it).
		It's okay if users have to do extra set up work to make it work (e.g add Pyroscope.configure in `fork` method block)
	""")

	newpid = os.fork()
	if newpid == 0:
		forked_process()
	else:
		old_process()
		os.wait()

