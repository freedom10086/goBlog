# GO博客系统
> 一个简单的go语言博客, 高效简洁，但不保证美观
## 配置psql
1. 设置并启动psql
    1. windows
    ```
    .\initdb.exe -D C:\app\pgsql-10.0.1\data -E UTF8
    .\pg_ctl -D "C:\app\pgsql-10.0.1\data" -l logfile start
    ```
    2. linux
    ```
    sudo apt-get install postgresql-client
    sudo apt-get install postgresql
    sudo -i -u postgres //切换到postgress用户
    psql//使用psql命令登录PostgreSQL控制台，以postgress用户
    \password postgres //更改密码
    ```
2. 创建用户 超级用户要谨慎
`createuser [--superuser] yang` or `createuser --interactive`

3. 设置密码 用postgres超级管理员登陆
> 为刚刚创建的用户设置密码

`psql postgres`

`\password yang`

`\q`

4. 创建数据库 指定所有权yang
`createdb -O yang -p 5432 blog` or `createdb blog` `grant all privileges on database blog to yang;`

5. 登陆数据库
`psql -U yang -d blog -h 127.0.0.1 -p 5432`

> -U 用户
> -d 数据库
> -h host
> -p 端口

`\q`

`exit`

6. 远程访问
```
$ find / -name "postgresql.conf" 或者 whereis postgresql.conf
/etc/postgresql/10/main/postgresql.conf

change 
listen_addresses = 'localhost'
to
listen_addresses = '*'

修改pg_hba.conf 添加
host    all(或者blog)     all(或者 yang)        0.0.0.0/0     md5(或者scram-sha-256)
host    all(或者blog)     all(或者 yang)        ::/0          md5(或者scram-sha-256)
```

7. 更新
`service postgresql reload` or `pg_ctl reload`
## 数据表定义
[goWebDB.sql](gowebDB.sql)

## 依赖及参考
- [Chart.js](http://chartjs.org/)
- [highlight.js](http://git.io/hljslicense)
- [bootstrap](https://getbootstrap.com)
- [marked.js](https://github.com/chjj/marked)

## bug 反馈
- [2351386755@qq.com](mailto://2351386755@qq.com)    
- [yangluo.chn@gmail.com](mailto://yangluo.chn@gmail.com)

## License
```
免费使用随便传播
```
