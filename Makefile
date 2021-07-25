.PHONY: all clean helm

all: build docker helm

clean: 
	-rm bassface bassface.amd64
	-docker rmi mars64/bassface
	UNTAGGED := docker images | grep "mars64/bassface.*<none>" | awk '{print $3}'
	$(foreach image,$(UNTAGGED), docker rmi $(image))

build:
	go build -o bassface.darwin64
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bassface.amd64 .

docker:
	docker build -t mars64/bassface:latest .
	docker push mars64/bassface:latest

helm:
	helm upgrade --install mars64-bassface helm/bassface --set badWords='${BASSFACE_BADWORDS}' --set discogsToken='${BASSFACE_DISCOGS_TOKEN}' --set password='${BASSFACE_PASSWORD}' --set join=${BASSFACE_JOIN_CHANNEL} --set nick=${BASSFACE_NICK} --set reportTo="${BASSFACE_REPORT_TO}" --set server="${BASSFACE_SERVER}"
