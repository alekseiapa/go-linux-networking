FROM ubuntu:latest

# Install iproute2 for ip command
RUN apt-get update && apt-get install -y iproute2

# Install iputils-ping for ping command
RUN apt-get install -y iputils-ping

# Install tcpdump
RUN apt-get install -y tcpdump

# Set the hostname
ARG CONTAINER_HOSTNAME
ENV HOSTNAME=$CONTAINER_HOSTNAME

# Update the PS1 environment variable to include the hostname
RUN echo 'export PS1="\[\e]0;\u@\h: \w\a\]\[\033[01;32m\]\u@\h\[\033[00m\]:\[\033[01;34m\]\w\[\033[00m\]\$ "' >> ~/.bashrc

CMD ["bash"]