HOLMES
======================================
--------------------------------------
## Section 1: Background

### 1.1 Web Traffic

There is Three main **Web Traffic** nowadays.

+ From **User**
+ From Normal Robots (googlebot, msnbot, etc.)
+ From Abnormal Web Crawler

As for User, they can classify into two type

+ **Normal User Clicks**
+ Malicious Clicks (Compete by Customers)

### 1.2 For Anjuke

+ User: **Quickly** response speed with **high** priority
 + Normal User: **Effective Click**
 + Malicious Click 
+ Normal Robots: Reasonable response speed with lower priority
+ Abnormal Web Crawler: Reject serve.

### 1.3 Difference between Previous Method

+ Javascript Code Counter
 + Currently, we detect by implicit human browsing behavior. A Javascript is embedded into the pages served to the client dynamically. And a event handler for mouse movement or key clicks is included. Robots and Crawlers do not execute the Javascript. But some people disable the Javasript in their browser, or other reasons will lead to the method failed.

+ We want to analyse the Web Server **Access Log** for more accurate result. Which means:
 + Complete history
 + More accuracy

### 1.4 Significant

+ Accurate effective count (money XD)
+ Forbid competitor robots
+ Protect customer (Malicious Click)

---------------------------

## Section 2: Input & Output

### 2.1 Format

Our method is to analyse the **Access Log**.

That means:

+ Input: Every single record r which belong to a set of record R.
+ Ouput: Single Record with tag.(Currently only effective click records)

### 2.2 Example

+ single input with many information:

 + 23	2013	28	59	103.0	103.0	xxx.xxx.xxx.xxx	xxx.xxx.xxx.xxx	XXX.XXX.com	GET	/abc/def/ghi	200	295	http://abc.def.ghi.com	Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/21.0.1180.89 Safari/537.1	0.03	-	xxx.xx.xxx.xxx	-	59	06	1482	80
+ single output with tag:
 + (xxx.xxx.xxx.xxx, human)

------------------------------------------

## Section 3: How

There is three main method

+ Syntactic log analysis
 + individual field parsing
 + user-agent mapping
 + multifaceted log analysis
+ Traffic pattern analysis
 + syntactic and pattern analysis
 + resource request patterns
 + query rate patterns
 + traffic metrics
+ Analytical learning
 + decision trees
 + neural networks
 + Bayesian network
 + Hidden Markov model

--------------------------------------------

## Section 4: Our implement

![](doc/image/implement.png) 

### 4.1 recognize web crawler by key words(such as "bot"、"spider")

Usually，some large legal search engine companies will declare themselves as a web crawler by give an user agent string which contain key words such as **bot** and **spider**. 

+ For example
 + Google: **Mozilla/5.0 (compatible; Googlebot/2.1; +http ://www.google.com/bot.html)** 
 + Baidu: **Mozilla/5.0 (compatible; Baiduspider/2.0; +http ://www.baidu.com/search/spider.html)** 

Based on the above facts,some normal web crawlers can be recognized using key words.

### 4.2 recognize web brower by normal user agent string pattern

In section 4.1, a method used to recognize normal web crawlers is described, however, in real world, there are so many web crawlers which do not declare themselves using user agent string which contained key words such as **bot** and **spider**. 

So, on the other hand, we focus on the normal browsers. Compared to web crawler, the amount of normal browser is small, and their user agent string pattern is more stable. 

+ For example
 + Google Chrome web browser: **Mozilla/5.0 (X11; Linux i686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1500.71 Safari/537.36**
 + So, user agent string pattern can be used to recognize normal browser as they declared. 
 + But,wait, some web crawlers also declare them as normal browser. By combine above two rules, a lot of web crawlers can be recognized.

### 4.3 recognize specific request by URI pattern

+ Among the access logs, there are many kinds of requests, but only some specific requests will be cared by us. 
+ These specific requests can be recognized using some uri patterns. 
+ For example, if we care request URI with **/a/b/** as the prefix, we can use regular pattern **^/a/b/** to match specific request URI.

### 4.4 recognize OK request by HTTP code 200

+ When client send request to server, the server may send response with different HTTP code which is depends on many condition. 
+ Among many HTTP code, the HTPP code **200** represents the response is OK. We use HTTP code **200** to recognize OK request.

--------------------------------------------

## Section 5: Result

### Classify Tree

![](doc/image/CountTree.png)


+ UA Filter 

![](doc/image/ua.png)

+ HttpCode Filter 

![](doc/image/httpcode.png) 
![](doc/image/4xx.png)

+ vppv 

![](doc/image/vppv.png)

+ Processed 

![](doc/image/processed.png)

### Compare Graph

![](doc/image/result.png)

### Analysis

+ The result is not stable at the begenning. 
 + The reason is that we do not have the history of the logs. 
 + After 1 hour, the result is a little higher than knowing.
 + The rule is not restrict.
+ One interesting thing is that in the watching list, there is some abvious bad spider.
 + Totally 11.7w watching list records
 + Top3: 3.3w, 2.3w, 2.1w

---------------------------------------------

## Section 6: What is more

+ More Rules
 + We use s.anjuke as the Rule to add into Whitelist, it is too restrict
 + In process, we only use referer filter, it is too loose
+ Better Algorithm
 + Watching list records should be analyse by there behavior to determine whether is effective or not.
+ Cluster for high speed
 + Currently, the speed of Holmes is still limited. We should use cluster to raise the speed of the program.
+ Other thing
 + Bad spider and bad human should be banned or challenged.
 + Bad spider should be banned, for protecting data and saving resource.

---------------------------------------------

## Section 7: Conclude

Currently the result is a little higher than knowing.

Probably reason:
 + The javascript code may not work properly.
 + Our data is the complete record.

There is still two thing we can do:
 + For those records we think is effective, there might be mang malicious click.
 + For those we kept in the watching list, classify the effective clicks.

---------------------------------------------
