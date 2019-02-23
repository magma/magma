# Orchestrator Remote CLI - Adding commands to NMS UI

## Setup
Set up NMS according to the NMS Setup README.

Tip: To test NMS with your local magma setup, in `.env` you can use `API_HOST` set to https://192.168.80.10:9443/, `API_CERT_FILENAME` with the path that points to `admin_operator.pem` (located in  `magma/.cache/test_certs/admin_operator.pem`), and `API_PRIVATE_KEY_FILENAME` with the path that points to `admin_operator.key.pem` (located in `magma/.cache/test_certs/admin_operator.key.pem`).

## Adding commands to NMS UI
Use `MagmaAPIUrls.command()` to get the url to the command endpoint `/networks/{network_id}/gateways/{gateway_id}/command/{command_name}`.

We can then make a request using that url, for example:
```
const url = MagmaAPIUrls.command(match, id, commandName);

axios
  .post(url)
  .then(_resp => {
    this.props.alert('Success');
  })
  .catch(error => this.props.alert(error.response.data.message));
```