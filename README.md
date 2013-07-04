HOLMES
======================================
--------------------------------------
## Background

There is Three main **Web Traffic** nowadays.

+ From User
+ From Normal Robots (googlebot, msnbot, etc.)
+ From Abnormal Web Crawler

We should give User quickly response with high priority, give Normal Robots right response with lower priority and reject serve for Abnormal Web Crawler.


+ Currently, we detect by implicit human browsing behavior. A Javascript is embedded into the pages served to the client dynamically. And a event handler for mouse movement or key clicks is included. Robots and Crawlers do not execute the Javascript. But some people disable the Javasript in their browser, or other reasons will lead to the method failed.

+ We want to analyse the Web Server **Access Log** for more accurate result.

---------------------------

## What

Our goal is to distinguish **Human User** from **Robots** by analyse the **Access Log**.

That means:

+ Input: Every single record r (belong to) a set of record R.
+ Ouput: Record with tag.

Example:

+ single input with many information:
 + 23	2013	28	59	103.0	103.0	xxx.xxx.xxx.xxx	xxx.xxx.xxx.xxx	XXX.XXX.com	GET	/abc/def/ghi	200	295	http://abc.def.ghi.com	Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/21.0.1180.89 Safari/537.1	0.03	-	xxx.xx.xxx.xxx	-	59	06	1482	80
+ single output with tag:
 + (xxx.xxx.xxx.xxx, human)
 + or (xxx.xxx.xxx.xxx, robot)

------------------------------------------

## How

There is three main method(Paper: [Web robot detection techniques: overview and limitations](http://iaumjournals.com/library/upload/article/af_2229553365_b.pdf)):

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

For more research.

---------------------------------------

## Deployment

Left Blank

----------------------------------------

## Milestone

Left Blank

----------------------------------------
