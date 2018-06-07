FROM ubuntu:latest
EXPOSE 22

# Setup ssh
RUN apt-get update
RUN apt-get install -y openssh-server ssh bash sudo
# RUN install the md5 tools hashdeep?
RUN mkdir /var/run/sshd; chmod 0755 /var/run/sshd

# user account
RUN adduser --disabled-password --gecos "" test_user
RUN echo 'test_user:password' | chpasswd
RUN usermod -aG sudo test_user

# ssh key
RUN mkdir /home/test_user/.ssh/
ADD testing_key.rsa.pub /home/test_user/.ssh/authorized_keys
RUN chmod -R 0700 /home/test_user/.ssh/
RUN chmod -R 0644 /home/test_user/.ssh/authorized_keys
RUN chown -R test_user:test_user /home/test_user/.ssh/

CMD ["/usr/sbin/sshd", "-D"]
