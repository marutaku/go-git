FROM ubuntu:20.04

RUN apt-get update && apt-get install -y \
  curl \
  unzip \
  libssl-dev \
  zlib1g-dev \
  make \
  gcc \ 
  bsdmainutils \ 
  python3



RUN curl -sL https://github.com/git/git/archive/e83c5163316f89bfbde7d9ab23ca2e25604af290.zip -o git.zip \
  && mkdir -p /${HOME}/bin \
  && unzip git.zip \
  && mv git-* original-git \
  && cd original-git \
  && sed -i 's/LIBS= -lssl/LIBS= -lssl -lz -lcrypto/' Makefile \
  && make \
  && make install

ENV PATH="/${HOME}/bin:${PATH}"

COPY decompress.py /original-git/decompress.py

CMD [ "/bin/bash" ]