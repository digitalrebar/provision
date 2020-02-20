#!/usr/bin/env python
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
import argparse
import json
import logging
import os

try:
    import requests
    import urllib3
except ImportError:
    print("Missing requests. pip install requests")
    raise SystemExit

'''
Usage: https://github.com/digitalrebar/provision/tree/v4/integrations/ansible

example: ansible -i drpmachines.py all -a "uname -a"
'''

urllib3.disable_warnings()


def setup_logging():
    fmat = '%(asctime)-15s %(name)s %(message)s'
    logging.basicConfig(format=fmat, level=logging.DEBUG)


def setup_parser():
    parser = argparse.ArgumentParser(
        description="Ansible dynamic inventory via DigitalRebar"
    )
    parser.add_argument(
        "--list",
        help="Ansible inventory of all of the deployments",
        action="store_true",
        dest="list_inventory"
    )
    parser.add_argument(
        "--host",
        help="Ansible inventory of a particular host",
        action="store",
        dest="ansible_host",
        type=str
    )
    parser.add_argument(
        "--debug",
        action="store_true",
        help="Enable debug logging."
    )
    return parser.parse_args()


def main():
    cli_args = setup_parser()
    if cli_args.debug:
        setup_logging()
    # change these values to match your DigitalRebar installation
    drp_addr = os.getenv('RS_ENDPOINT', "https://127.0.0.1:8092")
    logging.debug("RS_ENDPOINT: {0}".format(drp_addr))
    ups = os.getenv('RS_KEY', "rocketskates:r0cketsk8ts")
    logging.debug("RS_KEY: {0}".format(ups))
    profile = os.getenv('RS_ANSIBLE', "all_machines")
    logging.debug("RS_ANSIBLE: {0}".format(profile))
    host_address = os.getenv('RS_HOST_ADDRESS', "internal")
    logging.debug("RS_HOST_ADDRESS: {0}".format(host_address))
    ansible_user = os.getenv('RS_ANSIBLE_USER', "root")
    logging.debug("RS_ANSIBLE_USER: {0}".format(ansible_user))
    parent_key = os.getenv('RS_ANSIBLE_PARENT', "ansible/children")
    logging.debug("RS_ANSIBLE_PARENT: {0}".format(parent_key))
    user, password = ups.split(":")

    list_inventory = cli_args.list_inventory
    ansible_host = cli_args.ansible_host
    headers = {'content-type': 'application/json'}
    inventory = {'all': {'hosts': []}, '_meta': {'hostvars': {}}}
    inventory["_meta"]["rebar_url"] = drp_addr
    inventory["_meta"]["rebar_user"] = user
    inventory["_meta"]["rebar_profile"] = profile

    hostvars = {}

    url = drp_addr + "/api/v3/machines"
    if list_inventory:
        if profile != "all_machines":
            url += "?ansible=Eq({0})".format(profile)
    else:
        if ansible_host:
            url += "?Name={0}".format(ansible_host)
        else:
            if profile != "all_machines":
                url += "?ansible=Eq({0})".format(profile)

    raw = requests.get(
        url,
        headers=headers,
        auth=(user, password),
        verify=False
    )

    ignore_params = ["gohai-inventory", "inventory/data", "change-stage/map"]
    if raw.status_code == 200:
        for machine in raw.json():
            ansible_user = os.getenv('RS_ANSIBLE_USER', "root")
            name = machine[u'Name']
            if name == '' or name is None:
                continue
            if machine[u'Address'] == '' or machine[u'Address'] is None:
                continue
            inventory["all"]["hosts"].extend([name])
            myvars = hostvars.copy()
            if host_address == "internal":
                myvars["ansible_host"] = machine[u"Address"]
            else:
                myvars["ansible_host"] = machine[u"Params"][host_address]
            ansible_user = machine.get("Params").get("ansible_user", ansible_user)
            myvars["ansible_user"] = ansible_user
            myvars["rebar_uuid"] = machine[u"Uuid"]
            for k in machine[u'Params']:
                if k not in ignore_params:
                    myvars[k] = machine[u'Params'][k]
            hvs = machine.get(u'Params').get(u'ansible/hostvars', None)
            if hvs is not None:
                for k, v in hvs.items():
                    myvars[k] = v

            inventory["_meta"]["hostvars"][name] = myvars
    else:
        raise IOError(raw.text)

    if ansible_host is None:
        groups = requests.get(
            drp_addr + "/api/v3/profiles",
            headers=headers,
            auth=(user, password),
            verify=False
        )
        if groups.status_code == 200:
            for group in groups.json():
                name = group[u'Name']
                if name != "global" and name != "rackn-license":
                    inventory[name] = {"hosts": [], "vars": []}
                    gvars = hostvars.copy()
                    if 'Profiles' in group.keys() \
                            and len(group[u'Profiles']) > 0:
                        inventory[name]["children"] = group[u'Profiles']
                    for k in group[u'Params']:
                        v = group[u'Params'][k]
                        if k == parent_key:
                            inventory[name]["children"] = v
                        else:
                            gvars[k] = v
                    inventory[name]["vars"] = gvars
                    hosts = requests.get(
                        "{0}/api/v3/machines?slim=Params&Profiles=In({1})".format(
                            drp_addr, name
                        ),
                        headers=headers,
                        auth=(user, password),
                        verify=False
                    )
                    if hosts.status_code == 200:
                        inventory[name]["hosts"] = []
                        for host in hosts.json():
                            hostname = host[u'Name']
                            inventory[name]["hosts"].extend([hostname])
        else:
            raise IOError(groups.text)        

    print(json.dumps(inventory))


if __name__ == "__main__":
    main()  
