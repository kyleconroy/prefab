if [ ! -e /usr/local/go ]; then
	wget -q http://go.googlecode.com/files/go1.1.1.linux-amd64.tar.gz
	tar -C /usr/local -xzf go1.1.1.linux-amd64.tar.gz
	rm -f go1.1.1.linux-amd64.tar.gz 
fi

echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/go.sh
