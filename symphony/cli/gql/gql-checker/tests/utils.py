import re


ERROR_RX = re.compile("# ((GQL[0-9]+ ?)+)(: (.*))?$")


def extract_expected_errors(data):
    lines = data.splitlines()
    expected_codes = []
    expected_messages = []
    for line in lines:
        match = ERROR_RX.search(line)
        if match:
            codes = match.group(1).split()
            message = match.group(4)
            expected_codes.extend(codes)
            if message:
              expected_messages.append(message)
    return expected_codes, expected_messages
