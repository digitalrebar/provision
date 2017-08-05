#! /usr/bin/env python
# Copyright 2017, RackN
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# pip install requests
import requests, argparse, json, urllib3, os

'''
Usage: https://github.com/digitalrebar/provision/tree/master/integration/ansible

example: ansible -i inventory.py all -a "uname -a"
'''
    
def main():

    # change these values to match your DigitalRebar installation
    addr = os.getenv('RS_ENDPOINT', "https://127.0.0.1:8092")
    ups = os.getenv('RS_KEY', "rocketskates:r0cketsk8ts")
    arr = ups.split(":")
    user = arr[0]
    password = arr[1]

    # Argument parsing 
    parser = argparse.ArgumentParser(description="Ansible dynamic inventory via DigitalRebar")
    parser.add_argument("--list", help="Ansible inventory of all of the deployments", 
        action="store_true", dest="list_inventory")
    parser.add_argument("--host",
        help="Ansible inventory of a particular host", action="store",
        dest="ansible_host", type=str)

    cli_args = parser.parse_args()
    list_inventory = cli_args.list_inventory
    ansible_host = cli_args.ansible_host

    Headers = {'content-type': 'application/json'}
    urllib3.disable_warnings()
    print("# Digital Rebar URL " + addr + " via user " + user)

    profiles = {}
    profiles_raw = requests.get(addr + "/api/v3/profiles",headers=Headers,auth=(user,password),verify=False)
    if profiles_raw.status_code == 200: 
        for profile in profiles_raw.json():
            profiles[profile[u"Name"]] = [] 
    else:
        raise IOError(profiles_raw.text)

    if list_inventory:
        URL = addr + "/api/v3/machines"
    elif ansible_host:
        URL = addr + "/api/v3/machines?Name=" + ansible_host
    else:
        URL = addr + "/api/v3/machines"

    raw = requests.get(URL,headers=Headers,auth=(user,password),verify=False)

    if raw.status_code == 200: 
        for machine in raw.json():
            name = machine[u'Name']
            # TODO, should we only show machines that are in local bootenv?  others could be transistioning
            # if the machine has profiles, collect them
            if machine[u"Profiles"]:
                for profile in machine[u"Profiles"]:
                    profiles[profile].append(name)
            print name + " ansible_host=" + machine[u"Address"] 
    else:
        raise IOError(raw.text)

    for profile in profiles:
        print "[" + profile + "]"
        for machine in profiles[profile]:
            print machine

if __name__ == "__main__":
    main()  
