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
jsonResponse = ""

#######################################
# User Defined functions definitions
#######################################

def promote(sourceRepo: str,
            targetRepo: str,
            image: str,
            imageTag: str = "all"):

        #print(f"Started promoting {imageTag} tag of Image: {image}...")

        # Usage Example: fab promote:sourceRepo=feg-test,targetRepo=feg-prod,image=httpd,imageTag=1.0

        JFROG_API_KEY = _get_jfrog_apikey()

        try:
           headers = {
                       'content-type': 'application/json'
                     }
           response = requests.get(JFROG_URL + '/' + sourceRepo + '/v2/' + image + '/tags/list', 
                                   headers=headers, auth=(JFROG_USER, JFROG_API_KEY))
           if(response.status_code == 200):
              #print(f" Received Response Code={response.status_code}")
              
              if 'json' in response.headers.get('Content-Type'):
                 jsonResponse = response.json()
              else:
                 print(f"Response content from GET call is not in JSON format. Aborting...") 
                 abort("Exiting...")

              if ( imageTag != "all" ) and ( imageTag not in jsonResponse['tags'] ):
                   print(f"Image {image} tag {imageTag} does NOT EXIST in source repo:{sourceRepo}")
                   abort("Nothing to promote. Exiting...") 

              tags = jsonResponse['tags'] if imageTag == "all" else [imageTag]
              print(f"Started promoting {tags} tag of Image: {image}...")
              for tag in tags:
                  print(f"   Promoting {image}:{tag} from {sourceRepo} to {targetRepo}")
                  try:
                     payload = { 
                                 "targetRepo": targetRepo, 
                                 "dockerRepository": image, 
                                 "tag": tag, 
                                 "copy": "true" 
                               }
                     response = requests.post(JFROG_URL + '/' + sourceRepo + '/v2/promote', 
                                    headers=headers, json=payload, auth=(JFROG_USER, JFROG_API_KEY))
                     if(response.status_code == 200):
                        print(f"   SUCCESS promoting {image}:{tag} from {sourceRepo} to {targetRepo}.")
                        # print(f"   API Response Code={response.status_code}")
                        # print(f"   API Response Body={response.content}")
                     else:
                        print(f"FAILED promoting {image}:{tag} from {sourceRepo} to {targetRepo}.")
                        print(f"API Response Code={response.status_code}")
                        print(f"API Response Body={response.json()}")
                  except requests.exceptions.HTTPError as error:
                        print(f"   FAILED promoting {image}:{tag} from {sourceRepo} to {targetRepo}.")
                        print(f"   HTTPError Occurred on POST call. Error=")
                        print(error)
           else:
              print(f"GET call failed with status Code={response.status_code}. Response={response.json()}")
              abort("Exiting...")
        except requests.exceptions.HTTPError as error:
              print(f"HTTPError Occurred on GET call. Error=")
              print(error)
              abort("Exiting...")

def _get_jfrog_apikey():
    try:
        with open(os.path.expanduser(JFROG_APIKEY_PATH), 'r') as keyfile:
            return keyfile.read().rstrip()
    except IOError as e:
        print(e)
        abort("Please add the Jfrog API key for user %s to %s and re-run"
              % (JFROG_USER, JFROG_APIKEY_PATH))
