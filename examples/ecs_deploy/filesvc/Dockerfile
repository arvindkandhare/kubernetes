FROM ubuntu:latest

RUN apt-get update -qq
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y build-essential libfuse-dev libcurl4-openssl-dev libxml2-dev mime-support automake libtool wget tar

RUN apt-get install git -y 
RUN apt-get install -y openssh-server
RUN mkdir /var/run/sshd
RUN echo 'root:12345' | chpasswd
RUN sed -i 's/PermitRootLogin without-password/PermitRootLogin yes/' /etc/ssh/sshd_config

# SSH login fix. Otherwise user is kicked off after login
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd

ENV NOTVISIBLE "in users profile"
RUN echo "export VISIBLE=now" >> /etc/profile

RUN useradd git -p 12345

RUN apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    lxc \
    iptables

# Install Docker from Docker Inc. repositories.
RUN curl -sSL https://get.docker.com/ubuntu/ | sh

# Install the magic wrapper.
ADD ./wrapdocker /usr/local/bin/wrapdocker
RUN chmod +x /usr/local/bin/wrapdocker

RUN git clone https://github.com/arvindkandhare/s3fs-fuse

RUN cd s3fs-fuse && ./autogen.sh && ./configure --prefix=/usr && make && make install




ADD runme /usr/local/bin/

RUN chmod +x /usr/local/bin/runme

ENTRYPOINT ["/usr/local/bin/runme"]
