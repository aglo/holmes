library("rjson")
# Use rjson to get result
# JSON -> list
knowing <- fromJSON(file="http://                                  lse",method="C")
total_request <- fromJSON(file="h+--------------------------------+per_min",method="C")
ua_no_pass <- fromJSON(file="http|                                |s_per_min",method="C")
ua_pass <- fromJSON(file="http://|       +------------+           |min",method="C")
total_vppv <- fromJSON(file="http|           +        |           |n",method="C")
code200 <- fromJSON(file="http://|           |        |           |n",method="C")
code000 <- fromJSON(file="http://|           |        |           |n",method="C")
code301 <- fromJSON(file="http://|           |        |           |n",method="C")
code302 <- fromJSON(file="http://|           |        +           |n",method="C")
code400 <- fromJSON(file="http://|           +------------+       |n",method="C")
code403 <- fromJSON(file="http://|                        +       |n",method="C")
code404 <- fromJSON(file="http://|                        |       |n",method="C")
code408 <- fromJSON(file="http://|     +-------------+    |       |n",method="C")
code499 <- fromJSON(file="http://|                        |       |n",method="C")
code503 <- fromJSON(file="http://|                   +-+  |       |n",method="C")
watching <- fromJSON(file="http:/|                   | |  +       |in",method="C")
effective <- fromJSON(file="http:|                   +-++-+       |_min",method="C")
no_referer <- fromJSON(file="http+--------------------------------+er_min",method="C")
from_my <- fromJSON(file="http://                                  n",method="C")

# Split knowing list, change it into vector
tmp1 <- sapply(knowing,"[[",1)
tmp1 <- unlist(tmp1)/1000

tmp2 <- sapply(knowing,"[[",2)
tmp2[unlist(sapply(tmp2,is.null))] <- NA
tmp2 <- unlist(tmp2)

knowing <- cbind(tmp1,tmp2)
rm(tmp1,tmp2)

# function to process holmes result
p_holmes <- function(holmes) {
    tmp1 <- sapply(holmes,"[[",2)
    tmp1 <- unlist(tmp1)

    tmp2 <- sapply(holmes,"[[",3)
    tmp2[unlist(sapply(tmp2,is.null))] <- NA
    tmp2 <- unlist(tmp2)

    holmes <- cbind(tmp1,tmp2)
    rm(tmp1,tmp2)
    holmes
}

# Draw png file
png("result.png",width=640,height=480,units="px",pointsize=16)   #,type = c("cairo","cairo-png","Xlib","quartz"),antialias="none")
plot(knowing,xlab="Unix Timestamp",ylab="Click Num",type ="l",ylim=c(0,2000))
lines(p_holmes(holmes=effective),col="red")
legend("top",legend = c("From Knowing","From Holmes"),col = c("black","red"),lty = 1)
title(main="2013-7-22 9:00~13:00")
dev.off()

# Draw Pie graph (UA)
png("ua.png",width=640,height=480,units="px",pointsize=16)
ans <- c( sum(p_holmes(ua_no_pass)[,2]), sum((p_holmes(ua_pass)[,2])))
names(ans) <- c("UA no pass" , "UA pass")
pie(ans,main="UA Filter")
dev.off()


# Draw Pie graph (HTTP CODE)
png("httpcode.png",width=640,height=480,units="px",pointsize=16)
ans <- c( sum(p_holmes(code000)[,2]), sum((p_holmes(code200)[,2])) ,
        sum(sum(p_holmes(code301)[,2]) ,sum(p_holmes(code302)[,2])),
        sum(sum(p_holmes(code400)[,2]) ,sum(p_holmes(code403)[,2]) ,sum(p_holmes(code404)[,2]) ,sum(p_holmes(code408)[,2]) ,sum(p_holmes(code499)[,2])),
        sum(p_holmes(code503)[,2]))
names(ans) <- c("000" , "200" , "3xx", "4xx"  , "503")
pie(ans,main="HTTPCODE Filter")
dev.off()

# Draw Pie graph ( 4xx )
png("4xx.png",width=640,height=480,units="px",pointsize=16)
ans <- c(sum(p_holmes(code400)[,2]) ,sum(p_holmes(code403)[,2]) ,sum(p_holmes(code404)[,2]) ,sum(p_holmes(code408)[,2]) ,sum(p_holmes(code499)[,2]))
names(ans) <- c("400" , "403" , "404" , "408" , "499")
pie(ans,main="HTTPCODE = 4xx")
dev.off()

# Draw Pie graph ( vppv )
png("vppv.png",width=640,height=480,units="px",pointsize=16)
ans <- c(sum(p_holmes(watching)[,2]) , sum ( sum(p_holmes(effective)[,2]) , sum(p_holmes(from_my)[,2]) , sum(p_holmes(no_referer)[,2])))
names(ans) <- c("watching list","processed")
pie(ans ,main = "Analysis")
dev.off()

# Draw Pie graph ( Processed )
png("processed.png",width=640,height=480,units="px",pointsize=16)
ans <- c(sum(p_holmes(effective)[,2]) , sum(p_holmes(from_my)[,2]) , sum(p_holmes(no_referer)[,2]))
names(ans) <- c("Effective Click" , "From my.anjuke" , "No Referer")
pie(ans , main = "Processed")
dev.off()

