zero_channels_for_one_cbsd = {
    "spectrumInquiryResponse": [
        {
            "response": {
                "responseCode": 0
            },
            "cbsdId": "foo",
            "availableChannel": []
        }
    ]
}

single_channel_for_one_cbsd = {
    "spectrumInquiryResponse": [
        {
            "response": {
                "responseCode": 0
            },
            "cbsdId": "foo",
            "availableChannel": [
                {
                    "frequencyRange": {
                        "lowFrequency": 1,
                        "highFrequency": 1
                    },
                    "channelType": "test",
                    "ruleApplied": "test",
                    "maxEirp": 1,
                }
            ]
        }
    ]
}

two_channels_for_one_cbsd = {
    "spectrumInquiryResponse": [
        {
            "response": {
                "responseCode": 0
            },
            "cbsdId": "foo",
            "availableChannel": [
                {
                    "frequencyRange": {
                        "lowFrequency": 1,
                        "highFrequency": 10
                    },
                    "channelType": "test",
                    "ruleApplied": "test",
                    "maxEirp": 1,
                },
                {
                    "frequencyRange": {
                        "lowFrequency": 20,
                        "highFrequency": 30
                    },
                    "channelType": "test1",
                    "ruleApplied": "test1",
                    "maxEirp": 2,
                }
            ]
        }
    ]
}
