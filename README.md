GohJasmin
======

![Goh Jasmin](https://raw.githubusercontent.com/LordEidi/gohjasmin/master/ohjasmin.jpg)

**GohJasmin** (c) 2017-20 by [SwordLord - the coding crew](https://www.swordlord.com/)

Based on [OhJasmin](https://sourceforge.net/projects/ohjasmindns/) (c) 2010-20 by [SwordLord - the coding crew](https://www.swordlord.com/)

![Build and release gohjasmin](https://github.com/swordlordcodingcrew/gohjasmin/workflows/Build%20and%20release%20gohjasmin/badge.svg?branch=master) ![CodeQL Analysis](https://github.com/swordlordcodingcrew/gohjasmin/workflows/CodeQL%20Analysis/badge.svg?branch=master)

## Introduction ##

**GohJasmin** is a lightweight, self-hostable dynamic DNS server written in Golang. Just add an instance of [PowerDNS](https://powerdns.com/) and you have your very own dynamic DNS service with your own domains. 

(You could actually add any other DNS server to the mix. We used to do so and used [djbdns](https://cr.yp.to/djbdns.html). You will just need to write some glue code to export the dynamic data from the GohJasmin database into your own DNS server. Thats not very difficult.)

If you are looking for a service to manage your dynamic dns records and have your hosts send their current IP addresses to, then **GohJasmin** might be for you.

## Installation ##

Get the latest release here:

https://github.com/swordlordcodingcrew/gohjasmin/releases

## Configuration ##

All parameters which can be configured right now are in the file *gohjasmin.config.js*. A default configuration file will be written on the first run of the application.

It is important to know, that you should not have the sqlite DB in the same directory as you have the config as well as the authentication file. **gohjasmind** watches said directories and reacts on file changes. Since sqlite DB will change when being written to, this would create update events which are unecessary and slow down **gohjasmind**.

## Authentication ##

Have a look at the bin/gohjasmin.auth file. It contains demo data on how to authenticate users. The file is only read when starting gohjasmind. The fields are:

- domain
- user
- password hash (bcrypt, you can get your hash using the command below)
- last ip seen (which currently does not work as expected)

### Generating Password Hash ###

To generate a brypted hash for your password authentication, use this command. You will then be prompted for the password twice.

    htpasswd -nB user

## Run as a systemd service ##

Have a look at the bin/gohjasmind.service file.

The commands to install and run **gohjasmind** as a systemd service are:

    cp gohjasmin.service /lib/systemd/system/.
    chmod 755 /lib/systemd/system/gohjasmin.service
    systemctl enable gohjasmin.service
    systemctl start gohjasmin.service
    
Please see journalctl for errors in your log.

## Update your IP ##

Run gohjasmin.d in debug mode and it will dump some URLs you can use to update your DNS.

The quickest is (although not standard with any DDnS Clients):

    curl http://demouser:pwd@host:port/ohjasmin

### How to build **GohJasmin** ###

Get the code from here:

    git clone https://github.com/swordlordcodingcrew/gohjasmin.git

and build **GohJasmin** on your own with the following command:

    make build
        
## Dependencies ##

Dependencies are managed with the help of Go Modules. The Makefile as well as the Go toolchain will make sure that you will download the required dependencies if needed.

## License ##

**GohJasmin** is published under the GNU Affero General Public Licence version 3. See the LICENCE file for details.