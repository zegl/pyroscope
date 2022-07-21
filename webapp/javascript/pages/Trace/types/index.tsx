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

import { Router } from 'react-router-dom';
import { Trace } from './trace';
import DetailState from '../TraceTimelineViewer/SpanDetail/DetailState';

export type Points = {
  x: number;
  y: number | null;
};

export type DataAvg = {
  service_operation_call_rate: null | number;
  service_operation_error_rate: null | number;
  service_operation_latencies: null | number;
};

export type OpsDataPoints = {
  service_operation_call_rate: Points[];
  service_operation_error_rate: Points[];
  service_operation_latencies: Points[];
  avg: DataAvg;
};

export type ServiceOpsMetrics = {
  dataPoints: OpsDataPoints;
  errRates: number;
  impact: number;
  latency: number;
  name: string;
  requests: number;
  key: number;
};

export type ServiceMetricsObject = {
  serviceName: string;
  quantile: number;
  max: number;
  metricPoints: Points[];
};

export type ServiceMetrics = {
  service_latencies: null | ServiceMetricsObject[];
  service_call_rate: null | ServiceMetricsObject;
  service_error_rate: null | ServiceMetricsObject;
};

export type MetricsReduxState = {
  serviceError: {
    service_latencies_50: null | ApiError;
    service_latencies_75: null | ApiError;
    service_latencies_95: null | ApiError;
    service_call_rate: null | ApiError;
    service_error_rate: null | ApiError;
  };
  opsError: {
    opsLatencies: null | ApiError;
    opsCalls: null | ApiError;
    opsErrors: null | ApiError;
  };
  isATMActivated: null | boolean;
  loading: boolean;
  operationMetricsLoading: null | boolean;
  serviceMetrics: ServiceMetrics | null;
  serviceOpsMetrics: ServiceOpsMetrics[] | undefined;
};

export type TTraceTimeline = {
  childrenHiddenIDs: Set<string>;
  detailStates: Map<string, DetailState>;
  hoverIndentGuideIds: Set<string>;
  shouldScrollToFirstUiFindMatch: boolean;
  spanNameColumnWidth: number;
  traceID: string | TNil;
};

export type ApiError =  // eslint-disable-line import/prefer-default-export
  | string
  | {
      message: string;
      httpStatus?: any;
      httpStatusText?: string;
      httpUrl?: string;
      httpQuery?: string;
      httpBody?: string;
    };

type SearchQuery = {
  end: number | string;
  limit: number | string;
  lookback: string;
  maxDuration: null | string;
  minDuration: null | string;
  operation: string | TNil;
  service: string;
  start: number | string;
  tags: string | TNil;
};

type TTraceDiffState = {
  a?: string | TNil;
  b?: string | TNil;
  cohort: string[];
};

export type TNil = null | undefined;

export type FetchedState = 'FETCH_DONE' | 'FETCH_ERROR' | 'FETCH_LOADING';

export type FetchedTrace = {
  data?: Trace;
  error?: ApiError;
  id: string;
  state?: FetchedState;
};

export type ReduxState = {
  config: Config;
  dependencies: {
    dependencies: { parent: string; child: string; callCount: number }[];
    loading: boolean;
    error: ApiError | TNil;
  };
  // embedded: EmbeddedState;
  router: typeof Router & {
    location: Location;
  };
  services: {
    services: string[] | TNil;
    serverOpsForService: Record<string, string[]>;
    operationsForService: Record<string, string[]>;
    loading: boolean;
    error: ApiError | TNil;
  };
  trace: {
    traces: Record<string, FetchedTrace>;
    search: {
      error?: ApiError;
      results: string[];
      state?: FetchedState;
      query?: SearchQuery;
    };
  };
  traceDiff: TTraceDiffState;
  traceTimeline: TTraceTimeline;
  metrics: MetricsReduxState;
};

interface ITimeCursorUpdate {
  cursor: number | TNil;
}

interface ITimeReframeUpdate {
  reframe: {
    anchor: number;
    shift: number;
  };
}

interface ITimeShiftEndUpdate {
  shiftEnd: number;
}

interface ITimeShiftStartUpdate {
  shiftStart: number;
}

export type TUpdateViewRangeTimeFunction = (
  start: number,
  end: number,
  trackSrc?: string
) => void;

export type ViewRangeTimeUpdate =
  | ITimeCursorUpdate
  | ITimeReframeUpdate
  | ITimeShiftEndUpdate
  | ITimeShiftStartUpdate;

export interface IViewRangeTime {
  current: [number, number];
  cursor?: number | TNil;
  reframe?: {
    anchor: number;
    shift: number;
  };
  shiftEnd?: number;
  shiftStart?: number;
}

export interface IViewRange {
  time: IViewRangeTime;
}

export enum ETraceViewType {
  TraceTimelineViewer = 'TraceTimelineViewer',
  TraceGraph = 'TraceGraph',
  TraceStatistics = 'TraceStatistics',
  TraceSpansView = 'TraceSpansView',
  TraceFlamegraph = 'TraceFlamegraph',
}

type ConfigMenuItem = {
  label: string;
  url: string;
  anchorTarget?: '_self' | '_blank' | '_parent' | '_top';
};

type ConfigMenuGroup = {
  label: string;
  items: ConfigMenuItem[];
};

type TScript = {
  text: string;
  type: 'inline';
};

type LinkPatternsConfig = {
  type: 'process' | 'tags' | 'logs' | 'traces';
  key?: string;
  url: string;
  text: string;
};

type MonitorEmptyStateConfig = {
  mainTitle?: string;
  subTitle?: string;
  description?: string;
  button?: {
    text?: string;
    onClick?: Function;
  };
  info?: string;
  alert?: {
    message?: string;
    type?: 'success' | 'info' | 'warning' | 'error';
  };
};

export type Config = {
  archiveEnabled?: boolean;
  deepDependencies?: {
    menuEnabled?: boolean;
  };
  dependencies?: { dagMaxServicesLen?: number; menuEnabled?: boolean };
  menu: (ConfigMenuGroup | ConfigMenuItem)[];
  pathAgnosticDecorations?: any[];
  qualityMetrics?: {
    menuEnabled?: boolean;
    menuLabel?: string;
  };
  search?: { maxLookback: { label: string; value: string }; maxLimit: number };
  scripts?: TScript[];
  topTagPrefixes?: string[];
  tracking?: {
    cookieToDimension?: {
      cookie: string;
      dimension: string;
    }[];
    gaID: string | TNil;
    trackErrors: boolean | TNil;
  };
  linkPatterns?: LinkPatternsConfig;
  monitor?: {
    menuEnabled?: boolean;
    emptyState?: MonitorEmptyStateConfig;
    docsLink?: string;
  };
};
