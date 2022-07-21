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

/* eslint-disable jsx-a11y/anchor-is-valid, react/jsx-props-no-spreading */
import * as React from 'react';
import _get from 'lodash/get';
import _maxBy from 'lodash/maxBy';
import _values from 'lodash/values';
import MdKeyboardArrowRight from 'react-icons/lib/md/keyboard-arrow-right';
import SpanGraph from './SpanGraph';
import {
  TUpdateViewRangeTimeFunction,
  IViewRange,
  ViewRangeTimeUpdate,
} from '../types';
import LabeledList from '../common/LabeledList';
import TraceName from '../common/TraceName';
import { getTraceName } from './trace-viewer';
import { Trace } from '../types/trace';
import { formatDatetime, formatDuration } from '../utils/date';

import './TracePageHeader.css';

type TracePageHeaderEmbedProps = {
  canCollapse: boolean;
  hideMap: boolean;
  hideSummary: boolean;
  onSlimViewClicked: () => void;
  slimView: boolean;
  trace: Trace;
  updateNextViewRangeTime: (update: ViewRangeTimeUpdate) => void;
  updateViewRangeTime: TUpdateViewRangeTimeFunction;
  viewRange: IViewRange;
};

export const HEADER_ITEMS = [
  {
    key: 'timestamp',
    label: 'Trace Start',
    renderer: (trace: Trace) => {
      const dateStr = formatDatetime(trace.startTime);
      const match = dateStr.match(/^(.+)(\.\d+)$/);
      return match ? (
        <span className="TracePageHeader--overviewItem--value">
          {match[1]}
          <span className="TracePageHeader--overviewItem--valueDetail">
            {match[2]}
          </span>
        </span>
      ) : (
        dateStr
      );
    },
  },
  {
    key: 'duration',
    label: 'Duration',
    renderer: (trace: Trace) => formatDuration(trace.duration),
  },
  {
    key: 'service-count',
    label: 'Services',
    renderer: (trace: Trace) =>
      new Set(_values(trace.processes).map((p) => p.serviceName)).size,
  },
  {
    key: 'depth',
    label: 'Depth',
    renderer: (trace: Trace) =>
      _get(_maxBy(trace.spans, 'depth'), 'depth', 0) + 1,
  },
  {
    key: 'span-count',
    label: 'Total Spans',
    renderer: (trace: Trace) => trace.spans.length,
  },
];

export default function TracePageHeaderFn(props: TracePageHeaderEmbedProps) {
  const {
    canCollapse,
    hideMap,
    hideSummary,
    onSlimViewClicked,
    slimView,
    trace,
    updateNextViewRangeTime,
    updateViewRangeTime,
    viewRange,
  } = props;

  if (!trace) {
    return null;
  }

  const summaryItems =
    !hideSummary &&
    !slimView &&
    HEADER_ITEMS.map((item) => {
      const { renderer, ...rest } = item;
      return { ...rest, value: renderer(trace) };
    });

  const title = (
    <h1
      className={`TracePageHeader--title ${
        canCollapse ? 'is-collapsible' : ''
      }`}
    >
      <TraceName traceName={getTraceName(trace.spans)} />{' '}
      <small className="u-tx-muted">{trace.traceID.slice(0, 7)}</small>
    </h1>
  );

  return (
    <header className="TracePageHeader">
      <div className="TracePageHeader--titleRow">
        {canCollapse ? (
          <a
            className="TracePageHeader--titleLink"
            onClick={onSlimViewClicked}
            role="switch"
            aria-checked={!slimView}
            tabIndex={0}
            aria-hidden
          >
            <MdKeyboardArrowRight
              className={`TracePageHeader--detailToggle ${
                !slimView ? 'is-expanded' : ''
              }`}
            />
            {title}
          </a>
        ) : (
          title
        )}
      </div>
      {summaryItems && (
        <LabeledList
          className="TracePageHeader--overviewItems"
          items={summaryItems}
        />
      )}
      {!hideMap && !slimView && (
        <SpanGraph
          trace={trace}
          viewRange={viewRange}
          updateNextViewRangeTime={updateNextViewRangeTime}
          updateViewRangeTime={updateViewRangeTime}
        />
      )}
    </header>
  );
}
