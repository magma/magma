/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Button from '@fbcnms/ui/components/design-system/Button';
import React, {useState} from 'react';
import RelayEnvironment from '../common/RelayEnvironment';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import {commitMutation, graphql} from 'react-relay';
import type {
  CreateServiceMutation,
  CreateServiceMutationResponse,
  NetworkServiceInput,
} from './__generated__/CreateServiceMutation.graphql';

type formValues = {
  serviceName: string,
  customerName: string,
  serviceModelName: string,
  siteNumber: number,
  siteModelName: string,
  deviceName: string,
  deviceId: number,
  interfaceId: string,
  vlanId: string,
};

type CreateServiceResultTagProps = {
  response: CreateServiceMutationResponse,
};

const mutation = graphql`
  mutation CreateServiceMutation($nsi: NetworkServiceInput!) {
    createService(input: $nsi)
  }
`;

function commit(
  nsi: NetworkServiceInput,
  responseSetter: CreateServiceMutationResponse => void,
) {
  return commitMutation<CreateServiceMutation>(RelayEnvironment, {
    mutation,
    variables: {nsi},
    onError: (err: Error) => console.error(err),
    onCompleted: (response: CreateServiceMutationResponse, errors) => {
      if (errors) {
        console.error(errors);
      } else {
        responseSetter(response);
      }
    },
  });
}

type CreateServiceFormProps = {
  vals: formValues,
  setVals: formValues => void,
};

function CreateServiceForm(props: CreateServiceFormProps) {
  const {vals, setVals} = props;
  const genOnChange = (fieldName: string) => event =>
    setVals({...vals, [fieldName]: event.target.value});
  return (
    <>
      <Text>Service Name</Text>
      <TextInput
        value={vals.serviceName}
        onChange={genOnChange('serviceName')}
      />
      <Text>Service Model Name</Text>
      <TextInput
        value={vals.serviceModelName}
        onChange={genOnChange('serviceModelName')}
      />
      <Text>Customer</Text>
      <TextInput
        value={vals.customerName}
        onChange={genOnChange('customerName')}
      />
      <Text>Site Number</Text>
      <TextInput
        value={vals.siteNumber}
        onChange={e =>
          setVals({
            ...vals,
            siteNumber: e.target.value == '' ? 0 : parseInt(e.target.value),
          })
        }
      />
      <Text>Site Model Name</Text>
      <TextInput
        value={vals.siteModelName}
        onChange={genOnChange('siteModelName')}
      />
      <Text>Device Name</Text>
      <TextInput value={vals.deviceName} onChange={genOnChange('deviceName')} />
      <Text>Device ID</Text>
      <TextInput
        value={vals.deviceId}
        onChange={e =>
          setVals({
            ...vals,
            deviceId: e.target.value == '' ? 0 : parseInt(e.target.value),
          })
        }
      />
      <Text>Interface ID</Text>
      <TextInput
        value={vals.interfaceId}
        onChange={genOnChange('interfaceId')}
      />
      <Text>VLAN ID</Text>
      <TextInput value={vals.vlanId} onChange={genOnChange('vlanId')} />
    </>
  );
}

function CreateServiceResultTag(props: CreateServiceResultTagProps) {
  const parsed = JSON.parse(props.response.createService)[0];
  return (
    <>
      <Text>status: {parsed.status}</Text> <br />
      <Text>message: {parsed.message}</Text> <br />
      <Text>id: {parsed.id}</Text>
    </>
  );
}

function CreateService() {
  const [response, setResponse] = useState<?CreateServiceMutationResponse>();
  const [vals, setVals] = useState<formValues>({
    serviceName: '',
    customerName: '',
    serviceModelName: '',
    siteNumber: 0,
    siteModelName: '',
    deviceName: '',
    deviceId: 0,
    interfaceId: '',
    vlanId: '',
  });

  function createServiceButtonPressed() {
    const nsi: NetworkServiceInput = {
      Customer: vals.customerName,
      Name: vals.serviceName,
      Model: {
        Name: vals.serviceModelName,
      },
      DeviceSites: [
        {
          siteNumber: vals.siteNumber,
          siteModelName: vals.siteModelName,
          deviceName: vals.deviceName,
          deviceId: vals.deviceId,
          parameters: {
            interfaceID: vals.interfaceId,
            vlanID: vals.vlanId,
          },
          userPort: {
            id: -1,
            name: '',
          },
          accessMethod: 'CLI',
        },
      ],
    };
    commit(nsi, setResponse);
  }
  return (
    <>
      <CreateServiceForm vals={vals} setVals={setVals} />
      <Button onClick={createServiceButtonPressed}>Create Service</Button>
      {response && (
        <>
          <br />
          <CreateServiceResultTag response={response} />
        </>
      )}
    </>
  );
}

export default CreateService;
