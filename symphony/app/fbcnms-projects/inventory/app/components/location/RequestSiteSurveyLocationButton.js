/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  MarkSiteSurveyNeededMutationResponse,
  MarkSiteSurveyNeededMutationVariables,
} from '../../mutations/__generated__/MarkSiteSurveyNeededMutation.graphql';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import Button from '@fbcnms/ui/components/design-system/Button';
import MarkSiteSurveyNeededMutation from '../../mutations/MarkSiteSurveyNeededMutation';
import React from 'react';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {withStyles} from '@material-ui/core/styles';

const styles = _theme => ({});

type Props = {
  location: {id: string, siteSurveyNeeded: boolean},
} & WithAlert &
  WithStyles<typeof styles>;

class RequestSiteSurveyLocationButton extends React.Component<Props> {
  render() {
    const {location} = this.props;

    return (
      <Button onClick={this.requestSiteSurvey}>
        {location.siteSurveyNeeded
          ? 'Cancel Site Survey'
          : 'Request Site Survey'}
      </Button>
    );
  }

  requestSiteSurvey = () => {
    const {location} = this.props;

    const variables: MarkSiteSurveyNeededMutationVariables = {
      locationId: nullthrows(location.id),
      needed: !location.siteSurveyNeeded,
    };

    const callbacks: MutationCallbacks<MarkSiteSurveyNeededMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          this.props.alert('Failed requesting site survey');
        }
      },
      onError: (_error: Error) => {
        this.props.alert('Failed requesting site survey');
      },
    };

    MarkSiteSurveyNeededMutation(variables, callbacks);
  };
}

export default withStyles(styles)(withAlert(RequestSiteSurveyLocationButton));
