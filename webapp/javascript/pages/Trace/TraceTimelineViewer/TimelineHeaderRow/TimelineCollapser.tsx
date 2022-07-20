// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import React from 'react';
import Icon from '@webapp/ui/Icon';
import { faAngleRight } from '@fortawesome/free-solid-svg-icons/faAngleRight';
import { faAngleDown } from '@fortawesome/free-solid-svg-icons/faAngleDown';
import { faAngleDoubleDown } from '@fortawesome/free-solid-svg-icons/faAngleDoubleDown';
import { faAngleDoubleRight } from '@fortawesome/free-solid-svg-icons/faAngleDoubleRight';

import './TimelineCollapser.css';

type CollapserProps = {
  onCollapseAll: () => void;
  onCollapseOne: () => void;
  onExpandOne: () => void;
  onExpandAll: () => void;
};

export default class TimelineCollapser extends React.PureComponent<CollapserProps> {
  containerRef: React.RefObject<HTMLDivElement>;

  constructor(props: CollapserProps) {
    super(props);
    this.containerRef = React.createRef();
  }

  // TODO: Something less hacky than createElement to help TypeScript / AntD
  getContainer = () =>
    this.containerRef.current || document.createElement('div');

  render() {
    const { onExpandAll, onExpandOne, onCollapseAll, onCollapseOne } =
      this.props;
    return (
      <div className="TimelineCollapser" ref={this.containerRef}>
        <div aria-hidden="true" onClick={onExpandOne}>
          <Icon className="TimelineCollapser--btn" icon={faAngleDown} />
        </div>
        <div aria-hidden="true" onClick={onCollapseOne}>
          <Icon icon={faAngleRight} className="TimelineCollapser--btn" />
        </div>

        <div aria-hidden="true" onClick={onExpandAll}>
          <Icon icon={faAngleDoubleDown} className="TimelineCollapser--btn" />
        </div>

        <div aria-hidden="true" onClick={onCollapseAll}>
          <Icon icon={faAngleDoubleRight} className="TimelineCollapser--btn" />
        </div>
      </div>
    );
  }
}
