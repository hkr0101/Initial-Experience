# 简单的代码说明
实现了账户的注册、登录、登出，其中在储存密码时运用了简单的哈希函数。给予了admin
账号足够的权限。  
实现了在登录情况下添加、删除、修改、查看自己的问题，以
及在所有的情况下查看所有/特定问题。在登录情况下添加、删除、修改、查看自
己的答案，以及在所有的情况下查看某一个问题的答案  
一个小翻页，默认在显示一系列答案或者问题时每页20条内容  
在github上找到了一个关于调用chatgpt的项目用于生成ai答案，但是由于
我没有国外的手机号，无法获得chatgpt的key，这个内容仅仅停留在未测试可行性  
* main.go是主程序  
* routes中的是操作中涉及的函数
* mymodels中是三个实体Question、User、Answer
* myauth中是登录与登出的操作
* db中的是连接数据库以及在数据库中自动生成实体
* AI_answer中便是前文中提到的尚未完成的ai生成答案部分
api文档：https://apifox.com/apidoc/shared-86117e10-c314-4e57-a13f-494295b93689