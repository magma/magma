registration_requests = [
    {
        "registrationRequest": [
            {
                "fccId": "foo",
                "cbsdCategory": "B",
                "callSign": "WSD987",
                "userId": "John Doe",
                "airInterface": {
                    "radioTechnology": "E_UTRA"
                },
                "cbsdSerialNumber": "4321dcba",
                "measCapability": [
                    "RECEIVED_POWER_WITHOUT_GRANT"
                ],
                "installationParam": {
                    "latitude": 37.425056,
                    "longitude": -122.084113,
                    "height": 9.3,
                    "heightType": "AGL",
                    "indoorDeployment": False,
                    "antennaAzimuth": 271,
                    "antennaDowntilt": 3,
                    "antennaGain": 16,
                    "antennaBeamwidth": 30
                },
                "groupingParam": [
                    {
                        "groupId": "example-group-3",
                        "groupType": "INTERFERENCE_COORDINATION"
                    }
                ]
            }
        ]
    },
    {
        "registrationRequest": [
            {
                "fccId": "bar",
                "cbsdCategory": "A",
                "callSign": "CB987",
                "userId": "John Doe",
                "airInterface": {
                    "radioTechnology": "E_UTRA"
                },
                "cbsdSerialNumber": "abcd1234",
                "measCapability": [
                    "RECEIVED_POWER_WITHOUT_GRANT"
                ],
                "installationParam": {
                    "latitude": 37.419735,
                    "longitude": -122.072205,
                    "height": 6,
                    "heightType": "AGL",
                    "indoorDeployment": True
                },
                "groupingParam": [
                    {
                        "groupId": "example-group-1",
                        "groupType": "INTERFERENCE_COORDINATION"
                    },
                    {
                        "groupId": "example-group-2",
                        "groupType": "INTERFERENCE_COORDINATION"
                    }
                ]
            },
            {
                "fccId": "baz",
                "cbsdCategory": "A",
                "callSign": "CB987",
                "userId": "John Doe",
                "airInterface": {
                    "radioTechnology": "E_UTRA"
                },
                "cbsdSerialNumber": "abcd1234",
                "measCapability": [
                    "RECEIVED_POWER_WITHOUT_GRANT"
                ],
                "installationParam": {
                    "latitude": 45.419735,
                    "longitude": -100.072205,
                    "height": 6,
                    "heightType": "AGL",
                    "indoorDeployment": True
                },
                "groupingParam": [
                    {
                        "groupId": "example-group-1",
                        "groupType": "INTERFERENCE_COORDINATION"
                    },
                    {
                        "groupId": "example-group-2",
                        "groupType": "INTERFERENCE_COORDINATION"
                    }
                ]
            }
        ]
    }
]
