/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 1931547eb5dd8b93d17fa5f5bd47bafc
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type LocationCellScanCoverageMap_cellData$ref = any;
type LocationWiFiScanCoverageMap_wifiData$ref = any;
export type LocationCoverageMapTabQueryVariables = {|
  locationId: string
|};
export type LocationCoverageMapTabQueryResponse = {|
  +location: ?{|
    +cellData?: $ReadOnlyArray<?{|
      +$fragmentRefs: LocationCellScanCoverageMap_cellData$ref
    |}>,
    +wifiData?: $ReadOnlyArray<?{|
      +$fragmentRefs: LocationWiFiScanCoverageMap_wifiData$ref
    |}>,
  |}
|};
export type LocationCoverageMapTabQuery = {|
  variables: LocationCoverageMapTabQueryVariables,
  response: LocationCoverageMapTabQueryResponse,
|};
*/


/*
query LocationCoverageMapTabQuery(
  $locationId: ID!
) {
  location: node(id: $locationId) {
    __typename
    ... on Location {
      cellData {
        ...LocationCellScanCoverageMap_cellData
        id
      }
      wifiData {
        ...LocationWiFiScanCoverageMap_wifiData
        id
      }
    }
    id
  }
}

fragment LocationCellScanCoverageMap_cellData on SurveyCellScan {
  id
  latitude
  longitude
  networkType
  signalStrength
  mobileCountryCode
  mobileNetworkCode
  operator
}

fragment LocationWiFiScanCoverageMap_wifiData on SurveyWiFiScan {
  id
  latitude
  longitude
  frequency
  channel
  bssid
  ssid
  strength
  band
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "locationId",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "locationId"
  }
],
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitude",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitude",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "LocationCoverageMapTabQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "location",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "InlineFragment",
            "type": "Location",
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "cellData",
                "storageKey": null,
                "args": null,
                "concreteType": "SurveyCellScan",
                "plural": true,
                "selections": [
                  {
                    "kind": "FragmentSpread",
                    "name": "LocationCellScanCoverageMap_cellData",
                    "args": null
                  }
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "wifiData",
                "storageKey": null,
                "args": null,
                "concreteType": "SurveyWiFiScan",
                "plural": true,
                "selections": [
                  {
                    "kind": "FragmentSpread",
                    "name": "LocationWiFiScanCoverageMap_wifiData",
                    "args": null
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "LocationCoverageMapTabQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "location",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "__typename",
            "args": null,
            "storageKey": null
          },
          (v2/*: any*/),
          {
            "kind": "InlineFragment",
            "type": "Location",
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "cellData",
                "storageKey": null,
                "args": null,
                "concreteType": "SurveyCellScan",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  (v4/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "networkType",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "signalStrength",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "mobileCountryCode",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "mobileNetworkCode",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "operator",
                    "args": null,
                    "storageKey": null
                  }
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "wifiData",
                "storageKey": null,
                "args": null,
                "concreteType": "SurveyWiFiScan",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  (v4/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "frequency",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "channel",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "bssid",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "ssid",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "strength",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "band",
                    "args": null,
                    "storageKey": null
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "LocationCoverageMapTabQuery",
    "id": null,
    "text": "query LocationCoverageMapTabQuery(\n  $locationId: ID!\n) {\n  location: node(id: $locationId) {\n    __typename\n    ... on Location {\n      cellData {\n        ...LocationCellScanCoverageMap_cellData\n        id\n      }\n      wifiData {\n        ...LocationWiFiScanCoverageMap_wifiData\n        id\n      }\n    }\n    id\n  }\n}\n\nfragment LocationCellScanCoverageMap_cellData on SurveyCellScan {\n  id\n  latitude\n  longitude\n  networkType\n  signalStrength\n  mobileCountryCode\n  mobileNetworkCode\n  operator\n}\n\nfragment LocationWiFiScanCoverageMap_wifiData on SurveyWiFiScan {\n  id\n  latitude\n  longitude\n  frequency\n  channel\n  bssid\n  ssid\n  strength\n  band\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '16ca9fb7b04e14fdb62144f24b5ab0e5';
module.exports = node;
