This is the result of a hackathon with the Intel Edison.

What it does?
-------------

The idea behind this hack is that when you are in a conference, meetup, talk...
the organisers usually want to track the quality of the given talk. With this
project you will have:

- 2 buttons: vote the talk positive or negative
- Sound sensor to track the clapping

This information is going to be stored in a InfluxDB.

Part of the hackathon presentation was showind this data using a Grafana
frontend over that InfluxDB, but that's up to you. It's not 100% needed.

Requirements
------------

- Go
- An InfluxDB somewhere

If you have problems installing Influx do as I did and just use docker:

    docker run -d -p 8083:8083 -p 8086:8086 --expose 8090 --expose 8099 \
        --name influxdb -e PRE_CREATE_DB="intelmaker" tutum/influxdb

And Grafana (if you wish):

    docker run -d -p 3000:3000 --link influxdb:influxdb \
        --name grafana grafana/grafana

Deploy & run
------------

You will need to `make` previously setting up some environment variables:

    INFLUX_HOST=a INFLUX_PORT=b INFLUX_USER=c INFLUX_PWD=d

Or you can add them to the `Makefile` if you want.

Demo
----

You can see it working here: https://www.youtube.com/watch?v=UzrkYxbiYnY
