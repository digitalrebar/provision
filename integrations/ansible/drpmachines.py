#! /usr/bin/env python
# Copyright 2019, RackN
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

example: ansible -i drpmachines.py all -a "uname -a"
'''

def main():

    inventory = { "_meta": { "hostvars": {} } }

    # change these values to match your DigitalRebar installation
    addr = os.getenv('RS_ENDPOINT', "https://127.0.0.1:8092")
    ups = os.getenv('RS_KEY', "rocketskates:r0cketsk8ts")
    profile = os.getenv('RS_ANSIBLE', "all_machines")
    parent_key = os.getenv('RS_ANSIBLE_PARENT', "ansible/children")
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
    inventory["_meta"]["rebar_url"] = addr
    inventory["_meta"]["rebar_user"] = user
    inventory["_meta"]["rebar_profile"] = profile
    inventory["_meta"]["rebar_profile"] = profile
    inventory["_meta"]["all"] = {'hosts': [], 'children': {}}

    groups = []
    profiles = {}
    profiles_vars = {}
    hostvars = {}

    URL = addr + "/api/v3/machines"
    if list_inventory:
        if profile != "all_machines":
            URL += "?ansible=Eq(" + profile + ")"
    elif ansible_host:
        URL += "?Name=" + ansible_host
    elif profile != "all_machines":
        URL += "?ansible=Eq(" + profile + ")"

    raw = requests.get(URL,headers=Headers,auth=(user,password),verify=False)

    IGNORE_PARAMS = ["gohai-inventory","inventory/data","change-stage/map"]
    if raw.status_code == 200: 
        for machine in raw.json():
            name = machine[u'Name']
            inventory["_meta"]["all"]["hosts"].extend([name])
            myvars = hostvars.copy()
            myvars["ansible_host"] = machine[u"Address"]
            myvars["rebar_uuid"] = machine[u"Uuid"]
            for k in machine[u'Params']:
                if k not in IGNORE_PARAMS:
                    myvars[k] = machine[u'Params'][k]
            inventory["_meta"]["hostvars"][name] = myvars
    else:
        raise IOError(raw.text)

    if ansible_host is None:
        groups = requests.get(addr + "/api/v3/profiles",headers=Headers,auth=(user,password),verify=False)
        if groups.status_code == 200:
            for group in groups.json():
                name = group[u'Name']
                if name != "global" and name != "rackn-license":
                    inventory["_meta"]["all"]["children"][name] = { "vars": [] }
                    gvars = hostvars.copy()
                    for k in group[u'Params']:
                        v = group[u'Params'][k]
                        if k == parent_key:
                            inventory["_meta"]["all"]["children"][name]["children"] = v
                        else:
                            gvars[k] = v
                    inventory["_meta"]["all"]["children"][name]["vars"] = gvars
                    hosts = requests.get(addr + "/api/v3/machines?slim=Params&Profiles=In("+name+")",headers=Headers,auth=(user,password),verify=False)
                    if hosts.status_code == 200:
                        inventory["_meta"]["all"]["children"][name]["hosts"] = []
                        for host in hosts.json():
                            hostname = host[u'Name']
                            inventory["_meta"]["all"]["children"][name]["hosts"].extend([hostname])
        else:
            raise IOError(groups.text)        

    print json.dumps(inventory)

if __name__ == "__main__":
    main()  
