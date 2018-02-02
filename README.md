GohJasmin
======

![Goh Jasmin](https://raw.githubusercontent.com/LordEidi/gohjasmin/master/ohjasmin.jpg)

**GohJasmin** (c) 2017-18 by [SwordLord - the coding crew](https://www.swordlord.com/)

Based on OhJasmin (c) 2010-18 by [SwordLord - the coding crew](https://www.swordlord.com/)

## Introduction ##

**GohJasmin** is a lightweight, self-hostable dynamic DNS server written in Golang. Just add an instance of [djbdns](https://cr.yp.to/djbdns.html) or dbndns and you have your very own dynamic DNS service with your own domains. 

(You could actually add any other DNS server to the mix. You will just need to write some glue code to export the dynamic data from the GohJasmin database into your DNS server. Thats not very difficult.)

If you are looking for a service to manage your dynamic dns records and have your hosts send their current IP addresses to, then **GohJasmin** might be for you.

_This is still work in progress! Use at your own risk._


## Installation ##

### How to build **GohJasmin** ###

While we work on a binary release, you may 

    git clone https://github.com/LordEidi/gohjasmin.git

this repository and build **GohJasmin** on your own with the following command:

    gb build all
    
## Configuration ##

All parameters which can be configured right now are in the file *gohjasmin.config.js*. A default configuration file will be written on the first run of the application.

## Dependencies ##

Dependencies are managed in the "vendor" folder. Just go to the root of this project and "gb build all" to compile the projects binaries (for your platform).

## License ##

**GohJasmin** is published under the GNU Affero General Public Licence version 3. See the LICENCE file for details.