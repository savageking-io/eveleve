# eveleve
Discord Bot for Indie Game Developers

# What is this?
This is a self-hosted Discord Bot that can do the following things:
* Watch your GitHub repositories (You need to configure Webhooks) and notify about events in a special channel
* Watch your Travis CI builds 
* Watch your Patreon page

# How to setup
* Create new Discord Applcation and enable bot. Save the token into configuration yaml file. 
* Add Bot to your server
* Using Ansible Playbook deploy eveleve to your server
* Make sure that your server have a domain name or public IP - GitHub needs it
* Visit GitHub and enable webhook for each repository from configuration file
