simple-quote-monitor
===========================
A simple command line Go tool for monitoring stock quotes and sending desktop alerts when price drops or rises to certain amount


How it works
------------
Tired keeping track of the stock quotes ? Let simple-quote-monitor come rescue. It will periodically check the quote for
every 30 to 40 seconds and send desktop alerts when price drops or rises to certain amount.
![Desktop alert](../master/example.png?raw=true)

It scrapes google finance website and finds the price on it. Unfortunately, I didn't find great finance API and google also deprecated theirs so I am doing scraping instead. **However, when you use this tool, please keep in mind that the data from google finance are NOT real-time in all cases. Some stocks have real-time while many others have around 15 mins delay.** When you set the range of price you are watching, please consider this 15 mins delay. This might be addressed in the future if I find providers who have real-time stock quotes.

Dependencies
------------
Currently, the notification only works on OSX Mountain Lion's or after with the notification center
and you have to brew install [terminal-notifier](https://github.com/julienXX/terminal-notifier) for desktop alerts. Simply
run
```
$ brew install terminal-notifier
```
and you are done.
**This will be addressed in the future as I might embed the terminal-notifier in the source code so you do
not have to install terminal-notifier by yourself. I will also add supports to Windows and Linux**

Usage
-----
Install this tool with
```
go get github.com/ssdong/simple-quote-monitor
```

The command line options are
```
Usage: simple-quote-monitor [options]

-se <stock exchange>     The code of stock exchange
-ss <stock symbol>       The unique series of letters of a security
-min <number>            The minimum price watching for
-max <number>            The maximum price watching for
```

E.g. run
```
simple-quote-monitor -se=cve -ss=icc -min=0.5 -max=1
```
will monitor stock ICC in the stock exchange CVE with a minimum price of $0.5 and maximum of $1. It will
send you notifications when the price falls outside of this range(including these two numbers).

Tips
----
If you want the desktop alert to stay until you clicks it you would have to change your notification preferences for
terminal-notifier in notification centre as this [one](https://www.dropbox.com/s/n2kt0in8q6syiu6/Screenshot%202016-10-11%2014.27.27.png?dl=0)
Otherwise, the notification will go away in 5 seconds and you might miss it. Unfortunately this is something you have to
do it manually since OSX forces the notification to disappear within a certain amount of time.
