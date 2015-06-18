USER=root
IP=192.168.50.213
BINARY=/home/root/edison

all:
	make build
	make deploy
	make run

build:
	gox -arch="386" -os="linux"

deploy:
	scp edison_linux_386 $(USER)@$(IP):$(BINARY)

run:
	ssh $(USER)@$(IP) "INFLUX_HOST=$(INFLUX_HOST) INFLUX_PORT=$(INFLUX_PORT) INFLUX_USER=$(INFLUX_USER) INFLUX_PWD=$(INFLUX_PWD) $(BINARY)"
