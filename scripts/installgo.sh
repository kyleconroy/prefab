if [ ! -e /usr/local/go ]; then
	wget -q http://go.googlecode.com/files/go1.1.1.linux-amd64.tar.gz
	tar -C /usr/local -xzf go1.1.1.linux-amd64.tar.gz
	rm -f go1.1.1.linux-amd64.tar.gz 
fi

mkdir -p /home/vagrant/go/src/github.com/stackmachine

if [ ! -e /home/vagrant/go/src/github.com/stackmachine/stackgo ]; then
	ln -s /vagrant /home/vagrant/go/src/github.com/stackmachine/stackgo
fi

echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/go.sh
echo 'export GOPATH=/home/vagrant/go' >> /etc/profile.d/go.sh
