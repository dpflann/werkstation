import random
import urllib2

def loop_it(i, inc, mod):
  while True:
     #urllib2.urlopen("http://localhost/arduino/analog/6/%i" % i)
     print("i = %d" % i)
     i += inc
     i %= mod

