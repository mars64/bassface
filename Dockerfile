FROM alpine
RUN apk update \
	&& apk upgrade \
	&& apk add --no-cache \
	ca-certificates \
	&& update-ca-certificates 2>/dev/null || true
WORKDIR /
ADD bassface.amd64 /bassface
CMD /bassface
