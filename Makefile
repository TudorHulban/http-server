OUTPUT_BINARY = server-hx

compile:
	go build -o $(OUTPUT_BINARY) .

clean-cert:
	rm -f $(KEY_FILE) $(CERT_FILE) $(CERT_FILE).enc

COUNTRY := RO
STATE := Iasi
LOCALITY := Iasi
ORGANIZATION := TARA Works
DAYS := 365
KEY_FILE := key.pem
CERT_FILE := cert.pem

.PHONY: all clean cert

all: cert

cert: $(KEY_FILE) $(CERT_FILE)
cert-production: $(KEY_FILE).enc $(CERT_FILE)

$(KEY_FILE):
	openssl genrsa -out $@ 2048

$(CERT_FILE): $(KEY_FILE)
	openssl req -x509 -new -key $< -out $@ -days $(DAYS) -nodes \
		-subj "/C=$(COUNTRY)/ST=$(STATE)/L=$(LOCALITY)/O=$(ORGANIZATION)"

$(KEY_FILE).enc:
	openssl genrsa -aes256 -out $@ 2048

clean:
	rm -f $(KEY_FILE) $(CERT_FILE) $(KEY_FILE).enc

.DEFAULT_GOAL := all
