# Aranea
BUG :
+ 因为之前一直在开着ssr的状态下调试的，偶然关了ssr之后发现直接跑不起来了，估计是在node选择ip的时候出错了.

TODO:

+ ~~项目文件结构~~
+ ~~爬虫~~
+ ~~doing: master node 心跳~~
+ doing: master 接受/发送任务，收集返回内容,写到文件中.
+ 任务队列

***
+ cronJob
+ urlDelay : node连续爬取某个url需要等待的时间
+ HTTP to RPC
+ 支持任务中定义对于爬取后的数据进行什么处理之后再返回
