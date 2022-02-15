# quant_web_server_
a demo tool of showing mirco market of virture coin

这是一个用于显示虚拟货币微观市场结构的demo
url基于web前端，需要配置mysql和一些其他参数，具体请见配置文件

你目前可以用这个demo完成以下事情：
  1、获取Bianca实时行情信息（有限深度订单簿，归集交易，最新最优报价，卖方买房数量等），存储在本地mysql
  2、对于这些基本信息的预处理以及对订单簿特征的提取
  3、对某个节点的可视化，包括显示最优买房卖房报价的气泡图，订单簿具体信息，订单簿变化，汇集交易行情均价
 
这些功能都在本地前端（默认127.0.0.1:8080），可在配置文件更改）实现
 
关注此项目，后续会完成：
  1、okex实时行情
  2、合约行情
  3、更友善的微观结构可视化操作
  4、基于okex的实盘交易接口

所有代码均个人所写，可以放心直接使用，有问题加vx：13997171940交流

This is a demo for displaying the micro market structure of virtual money
The URL is based on the web front end. You need to configure MySQL and some other parameters. See the configuration file for details

You can now use this demo to do the following:
1. Obtain Bianca real-time market information (limited depth order book, collection transaction, latest and best quotation, number of houses bought by the seller, etc.) and store it in local mysql
2. For the preprocessing of these basic information and the extraction of order book features
3. The visualization of a node includes the bubble chart showing the optimal quotation for buying and selling houses, the specific information of the order book, the change of the order book, and the collection of the average price of the transaction market

These functions are implemented in the local front end (127.0.0.1:8080 by default) and can be changed in the configuration file

Pay attention to this project, and the following will be completed:
1. Okex real time quotes
2. Contract price
3. More friendly microstructure visualization
4. Interface of firm offer transaction based on okex

All codes are written by individuals and can be used directly. If there is a problem, add VX: 13997171940 for communication
