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

import _get from 'lodash/get';
import memoizeOne from 'memoize-one';
import deepFreeze from 'deep-freeze';
import processDeprecation from './process-deprecation';

let haveWarnedFactoryFn = false;
let haveWarnedDeprecations = false;

const defaultConfig = deepFreeze(
  Object.defineProperty(
    {
      archiveEnabled: false,
      dependencies: {
        dagMaxNumServices: 100,
        menuEnabled: true,
      },
      linkPatterns: [],
      qualityMetrics: {
        menuEnabled: false,
        menuLabel: 'Trace Quality',
      },
      menu: [
        {
          label: 'About Jaeger',
          items: [
            {
              label: 'Website/Docs',
              url: 'https://www.jaegertracing.io/',
            },
            {
              label: 'Blog',
              url: 'https://medium.com/jaegertracing/',
            },
            {
              label: 'Twitter',
              url: 'https://twitter.com/JaegerTracing',
            },
            {
              label: 'Discussion Group',
              url: 'https://groups.google.com/forum/#!forum/jaeger-tracing',
            },
            {
              label: 'Online Chat',
              url: 'https://cloud-native.slack.com/archives/CGG7NFUJ3',
            },
            {
              label: 'GitHub',
              url: 'https://github.com/jaegertracing/',
            },
            {
              label: `Jaeger`,
            },
            {
              label: `Commit`,
            },
            {
              label: `Build`,
            },
            {
              label: `Jaeger UI v${1}`,
            },
          ],
        },
      ],
      search: {
        maxLookback: {
          label: '2 Days',
          value: '2d',
        },
        maxLimit: 1500,
      },
      tracking: {
        gaID: null,
        trackErrors: true,
        customWebAnalytics: null,
      },
      monitor: {
        menuEnabled: true,
        emptyState: {
          mainTitle: 'Get started with Service Performance Monitoring',
          subTitle:
            'A high-level monitoring dashboard that helps you cut down the time to identify and resolve anomalies and issues.',
          description:
            'Service Performance Monitoring aggregates tracing data into RED metrics and visualizes them in service and operation level dashboards.',
          button: {
            text: 'Read the Documentation',
            onClick: () =>
              window.open('https://www.jaegertracing.io/docs/latest/spm/'),
          },
          alert: {
            message:
              'Service Performance Monitoring requires a Prometheus-compatible time series database.',
            type: 'info',
          },
        },
        docsLink: 'https://www.jaegertracing.io/docs/latest/spm/',
      },
    },
    // fields that should be individually merged vs wholesale replaced
    '__mergeFields',
    { value: ['dependencies', 'search', 'tracking'] }
  )
);

export const deprecations = [
  {
    formerKey: 'dependenciesMenuEnabled',
    currentKey: 'dependencies.menuEnabled',
  },
  {
    formerKey: 'gaTrackingID',
    currentKey: 'tracking.gaID',
  },
];

const getConfig = memoizeOne(function getConfig() {
  const getJaegerUiConfig = window.getJaegerUiConfig;
  if (typeof getJaegerUiConfig !== 'function') {
    if (!haveWarnedFactoryFn) {
      // eslint-disable-next-line no-console
      console.warn('Embedded config not available');
      haveWarnedFactoryFn = true;
    }
    return { ...defaultConfig };
  }
  const embedded = getJaegerUiConfig();
  if (!embedded) {
    return { ...defaultConfig };
  }
  // check for deprecated config values
  if (Array.isArray(deprecations)) {
    deprecations.forEach((deprecation) =>
      processDeprecation(embedded, deprecation, !haveWarnedDeprecations)
    );
    haveWarnedDeprecations = true;
  }
  const rv = { ...defaultConfig, ...embedded };
  // __mergeFields config values should be merged instead of fully replaced
  const keys = defaultConfig.__mergeFields || [];
  for (let i = 0; i < keys.length; i++) {
    const key = keys[i];
    if (typeof embedded[key] === 'object' && embedded[key] !== null) {
      rv[key] = { ...defaultConfig[key], ...embedded[key] };
    }
  }
  return rv;
});

export default getConfig;

export function getConfigValue(path: string) {
  return _get(getConfig(), path);
}
