/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Suggestion} from '@fbcnms/ui/components/Typeahead';

import * as React from 'react';
import RelayEnvironment from '../../common/RelayEnvironment.js';
import Typeahead from '@fbcnms/ui/components/Typeahead';
import {fetchQuery, graphql} from 'relay-runtime';

type Props = {
  className?: string,
  required?: boolean,
  headline?: string,
  selectedProject?: ?{
    id: string,
    name: string,
  },
  margin?: ?string,
  onProjectSelection: (?{id: string, name: string}) => void,
};

type State = {
  projectSuggestions: Array<Suggestion>,
  projects: Array<Suggestion>,
};

const projectSearchQuery = graphql`
  query ProjectTypeahead_ProjectsQuery(
    $limit: Int
    $filters: [ProjectFilterInput!]!
  ) {
    projectSearch(limit: $limit, filters: $filters) {
      id
      name
      type {
        name
      }
    }
  }
`;

class ProjectTypeahead extends React.Component<Props, State> {
  state = {
    projectSuggestions: [],
    projects: [],
  };
  componentDidMount() {
    fetchQuery(RelayEnvironment, projectSearchQuery, {
      limit: 1000,
      filters: [],
    }).then(response => {
      if (!response || !response.projectSearch) {
        return;
      }
      this.setState({
        projects: response.projectSearch.map(p => ({
          name: p.name,
          entityId: p.id,
          entityType: 'project',
          type: p?.type.name,
        })),
      });
    });
  }

  fetchNewProjectSuggestions = (searchTerm: string) => {
    const searchTermLC = searchTerm.toLowerCase();
    const projects = this.state.projects;
    const suggestions = projects.filter(e =>
      e.name.toLowerCase().includes(searchTermLC),
    );
    const projectSuggestions = suggestions;
    this.setState({
      projectSuggestions,
    });
  };

  render() {
    const {selectedProject, headline, required, className} = this.props;
    const {projectSuggestions} = this.state;
    return (
      <div className={className}>
        <Typeahead
          margin={this.props.margin}
          required={!!required}
          suggestions={projectSuggestions}
          onSuggestionsFetchRequested={this.fetchNewProjectSuggestions}
          onEntitySelected={suggestion =>
            this.props.onProjectSelection({
              id: suggestion.entityId,
              name: suggestion.name,
            })
          }
          onEntriesRequested={() => {}}
          onSuggestionsClearRequested={() =>
            this.props.onProjectSelection(null)
          }
          placeholder={headline}
          value={
            selectedProject
              ? {
                  name: selectedProject.name,
                  entityId: selectedProject.id,
                  entityType: '',
                  type: '',
                }
              : null
          }
        />
      </div>
    );
  }
}

export default ProjectTypeahead;
