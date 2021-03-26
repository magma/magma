#!/usr/bin/python

import sys
import requests
from requests.auth import HTTPBasicAuth

#################################
# Global variables definitions
#################################

JFROG_URL = 'http://artifactory.magmacore.org.com:8081/artifactory/api/docker'
JFROG_KEY = 'AP8nwtQtPAf5M3yXD3dSosJbLjkvFicWMUN2BF'
reqHeaders = { "Content-Type": "application/json", "X-JFrog-Art-Api": "AKCp8ihewtkTv45dL9qevJFGQkhPr4xAU5k69VhsfQJvsGBT1X646DWxo8yL3KzmbhudBiCXv" }
sourceRepo = ""
targetRepo = ""
imgName = "gateway_go"
imgVersion = "1.0"
jsonResponse = ""


#######################################
# User Defined functions definitions
#######################################

def callJFROGAPI(method='GET',
                 url=JFROG_URL,
                 reqHeaders={ "Content-Type": "application/json", "X-JFrog-Art-Api": JFROG_KEY },
                 payload=""):

      print(f"callJFROGAPI() invoked with \n method:{method}\n url={url}\n reqHeaders={reqHeaders}\n payload={payload}")
    if method == 'GET':
        try:
           response = requests.get(url + '/' + sourceRepo + '/v2/' + imgName + '/tags/list', headers=reqHeaders)
           if(response.status_code == 200):
              print(f" Received Response Code={response.status_code}")
              jsonResponse = response.json()
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
           response = requests.post(jfrogURL + '/' + sourceRepo + '/v2/promote', headers=reqHeaders, json=payload)
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
    status = callJFROGAPI(method='POST', payload={ "targetRepo": targetRepo, "dockerRepository": imgName, "tag": imgVersion })
    return status


if __name__ == "__main__":
    modules = [ 'docker', 'all' ]
    argCount = len(sys.argv)

    if argCount != 3:
       print(f"Invalid arguments. Usage: {sys.argv[0]} docker|feg|cwf|all version")
       quit()

    if sys.argv[1] not in modules:
       print(f"Module {sys.argv[1]} not found. Usage: {sys.argv[0]} docker|feg|cwf|all version")
       quit()

    module = sys.argv[1]
    imgVersion = sys.argv[2]
    imgName = "gateway_go"
    sourceRepo = module + "-local"
    targetRepo = module + "-kubernetesapps"

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
