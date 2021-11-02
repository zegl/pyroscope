#!/usr/bin/env python3

import os

def work(n):
	i = 0
	while i < n:
		i += 1

def fast_function():
	work(20000)

def slow_function():
	print("test print")
	work(800000)

if __name__ == "__main__":
	while True:
		fast_function()
		slow_function()
