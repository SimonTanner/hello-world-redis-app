
# hello-world-redis-app

## Installation & use

### Macos

Using homebrew, simply enter `brew install redis` in the command line.

In order to use redis as a cache & not just an in memory db, create a new folder using sudo mkdir /etc/redis and then copy the redis.conf file:
 `cp /usr/local/etc/redis.conf /etc/redis/redis.conf`
 Then open the new redis.conf file and find where it says `daemonize no` and set this to yes to run redis as a daemon. More advanced info on configuration can be found at https://www.miarec.com/doc/administration-guide/doc735 however a lot of this is relevant for linux.

To run the server simply enter `redis-server /etc/redis/redis.conf` which will start the server on localhost:6379 which is the default address & port number on a mac.

To test the connection you can either simply run the hello-world-redis-app which will ping the server, or simply use the redis-cli tool which will have automatically have been installed. Then simply enter ping to which you should get PONG as a response.

to shutdown the redis server enter `redis-cli shutdown`

### Windows 7