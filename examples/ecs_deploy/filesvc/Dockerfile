FROM ubuntu

RUN apt-get update -qq
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y build-essential libfuse-dev libcurl4-openssl-dev libxml2-dev mime-support automake libtool wget tar

RUN apt-get install git -y 
RUN git clone https://github.com/s3fs-fuse/s3fs-fuse

RUN cd s3fs-fuse && ./autogen.sh && ./configure --prefix=/usr && make && make install


ADD runme /usr/local/bin/

RUN chmod +x /usr/local/bin/runme

ENTRYPOINT ["/usr/local/bin/runme"]
