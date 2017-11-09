FROM quay.io/prometheus/busybox:latest

COPY prom-elasticsearch-exporter /bin/prom-elasticsearch-exporter

EXPOSE      9108
ENTRYPOINT  [ "/bin/prom-elasticsearch-exporter" ]
