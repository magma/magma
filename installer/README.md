# FB-Magma Infra Installer

Installer to configure infra for proper Magma installation

## Getting Started
One of the biggest challenge for anyone to get started with Magma is to satisfy the dependencies including host configuration, packages and environment variables required to successfully install Magma. We want to provide this Magma Infra installer that should take care of not just the pre-requisites but also configurations required before setting up a Magma installation.

This installer will set up the underlying infra for the Magma project up and running on your local machine for development and testing purposes.  Please see Production Deployment section for important notes on how to deploy the project on a production system with automation , HA and scalability. Please reach out to info@irsols.com for access to a pre-installed instance or for production setup.

### Prerequisites

A Baremetal or VM Host (with Intel VT-X) enabled with at least: 
```
	- 4 CPUs
        - 12GB RAM
        - 80 GB HDD Space
	- Centos 7 Minimum installed
```

### Instructions
This is a multi-stage installer with some manual intervention (eg for reboots). The installer will check if your host meets the requirements for proper Magma installation. It will download and setup the Vagrant, Docker, Virtual Box , Python3 , Docker-compose and other necessary packages that magma requires in the CORRECT order. Once all of the packages are installed you are all set to install the Magma software from the official documentation which is getting updated very frequently . 
Once you have cloned the official IRSOLS repo from "https://github.com/irsols-devops/magma.git" , please open up setup-magma-infra.sh to go through the installer and adjust to your environment if necessary

### Quick-Start
Clone the repository from irsols-devops github , change to magma/installer directory and run 'setup-magma-infra.sh'

```
# git clone https://github.com/irsols-devops/magma.git
# cd magma/installer/
# ./setup-magma-infra.sh

```

## Production Deployment

This installer is for a single host installation . Typical Magma installation doesnt account for multiple factors required in a production system , e.g Security, Scalability and HA. These factors need to be manually implemented or adjusted for specific service provider environment.Ideally this should be installed in an HA format and multiple instances to off-load the MNO EPC. 
IRSOLS can provide a production ready automated installer and can customize per specific environments . Please reach out to info@irsols.com for such requests. 


## Contributing

Please feel free to clone / fork the repo and let us know of any feature requests that you'd like to add with pull requests 

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **IRSOLS DevOps Team ** - *Initial Work* - [IRSOLS Inc](https://github.com/irsols-devops)
* ** Zeeshan Rizvi ** - [Zeeshan Rizvi] (https://linkedin.com/in/zeeshanrizvi)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Inspiration : FB Magma Open Source project


