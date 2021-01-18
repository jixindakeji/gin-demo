# 安装步骤
### 数据库部署
1. 先建立数据库名称，修改配置文件处的数据库配置，启动程序，如果无误会自动创建表
2. 将db/init.sql 的数据导入到数据库对应的表，这是基础配置项，导入对应的用户/角色/权限等
3. admin用户修改密码时，需要配置将app.yaml中的 admin-secret，会重置密码，然后将该项删除，即可重置密码，如果存在该项，会一直重置密码
4. 配置对应的vue前端, 对应的版本vue前端需要对应master分支的tag一致，比如 0.1.* 对应的tag是 0.1.* (*忽略)
5. 对应的 vue web: https://github.com/jixindakeji/vue-demo.git


## 
if any bug found, please let me known, thanks for your sharing :)
