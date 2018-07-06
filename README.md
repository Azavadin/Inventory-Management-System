# Cloud_computing_Fabic
Frabic inventory
Item list(JSON format front end):
  key: item name + catalog#
  Fields:
    Name:
    Catalog:
    Comments:
    Purchase_date:
    Supplier:
    Purchase_by:


Required packages for Fabic under Ubuntu 16.04 :

1. Docker:

    a. Install dependencies for Docker:
        sudo apt install apt-transport-https ca-certificates curl software-properties-common

    b. Install Docker: 
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add - && \
    sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" && \
    sudo apt update && \
    sudo apt install -y docker-ce

    c. Add the current user to the docker group:
    sudo groupadd docker ; \
    sudo gpasswd -a ${USER} docker && \
    sudo service docker restart

    d. logout and log in again

2. Docker composer:

    a. Install:
    sudo curl -L https://github.com/docker/compose/releases/download/1.18.0/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose && \
    sudo chmod +x /usr/local/bin/docker-compose
    
3. Go:

    a. Install:
    wget https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz && \
    sudo tar -C /usr/local -xzf go1.9.2.linux-amd64.tar.gz && \
    rm go1.9.2.linux-amd64.tar.gz && \
    echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile && \
    echo 'export GOPATH=$HOME/go' | tee -a $HOME/.bashrc && \
    echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' | tee -a $HOME/.bashrc && \
    mkdir -p $HOME/go/{src,pkg,bin}
    

For Making the First Blockchain network in your system:
    
    a. Make a new directory in the src folder of your GOPATH:
    mkdir -p $GOPATH/src/github.com/
    
    b.Copy the SimpleInventory package and paste under $GOPATH/src/github.com
       
    c. On the terminal go in SimpleInventory directory
    cd $GOPATH/src/github.com/SimpleInventory
    
    d.Run make command.
    Once it starts listening to port 3000.
    Open url: http://localhost:3000/ in browser.
    
    
    

    
    
    
    
