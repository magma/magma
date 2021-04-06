#!/usr/bin/python3

import os
import json
import sys
import requests
from requests.auth import HTTPBasicAuth
from fabric.utils import abort

#################################
# Global variables definitions
#################################

JFROG_URL = 'https://artifactory.magmacore.org/artifactory/api/docker'
JFROG_USER = 'ci-bot'
# echo $JFROG_CIBOT_APIKEYS > ~/.magma/secrets/jfrog_apikey
JFROG_APIKEY_PATH = "~/.magma/secrets/jfrog_apikey"
JFROG_API_KEY = ""

sourceRepo = ""
targetRepo = ""
imgName = "httpd"
imgVersion = "1.0"
jsonResponse = ""

#######################################
# User Defined functions definitions
#######################################

def callJFROGAPI(method='GET',
                 url=JFROG_URL,
                 reqHeaders={ "Content-Type": "application/json" },
                 payload=""):

    print(f"callJFROGAPI() invoked with \n method:{method}\n url={url}\n reqHeaders={reqHeaders}\n payload={payload}")
    if method == 'GET':
        try:
           response = requests.get(url + '/' + sourceRepo + '/v2/' + imgName + '/tags/list', headers=reqHeaders, auth=(JFROG_USER, JFROG_API_KEY))
           if(response.status_code == 200):
              print(f" Received Response Code={response.status_code}")
              
              if 'json' in response.headers.get('Content-Type'):
                 jsonResponse = response.json()
              else:
                 print(f"Response content is not in JSON format.") 
                 return -1

              print(f"Image:{imgName} list of tags:{jsonResponse['tags']}")
              if imgVersion in jsonResponse["tags"]:
                 print(f"{imgName}:{imgVersion} found in sourece repo:{sourceRepo}. Promoting it to {targetRepo}")
              else:
                 print(f"{imgName}:{imgVersion} NOT found in sourece repo:{sourceRepo}. NOT Promoting. Quitting...")
                 return -1
           else:
              print(f"GET call failed with status Code={response.status_code}. Response={response.json()}")
              return -1
        except requests.exceptions.HTTPError as error:
              print(f"HTTPError Occurred. Error=")
              print(error)
              return -1
    elif method == 'POST':
        try:
           response = requests.post(url + '/' + sourceRepo + '/v2/promote', headers=reqHeaders, json=payload, auth=(JFROG_USER, JFROG_API_KEY))
           if(response.status_code == 200):
              print(f"SUCCESS promoting {imgName}:{imgVersion} from {sourceRepo} to {targetRepo}.")
              print(f"Response Code={response.status_code}")
              print(f"Response Body={response.content}")
           else:
              print(f"Response Code={response.status_code}")
              print(f"Response Body={response.json()}")
              return -1
        except requests.exceptions.HTTPError as error:
              print("HTTPError Occurred. Error=")
              print(error)
              return -1
    else:
        print(f"INCORRECT API method {method} invoked. Please use GET or POST.")
        return -1

    return 0

def validate(module):
    print(f"Validating {imgName}:{imgVersion} in repop:{sourceRepo}")
    status = callJFROGAPI(method='GET')
    return status


def promote(module):
    print(f"Promoting {imgName}:{imgVersion} from {sourceRepo} to {targetRepo}")
    status = callJFROGAPI(method='POST', payload={ "targetRepo": targetRepo, "dockerRepository": imgName, "tag": imgVersion, "copy": "true" })
    return status

def _get_jfrog_apikey():
    try:
        with open(os.path.expanduser(JFROG_APIKEY_PATH), 'r') as keyfile:
            return keyfile.read().rstrip()
    except IOError as e:
        print(e)
        abort("Please add the Jfrog API key for user %s to %s and re-run"
              % (JFROG_USER, JFROG_APIKEY_PATH))

JFROG_API_KEY = _get_jfrog_apikey()

if __name__ == "__main__":
    modules = [ 'docker', 'feg', 'cwf', 'orc8r', 'all' ]
    argCount = len(sys.argv)

    if argCount != 4:
       print(f"Invalid arguments. Usage: {sys.argv[0]} {modules} imageName version. Ex: promote.py feg gateway_go 1a4b2b30")
       abort("Exiting...")

    if sys.argv[1] not in modules:
       print(f"Module {sys.argv[1]} not found in allowed list {modules}")
       print(f"Usage: {sys.argv[0]} {modules} imageName version")
       abort("Exiting...")

    module = sys.argv[1]
    imgName = sys.argv[2]
    imgVersion = sys.argv[3]
    sourceRepo = module + "-test"
    targetRepo = module + "-prod"


    status = validate(module)
    if status == 0:
       print(f"Validation of {imgName}:{imgVersion} in repop:{sourceRepo} is successful.")
    else:
       print(f"Validation of {imgName}:{imgVersion} in repop:{sourceRepo} FAILED. CANNOT Promote. Quitting...")
       quit()

    status = promote(module)
    if status == 0:
       print("Promotion is successful")
    else:
       print("Promotion FAILED")
