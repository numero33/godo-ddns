# GODO-DDNS

Golang DigitalOcean DynDNS

> Works for Synology
![image](https://user-images.githubusercontent.com/12273891/58832261-daec4880-864e-11e9-880c-b001c9c38322.png)

## Usage

```bash
$ curl 'https://yourdomain.com/update?h=__HOSTNAME__&ip=__MYIP__&p=__PASSWORD__'
```

- \_\_HOSTNAME\_\_: Hostname
- \_\_MYIP\_\_: IPv4 address
- \_\_PASSWORD\_\_: DigitalOcean API Token [Apps & API](https://cloud.digitalocean.com/settings/applications) (https://cloud.digitalocean.com/settings/applications)


## Docker
[numero33/godo-ddns](https://hub.docker.com/r/numero33/godo-ddns)
```bash
$ docker run -d -p 80:8080 -e GODOHOST=8080 numero33/godo-ddns
```

## Build
```bash
$ go build -o godo-ddns
$ chmod +x godo-ddns
$ GODOHOST=":80" ./godo-ddns
```
