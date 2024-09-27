FROM centos:7

LABEL org.opencontainers.image.authors="liuzihao<2115883273@qq.com>"

RUN mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup && \
    curl -o /etc/yum.repos.d/CentOS-Base.repo https://mirrors.aliyun.com/repo/Centos-7.repo && \
    sed -i 's/gpgcheck=1/gpgcheck=0/g' /etc/yum.repos.d/CentOS-Base.repo

# 更新系统并安装常用工具
RUN yum clean all && \
    yum makecache && \
    yum update -y && \
    yum install -y \
        sudo \
        which \
        vim \
        wget \
        curl \
        net-tools \
        iputils \
        tar \
        gzip \
        unzip \
        git \
        openssh-server \
        passwd \
        cronie \
        yum-utils \
        epel-release && \
    yum clean all

ENV PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:${PATH}"

RUN useradd -m -s /bin/bash lijingwoquan && \
    echo "lijingwoquan:123456" | chpasswd && \
    usermod -aG wheel lijingwoquan

RUN echo "lijingwoquan ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

USER lijingwoquan

WORKDIR /

CMD ["/bin/bash"]
